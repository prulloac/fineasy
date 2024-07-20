package middleware

import (
	"crypto/rsa"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prulloac/fineasy/internal/auth"
	"github.com/prulloac/fineasy/pkg/logging"
)

var logger = logging.NewLoggerWithPrefix("[Middleware]")

func CaptureTokenFromHeader(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.AbortWithStatusJSON(401, gin.H{"error": "missing Authorization header"})
		return
	}

	token, err := jwt.Parse(strings.TrimPrefix(tokenString, "Bearer "), func(token *jwt.Token) (interface{}, error) {
		return loadVerifyKey(), nil
	}, jwt.WithAudience("fineasy"), jwt.WithIssuer("fineasy"), jwt.WithExpirationRequired())
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

func GenerateBearerToken(user *auth.User) string {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  user.Hash,
		"iss":  "fineasy",
		"aud":  "fineasy",
		"exp":  now.Add(time.Hour * 24).Unix(),
		"iat":  now.Unix(),
		"uid":  user.ID,
		"mail": user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, _ := token.SignedString(loadSignKey())
	return "Bearer " + tokenString
}

func GetClientToken(c *gin.Context) *jwt.Token {
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client user not found"})
		return nil
	}
	return token.(*jwt.Token)
}

func loadSignKey() *rsa.PrivateKey {
	if os.Getenv("JWT_PRIVATE_KEY") == "" {
		panic("JWT_PRIVATE_KEY not set")
	}
	privateKey := os.Getenv("JWT_PRIVATE_KEY")
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		logger.Fatalf("error parsing private key: %s", err)
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
		logger.Fatalf("error parsing public key: %s", err)
	}
	return verifyKey
}
