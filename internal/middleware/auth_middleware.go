package middleware

import (
	"net/http"

	"media-equipment-tracker/internal/adapters/inmem"
	authzservice "media-equipment-tracker/internal/application/authz_service"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/utils"

	"github.com/gin-gonic/gin"
)

type TokenVerifier interface {
	VerifyByToken(tokenStr string, needRoles []domain.RoleAuth) (*domain.TokenPayload, error)
}

func AuthMiddleware(authZ authzservice.AuthZ, tokenRep inmem.TokenRepository, authServ TokenVerifier, needRoles []domain.RoleAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := utils.TokenFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

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
