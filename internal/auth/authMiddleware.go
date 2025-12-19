package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// HashPassword hashes a plaintext password using bcrypt

// CheckPassword compares a plaintext password with a hash

// jwt

// Generate a JWT token for a user

// Verify JWT token and return claims

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
	claims, err := VerifyToken(token)
	if err != nil {
		// get the refresh token cookie
		refreshToken, refreshTokenErr := c.Cookie("nvbx_ref_token")
		if refreshTokenErr != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session, please try to login"})
			return
		}
		// check the refresh token
		rerfeshTokenClaims, invalidRefreshTokenErr := VerifyToken(refreshToken)
		if invalidRefreshTokenErr != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session, please try to login"})
			return
		}
		//generate new access token
		newAccessToken, newAccessTokenErr := GenerateToken(rerfeshTokenClaims.UserSub, 30)
		if newAccessTokenErr != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not create a new session token, please try again later or report the error"})
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "a new access token generated", "data": gin.H{"accessToken": newAccessToken}})
		return
	}
	c.Set("user", claims.UserSub)
	c.Set("userId", claims.UserSub.Id)
	c.Next()

}
