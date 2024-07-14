package pkg

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"log"
	random "math/rand"

	"golang.org/x/crypto/sha3"
)

func HashPassword(password string, salt string, algorithm string) string {
	log.Printf("🔐 Hashing password with algorithm: %s", algorithm)
	var out string
	switch algorithm {
	case "SHA256":
		out = hashWithSHA256(password, salt)
	case "SHA512":
		out = hashWithSHA512(password, salt)
	case "SHA3_256":
		out = hashWithSHA3_256(password, salt)
	case "SHA3_512":
		out = hashWithSHA3_512(password, salt)
	default:
		panic("invalid algorithm")
	}
	return out[:24] // 24 characters
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

func hashWithSHA512(password string, salt string) string {
	concatenated := []byte(password + salt)
	ch := sha512.Sum512([]byte(concatenated))
	out := hex.EncodeToString(ch[:])
	log.Printf("🔐 Hashed password: %s", out)
	return out
}

func hashWithSHA3_256(password string, salt string) string {
	concatenated := []byte(password + salt)
	ch := sha3.Sum256([]byte(concatenated))
	out := hex.EncodeToString(ch[:])
	log.Printf("🔐 Hashed password: %s", out)
	return out
}

func hashWithSHA3_512(password string, salt string) string {
	concatenated := []byte(password + salt)
	ch := sha3.Sum512([]byte(concatenated))
	out := hex.EncodeToString(ch[:])
	log.Printf("🔐 Hashed password: %s", out)
	return out
}
