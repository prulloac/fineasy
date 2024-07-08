package pkg

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	random "math/rand"
)

func HashPassword(password string, salt string, algorithm string) string {
	switch algorithm {
	case "SHA256":
		return hashWithSHA256(password, salt)
	default:
		panic("invalid algorithm")
	}
}

func GenerateSalt() string {
	stdchars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#?!@$%^&*-" // 74 characters
	l := 16 + random.Intn(10)
	salt := make([]byte, l)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}
	for i, b := range salt {
		salt[i] = stdchars[b%byte(len(stdchars))]
	}
	return string(salt)
}

func hashWithSHA256(password string, salt string) string {
	concatenated := []byte(password + salt)
	ch := sha256.Sum256([]byte(concatenated))
	out := hex.EncodeToString(ch[:])
	log.Printf("🔐 Hashed password: %s", out)
	return out
}
