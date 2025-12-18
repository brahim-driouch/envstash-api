package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plaintext password using bcrypt
func HashPassword(password string) (string, error) {
	// Generate hash with default cost (10)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// CheckPassword compares a plaintext password with a hash
func ComparePasswords(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// jwt

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // use env variable in production

type JWTClaims struct {
	UserSub TokenSub
	jwt.RegisteredClaims
}

type TokenSub struct {
	Id         string `json:"id"`
	Fullname   string `json:"fullname"`
	Email      string `json:"email"`
	IsVerified bool   `json:"isVerified"`
	IsAdmin    bool   `json:"isAdmin"`
}

// Generate a JWT token for a user
func GenerateToken(userSub TokenSub, minutes int32) (string, error) {
	claims := JWTClaims{
		UserSub: userSub,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(minutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Verify JWT token and return claims
func VerifyToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func getAuthorizationHeader(c *gin.Context) (string, error) {
	tokenHeader := c.GetHeader("Authorization")

	// checl for AUTHORIZATION header
	if tokenHeader == "" {
		return "", errors.New(" authorization header required")

	}

	// check for authorzation format
	parts := strings.SplitN(tokenHeader, "", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization format")
	}
	return parts[0], nil

}
func AuthMidleware(c *gin.Context) {
	token, err := getAuthorizationHeader(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// check is token is valid
	_, err = VerifyToken(token)
	// get the refresh token cookie
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

}
