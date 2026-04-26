package userapi

import (
	"errors"
	"media-equipment-tracker/internal/application/userservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/handlers/userapi/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	service userservice.UserService
}

func NewUserRouter(router *gin.RouterGroup, service userservice.UserService) UserRouter {
	r := UserRouter{
		service: service,
	}
	gr := router.Group("users")
	gr.GET("/me", r.GetMe)
	gr.PATCH("/me", r.UpdateMe)
	//gr.PATCH("/me/invocations", r.GetMyInvocations) // TODO: нужен invocationService
	//gr.POST("/logout", r.Logout)
	return r
}

func (r *UserRouter) GetMe(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := r.service.GetCurrent(ctx)
	if err != nil {
		if errors.Is(err, errs.EntityNotFoundError{}) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, dto.UserToUserResponse(user))
}

func (r *UserRouter) UpdateMe(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.UpdateUserDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := r.service.GetCurrent(ctx)
	if err != nil {
		if errors.Is(err, errs.EntityNotFoundError{}) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		return
	}
	user.FullName = req.FullName
	user.Email = req.Email
	updatedUser, err := r.service.Update(ctx, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.UserToUserResponse(updatedUser))
}

// TODO: нужен invocationService
func (r *UserRouter) GetMyInvocations(c *gin.Context) {
	ctx := c.Request.Context()
	var filter dto.InvocationFilterDto
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users, err := r.service.GetAll(ctx, domain.UserWithEquipmentInvocations(), domain.UserWithStudioInvocations())
	if err != nil {
		if errors.Is(err, errs.EntityNotFoundError{}) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		return
	}
	userInvocationsResponses := make([]dto.UserInvocationsResponse, 0)
	for _, user := range users {
		userInvocationsResponses = append(userInvocationsResponses, dto.UserToUserInvocationsResponse(user))
	}
	c.JSON(http.StatusOK, userInvocationsResponses)
}
