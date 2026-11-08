package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(bcryptInput(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(hash, plain string) bool {
	if strings.HasPrefix(hash, "$2b$") || strings.HasPrefix(hash, "$2y$") {
		hash = "$2a$" + hash[4:]
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), bcryptInput(plain)) == nil
}

func bcryptInput(plain string) []byte {
	if len(plain) <= 72 {
		return []byte(plain)
	}
	digest := sha256.Sum256([]byte(plain))
	return []byte(hex.EncodeToString(digest[:]))
}
