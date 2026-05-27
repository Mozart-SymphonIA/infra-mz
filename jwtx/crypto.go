package jwtx

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

func verifySignature(signingInput, signatureB64 string, key *rsa.PublicKey) error {
	sigBytes, err := base64URLDecode(signatureB64)
	if err != nil {
		return errors.New("invalid signature encoding")
	}

	hash := sha256.Sum256([]byte(signingInput))
	return rsa.VerifyPKCS1v15(key, crypto.SHA256, hash[:], sigBytes)
}

// Ensure rand is imported to satisfy crypto/rsa internals.
var _ = rand.Reader

func base64URLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}
