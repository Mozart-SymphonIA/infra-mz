// Package jwtx provides shared JWT verification for Mozart-SymphonIA services.
// It fetches the public key from Ganesha's JWKS endpoint, caches it, and
// verifies RS256 tokens without calling Ganesha on every request.
package jwtx

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sync"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

// Claims represents the decoded payload of a Ganesha-issued JWT.
type Claims struct {
	Subject     string
	TelegramID  int64
	Username    string
	Permissions []string
	IssuedAt    time.Time
	ExpiresAt   time.Time
}

// Verifier fetches the JWKS from Ganesha once, caches the public key,
// and verifies tokens locally on every call.
type Verifier struct {
	jwksURL    string
	httpClient *http.Client

	mu     sync.RWMutex
	pubKey *rsa.PublicKey
}

// NewVerifier creates a Verifier that fetches the public key from jwksURL.
// Typical value: "http://ganesha:5180/.well-known/jwks.json"
func NewVerifier(jwksURL string) *Verifier {
	return &Verifier{
		jwksURL:    jwksURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Verify parses and validates a JWT string, returning the decoded Claims.
// The public key is fetched lazily on first call and cached thereafter.
func (v *Verifier) Verify(token string) (Claims, error) {
	key, err := v.publicKey()
	if err != nil {
		return Claims{}, fmt.Errorf("jwtx: fetch public key: %w", err)
	}
	return verifyRS256(token, key)
}

// RefreshKey forces a re-fetch of the JWKS (e.g. after a key rotation).
func (v *Verifier) RefreshKey() error {
	key, err := v.fetchKey()
	if err != nil {
		return err
	}
	v.mu.Lock()
	v.pubKey = key
	v.mu.Unlock()
	return nil
}

func (v *Verifier) publicKey() (*rsa.PublicKey, error) {
	v.mu.RLock()
	cached := v.pubKey
	v.mu.RUnlock()
	if cached != nil {
		return cached, nil
	}
	return v.fetchAndCache()
}

func (v *Verifier) fetchAndCache() (*rsa.PublicKey, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	// Double-check after acquiring write lock.
	if v.pubKey != nil {
		return v.pubKey, nil
	}
	key, err := v.fetchKey()
	if err != nil {
		return nil, err
	}
	v.pubKey = key
	return key, nil
}

func (v *Verifier) fetchKey() (*rsa.PublicKey, error) {
	resp, err := v.httpClient.Get(v.jwksURL)
	if err != nil {
		return nil, fmt.Errorf("get jwks: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read jwks: %w", err)
	}

	var jwks struct {
		Keys []struct {
			N string `json:"n"`
			E string `json:"e"`
		} `json:"keys"`
	}
	if err := json.Unmarshal(body, &jwks); err != nil || len(jwks.Keys) == 0 {
		return nil, errors.New("invalid jwks response")
	}

	k := jwks.Keys[0]
	nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
	if err != nil {
		return nil, fmt.Errorf("decode n: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
	if err != nil {
		return nil, fmt.Errorf("decode e: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := int(new(big.Int).SetBytes(eBytes).Int64())

	return &rsa.PublicKey{N: n, E: e}, nil
}
