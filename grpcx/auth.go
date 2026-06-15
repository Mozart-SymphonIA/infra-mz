package grpcx

import (
	"context"
	"strings"
	"unicode"

	"github.com/Mozart-SymphonIA/infra-mz/jwtx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type claimsKey struct{}

// ClaimsFromContext returns the JWT claims injected by NewAuthInterceptor.
func ClaimsFromContext(ctx context.Context) (jwtx.Claims, bool) {
	c, ok := ctx.Value(claimsKey{}).(jwtx.Claims)
	return c, ok
}

// NewAuthInterceptor returns a gRPC unary interceptor that verifies the JWT
// from the incoming metadata, checks that the caller has the permission
// required for the invoked method (looked up from actionPolicy), and injects
// the full Claims into the context for downstream handlers.
//
// actionPolicy maps snake_case action names to required permissions,
// e.g. "register_transaction" → "finance:basic".
// It is typically loaded at startup from Ganesha's GetActionPolicy RPC.
func NewAuthInterceptor(v *jwtx.Verifier, actionPolicy map[string]string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization")
		}

		raw := strings.TrimPrefix(tokens[0], "Bearer ")
		claims, err := v.Verify(raw)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		action := methodToAction(info.FullMethod)
		requiredPermission, ok := actionPolicy[action]
		if !ok {
			return nil, status.Errorf(codes.PermissionDenied, "no policy defined for action: %s", action)
		}

		if !claims.HasPermission(requiredPermission) {
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}

		ctx = context.WithValue(ctx, claimsKey{}, claims)
		return handler(ctx, req)
	}
}

// NewSharedAccountInterceptor returns a gRPC unary interceptor that checks
// whether the caller has the invoked action on the shared account referenced
// in the request. It relies on Claims already being in the context (i.e. must
// run after NewAuthInterceptor).
//
// The request must implement GetSharedAccountId() string; if it does not, the
// interceptor passes through without checking (not a shared-account method).
func NewSharedAccountInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		sar, ok := req.(interface{ GetSharedAccountId() string })
		if !ok || sar.GetSharedAccountId() == "" {
			return handler(ctx, req)
		}

		claims, ok := ClaimsFromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing claims")
		}

		action := methodToAction(info.FullMethod)
		if !claims.HasSharedAccountPermission(sar.GetSharedAccountId(), action) {
			return nil, status.Error(codes.PermissionDenied, "not authorized for this shared account")
		}

		return handler(ctx, req)
	}
}

// methodToAction converts a gRPC FullMethod to a snake_case action name.
// e.g. "/finance.FinanceService/RegisterTransaction" → "register_transaction"
func methodToAction(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	method := parts[len(parts)-1]
	var b strings.Builder
	for i, r := range method {
		if unicode.IsUpper(r) && i > 0 {
			b.WriteRune('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}
