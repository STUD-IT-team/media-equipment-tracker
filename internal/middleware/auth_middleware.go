package middleware

import (
	"media-equipment-tracker/internal/adapters/in_mem"
	"media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	tokenmaker "media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenVerifier interface {
	VerifyByToken(tokenStr string, needRoles []domain.RoleAuth) (*tokenmaker.TokenPayload, error)
}

func AuthMiddleware(authServ TokenVerifier, needRoles []domain.RoleAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		authZ := authz_service.GetAuthZ()
		accessToken, err := utils.TokenFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		tokenRep := in_mem.GetTokenRepository()
		if !tokenRep.Check(accessToken) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		payload, err := authServ.VerifyByToken(accessToken, needRoles)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		ctx = authZ.Authorize(ctx, *payload)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
