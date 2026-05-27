package httpx

import (
	"context"
	"net/http"
	"strings"

	"github.com/Mozart-SymphonIA/infra-mz/jwtx"
)

type ctxKey struct{}

// IdentityIDFromContext returns the identity UUID injected by NewAuthMiddleware.
func IdentityIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ctxKey{}).(string)
	return id, ok
}

// NewAuthMiddleware returns an HTTP middleware that:
//   - extracts the JWT from the named cookie or Authorization: Bearer header,
//   - verifies it with the provided Verifier,
//   - checks that the caller holds requiredPermission,
//   - injects the identity UUID into the request context.
//
// Returns 401 on missing/invalid token, 403 on insufficient permission.
func NewAuthMiddleware(v *jwtx.Verifier, cookieName, requiredPermission string) func(http.Handler) http.Handler {
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

			if !claims.HasPermission(requiredPermission) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), ctxKey{}, claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
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
