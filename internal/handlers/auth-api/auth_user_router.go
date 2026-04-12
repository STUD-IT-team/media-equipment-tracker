package authapi

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	authuser "media-equipment-tracker/internal/application/auth_service/auth_user"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils"
)

type AuthUserRouter struct {
	service authuser.AuthUserService
}

func NewAuthUserRouter(router *gin.RouterGroup, service authuser.AuthUserService) AuthUserRouter {
	r := AuthUserRouter{
		service: service,
	}
	gr := router.Group("auth")
	gr.POST("/register", r.Register)
	gr.POST("/login", r.Login)
	gr.POST("/logout", r.Logout)
	return r
}

func (r *AuthUserRouter) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req authuser.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := r.service.RegisterUser(ctx, req); err != nil {
		if errors.Is(err, errs.EntityNotFoundError{}) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (r *AuthUserRouter) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req authuser.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accessToken, err := r.service.LoginUser(ctx, req)
	if err != nil {
		if errors.Is(err, errs.EntityNotFoundError{}) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	rsp := authuser.LoginUserResponse{
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, rsp)
}

func (r *AuthUserRouter) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	accessToken, err := utils.TokenFromHeader(c)
	if err != nil {
		return
	}
	r.service.LogoutUser(ctx, accessToken)
	c.JSON(http.StatusOK, gin.H{})
}
