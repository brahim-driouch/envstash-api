package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/brahim-driouch/envstash.git/internal/services"
	"github.com/brahim-driouch/envstash.git/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetAuthorizationHeader(c *gin.Context) string {
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

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// get access token
		token := GetAuthorizationHeader(c)

		// check if access token is valid
		claims, err := utils.VerifyToken(token)
		if err != nil {
			// if not valid, get the refresh token cookie
			refreshToken, refreshTokenErr := c.Cookie("nvstash_ref_token")
			// if refresh token is not valid, return unauthorized
			if refreshTokenErr != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session, please try to login"})
				return
			}

			// check the refresh token
			refreshTokenClaims, invalidRefreshTokenErr := authService.FindRefreshToken(ctx, refreshToken)
			// if refresh token is not valid, return unauthorized
			if invalidRefreshTokenErr != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session, please try to login"})
				return
			}
			//check if refresh token is expired
			if refreshTokenClaims.ExpiresAt.Before(time.Now()) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session expired, please login"})
				return
			}
			//check if refresh token is revoked
			if refreshTokenClaims.RevokedAt != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session revoked, please login"})
				return
			}
			// get user from db
			user, err := authService.FindUserByID(ctx, refreshTokenClaims.UserID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not find user"})
				return
			}
			// create new access token sub
			newAccessTokenSub := utils.TokenSub{
				Id:         user.ID,
				Fullname:   user.Fullname,
				Email:      user.Email,
				IsVerified: user.IsVerified,
				IsAdmin:    user.IsAdmin,
			}
			// generate new access token
			newAccessToken, newAccessTokenErr := utils.GenerateAccessToken(newAccessTokenSub, 15)
			// if new access token is not generated, return internal server error
			if newAccessTokenErr != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not create a new session token"})
				return
			}

			// Set new token in header
			c.Header("X-New-Access-Token", newAccessToken)

			c.Set("user", newAccessTokenSub)
			c.Set("userId", newAccessTokenSub.Id)
			c.Next()
			return
		}

		// Access token is valid
		c.Set("user", claims.UserSub)
		c.Set("userId", claims.UserSub.Id)
		c.Next()
	}
}
