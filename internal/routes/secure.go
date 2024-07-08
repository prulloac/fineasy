package routes

import (
	"crypto/rsa"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prulloac/fineasy/internal/auth"
)

func createJwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing Authorization header"})
			return
		}
		log.Printf("🔐 Token: %s", tokenString)

		token, err := jwt.Parse(strings.TrimPrefix(tokenString, "Bearer "), func(token *jwt.Token) (interface{}, error) {
			return loadVerifyKey(), nil
		}, jwt.WithAudience("fineasy"), jwt.WithIssuer("fineasy"))
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		if !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		c.Set("token", token)
		c.Next()
	}
}

func generateBearerToken(user auth.User) string {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": user.Hash,
		"iss": "fineasy",
		"aud": "fineasy",
		"exp": now.Add(time.Hour * 24).Unix(),
		"iat": now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, _ := token.SignedString(loadSignKey())
	return "Bearer " + tokenString
}

func loadSignKey() *rsa.PrivateKey {
	if os.Getenv("JWT_PRIVATE_KEY") == "" {
		panic("JWT_PRIVATE_KEY not set")
	}
	privateKey := os.Getenv("JWT_PRIVATE_KEY")
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		log.Fatalf("error parsing private key: %s", err)
	}
	return signKey
}

func loadVerifyKey() *rsa.PublicKey {
	if os.Getenv("JWT_PUBLIC_KEY") == "" {
		panic("JWT_PUBLIC_KEY not set")
	}
	publicKey := os.Getenv("JWT_PUBLIC_KEY")
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		log.Fatalf("error parsing public key: %s", err)
	}
	return verifyKey
}
