package jwtx

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

// verifyRS256 validates the signature and claims of an RS256 JWT.
// Intentionally avoids external JWT libraries to keep infra-mz dependency-free.
func verifyRS256(token string, key *rsa.PublicKey) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}

	if err := verifySignature(parts[0]+"."+parts[1], parts[2], key); err != nil {
		return Claims{}, ErrInvalidToken
	}

	payload, err := base64URLDecode(parts[1])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	return parseClaims(payload)
}

func parseClaims(payload []byte) (Claims, error) {
	var raw map[string]any
	if err := json.Unmarshal(payload, &raw); err != nil {
		return Claims{}, ErrInvalidToken
	}

	exp, err := toInt64(raw["exp"])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	expiresAt := time.Unix(exp, 0).UTC()
	if time.Now().UTC().After(expiresAt) {
		return Claims{}, ErrTokenExpired
	}

	iat, _ := toInt64(raw["iat"])
	tgID, err := toInt64(raw["tg"])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	sub, _ := raw["sub"].(string)
	username, _ := raw["u"].(string)
	perms := parsePerms(raw["p"])

	return Claims{
		Subject:     sub,
		TelegramID:  tgID,
		Username:    username,
		Permissions: perms,
		IssuedAt:    time.Unix(iat, 0).UTC(),
		ExpiresAt:   expiresAt,
	}, nil
}

func parsePerms(v any) []string {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, p := range raw {
		if s, ok := p.(string); ok && s != "" {
			out = append(out, s)
		}
	}
	return out
}

func toInt64(v any) (int64, error) {
	switch x := v.(type) {
	case float64:
		return int64(x), nil
	case int64:
		return x, nil
	case json.Number:
		return x.Int64()
	case string:
		return strconv.ParseInt(x, 10, 64)
	default:
		return 0, errors.New("not a number")
	}
}
