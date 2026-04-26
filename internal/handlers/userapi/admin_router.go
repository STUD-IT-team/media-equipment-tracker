package userapi

import (
	"errors"
	"media-equipment-tracker/internal/application/userservice"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/handlers/userapi/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRouterForAdmin struct {
	service userservice.UserService
}

func NewUserRouterForAdmin(router *gin.RouterGroup, service userservice.UserService) UserRouterForAdmin {
	r := UserRouterForAdmin{
		service: service,
	}
	gr := router.Group("users")
	gr.GET("", r.GetAll)
	gr.GET("/:id", r.GetByID)
	gr.PATCH("/:id", r.Update)
	return r
}

func (r *UserRouterForAdmin) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	var filter domain.UserFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users, err := r.service.GetAll(ctx, domain.UserWithFilter(filter))
	if err != nil {
		if errors.Is(err, errs.EntityNotFoundError{}) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		return
	}
	userResponses := make([]dto.UserResponse, 0)
	for _, user := range users {
		userResponses = append(userResponses, dto.UserToUserResponse(user))
	}
	c.JSON(http.StatusOK, userResponses)
}

func (r *UserRouterForAdmin) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id format"})
		return
	}
	user, err := r.service.GetByID(ctx, userID,
		domain.UserWithOrganizations(), domain.UserWithDepartments(), domain.UserWithStudioInvocations(), domain.UserWithEquipmentInvocations())
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

func (r *UserRouterForAdmin) Update(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id format"})
		return
	}
	var userDto dto.UpdateUserByAdminDto
	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := dto.UpdateUserByAdminDtoToUser(userDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = userID
	if _, err := r.service.Update(ctx, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
