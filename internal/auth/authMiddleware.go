package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func getAuthorizationHeader(c *gin.Context) (string, error) {
	tokenHeader := c.GetHeader("Authorization")

	if tokenHeader == "" {
		return "", errors.New("authorization header required")
	}

	parts := strings.SplitN(tokenHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization format")
	}

	return parts[1], nil
}

func AuthMiddleware(c *gin.Context) {
	token, err := getAuthorizationHeader(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// check if token is valid
	claims, err := VerifyToken(token)
	if err != nil {
		// get the refresh token cookie
		refreshToken, refreshTokenErr := c.Cookie("nvstash_ref_token")
		if refreshTokenErr != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session, please try to login"})
			return
		}

		// check the refresh token
		refreshTokenClaims, invalidRefreshTokenErr := VerifyToken(refreshToken)
		if invalidRefreshTokenErr != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session, please try to login"})
			return
		}

		// generate new access token
		newAccessToken, newAccessTokenErr := GenerateToken(refreshTokenClaims.UserSub, 15)
		if newAccessTokenErr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not create a new session token"})
			return
		}

		// Set new token in header
		c.Header("X-New-Access-Token", newAccessToken)

		c.Set("user", refreshTokenClaims.UserSub)
		c.Set("userId", refreshTokenClaims.UserSub.Id)
		c.Next()
		return
	}

	// Access token is valid
	c.Set("user", claims.UserSub)
	c.Set("userId", claims.UserSub.Id)
	c.Next()
}
