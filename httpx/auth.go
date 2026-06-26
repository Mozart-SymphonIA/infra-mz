package httpx

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/Mozart-SymphonIA/infra-mz/jwtx"
)

type ctxKey struct{}

var rePathParam = regexp.MustCompile(`\{[^}]+\}/?`)

// IdentityIDFromContext returns the identity UUID injected by NewAuthMiddleware.
func IdentityIDFromContext(ctx context.Context) (string, bool) {
	c, ok := ctx.Value(ctxKey{}).(jwtx.Claims)
	if !ok {
		return "", false
	}
	return c.Subject, true
}

// ClaimsFromContext returns the full JWT claims injected by NewAuthMiddleware.
func ClaimsFromContext(ctx context.Context) (jwtx.Claims, bool) {
	c, ok := ctx.Value(ctxKey{}).(jwtx.Claims)
	return c, ok
}

// RouteToAction converts an HTTP route pattern to a snake_case action name,
// mirroring the gRPC methodToAction convention.
//
// Examples:
//
//	"GET /api/finance/summary"                       → "get_finance_summary"
//	"GET /api/finance/phantom-expenses"              → "get_finance_phantom_expenses"
//	"GET /api/finance/shared-accounts/{id}/summary"  → "get_finance_shared_accounts_summary"
func RouteToAction(pattern string) string {
	parts := strings.SplitN(pattern, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	method := strings.ToLower(parts[0])
	path := parts[1]

	path = strings.TrimPrefix(path, "/api/")
	path = rePathParam.ReplaceAllString(path, "")
	path = strings.NewReplacer("/", "_", "-", "_").Replace(path)
	path = strings.Trim(path, "_")

	return method + "_" + path
}

// NewAuthMiddleware returns an HTTP middleware that:
//   - extracts the JWT from the named cookie or Authorization: Bearer header,
//   - verifies it with the provided Verifier,
//   - derives the required permission from actionPolicy using the matched route pattern,
//   - injects the full Claims into the request context.
//
// actionPolicy maps snake_case action names to required permissions, as loaded
// from Ganesha's GetActionPolicy RPC. Returns 401 on missing/invalid token,
// 403 on unknown route or insufficient permission.
func NewAuthMiddleware(v *jwtx.Verifier, cookieName string, actionPolicy map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := extractToken(r, cookieName)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := v.Verify(token)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			action := RouteToAction(r.Pattern)
			requiredPermission, ok := actionPolicy[action]
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			if !claims.HasPermission(requiredPermission) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), ctxKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// NewSharedAccountMiddleware returns an HTTP middleware that checks whether
// the caller holds the required per-account permission for the matched route.
// It relies on Claims already being in the context (must chain after NewAuthMiddleware).
//
// policy maps snake_case action names to the per-account permission they require,
// as loaded from Ganesha's GetSharedAccountActionPolicy RPC. Routes not present
// in the policy pass through without shared-account checking.
func NewSharedAccountMiddleware(policy map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			action := RouteToAction(r.Pattern)
			requiredPermission, inPolicy := policy[action]
			if !inPolicy {
				next.ServeHTTP(w, r)
				return
			}

			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			accountID := r.PathValue("id")
			if !claims.HasSharedAccountPermission(accountID, requiredPermission) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func extractToken(r *http.Request, cookieName string) (string, bool) {
	if c, err := r.Cookie(cookieName); err == nil && c.Value != "" {
		return c.Value, true
	}
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer "), true
	}
	return "", false
}
