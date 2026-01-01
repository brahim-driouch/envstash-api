package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type TokenSub struct {
	Id         string `json:"id"`
	Fullname   string `json:"fullname"`
	Email      string `json:"email"`
	IsVerified bool   `json:"isVerified"`
	IsAdmin    bool   `json:"isAdmin"`
}
type JWTClaims struct {
	UserSub TokenSub
	jwt.RegisteredClaims
}
type EmailVerificationClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateEmailVerifcationToken(userId string) (string, error) {
	if userId == "" {
		return "", fmt.Errorf("user ID cannot be empty")
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})
	tokenString, err := claims.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateAccessToken(userSub TokenSub, minutes int32) (string, error) {
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

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(bytes)
	return token, nil
}

func VerifyAndDecodeEmailVerificationToken(ctx context.Context, token string) (*EmailVerificationClaims, error) {
	// Parse the token to check if it's valid
	parsedToken, err := jwt.ParseWithClaims(token, &EmailVerificationClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}
	claims, ok := parsedToken.Claims.(*EmailVerificationClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims type")
	}

	// Check if token is expired- token is 24 hours
	if claims.RegisteredClaims.ExpiresAt != nil && claims.RegisteredClaims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	// Token is valid
	return claims, nil
}
