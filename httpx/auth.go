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

// NewAuthMiddleware returns an HTTP middleware that:
//   - extracts the JWT from the named cookie or Authorization: Bearer header,
//   - verifies it with the provided Verifier,
//   - checks that the caller holds requiredPermission,
//   - injects the full Claims into the request context.
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

			ctx := context.WithValue(r.Context(), ctxKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// NewSharedAccountMiddleware returns an HTTP middleware that checks whether
// the caller has the given action on the shared account identified by the
// {id} path parameter. It relies on Claims already being in the context
// (i.e. must be chained after NewAuthMiddleware).
//
// action matches the operation name in snake_case,
// e.g. "get_shared_account_summary".
func NewSharedAccountMiddleware(action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			accountID := r.PathValue("id")
			if !claims.HasSharedAccountPermission(accountID, action) {
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
