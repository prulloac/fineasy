package pkg

import (
	"github.com/google/uuid"
)

func HashPassword(password string, salt string, algorithm string) string {
	return ""
}

func GenerateSalt() string {
	return uuid.NewString()
}
