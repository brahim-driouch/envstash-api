package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func getAuthorizationHeader(c *gin.Context) string {
	tokenHeader := c.GetHeader("Authorization")

	if tokenHeader == "" {
		return ""
	}

	parts := strings.SplitN(tokenHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func AuthMiddleware(c *gin.Context) {
	token := getAuthorizationHeader(c)

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
