package userapi

import (
	"media-equipment-tracker/internal/application/userservice"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/handlers/userapi/dto"
	"media-equipment-tracker/internal/utils/ginerror"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	service userservice.UserService
}

func NewUserRouter(router *gin.RouterGroup, service userservice.UserService) UserRouter {
	r := UserRouter{
		service: service,
	}
	gr := router.Group("/users")
	gr.GET("", r.GetAllUsers)
	gr.GET("/:id", r.GetUser)
	gr.PATCH("/:id", r.UpdateUser)
	gr.PATCH("/:id/nice", r.UpdateNice)

	gr.GET("/me", r.Me)
	gr.GET("/me/invocations", r.MyInvocations)
	gr.PATCH("/me", r.UpdateSelf)

	return r
}

func (r *UserRouter) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := r.service.GetAll(ctx)

	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		case errs.IsValidationError(err):
			c.JSON(400, ginerror.ErrJSONBody(err))
		case errs.IsRoleAuthError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeGetAllUsersResponse(c, users))
}

func (r *UserRouter) GetUser(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeGetUserRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	user, err := r.service.Get(ctx, id)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		case errs.IsValidationError(err):
			c.JSON(400, ginerror.ErrJSONBody(err))
		case errs.IsRoleAuthError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeGetUserResponse(c, user))
}

func (r *UserRouter) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeUpdateUserRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("UpdateUserRequest", err.Error())))
		return
	}

	user, err := r.service.Update(ctx, req)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		case errs.IsEntityAlreadyExistsError(err):
			c.JSON(409, ginerror.ErrJSONBody(err))
		case errs.IsValidationError(err):
			c.JSON(400, ginerror.ErrJSONBody(err))
		case errs.IsRoleAuthError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeUpdateUserResponse(c, user))
}

func (r *UserRouter) UpdateNice(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeUpdateNiceRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("UpdateNiceRequest", err.Error())))
		return
	}

	user, err := r.service.UpdateNice(ctx, req)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		case errs.IsValidationError(err):
			c.JSON(400, ginerror.ErrJSONBody(err))
		case errs.IsRoleAuthError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeUpdateNiceResponse(c, user))
}

func (r *UserRouter) Me(c *gin.Context) {
	ctx := c.Request.Context()

	user, err := r.service.Me(ctx)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		case errs.IsValidationError(err):
			c.JSON(400, ginerror.ErrJSONBody(err))
		case errs.IsRoleAuthError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeMeUserResponse(c, user))
}

func (r *UserRouter) UpdateSelf(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeUpdateSelfRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("UpdateUserRequest", err.Error())))
		return
	}

	user, err := r.service.UpdateSelf(ctx, req)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		case errs.IsValidationError(err):
			c.JSON(400, ginerror.ErrJSONBody(err))
		case errs.IsRoleAuthError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeUpdateSelfResponse(c, user))
}

func (r *UserRouter) MyInvocations(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeMyInvocationsRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("MyInvocationsRequest", err.Error())))
		return
	}

	invs, err := r.service.Invocations(ctx, req)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		case errs.IsValidationError(err):
			c.JSON(400, ginerror.ErrJSONBody(err))
		case errs.IsRoleAuthError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeMyInvocationsResponse(c, invs))
}
