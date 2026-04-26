package organizationapi

import (
	"errors"

	"media-equipment-tracker/internal/application/organizationservice"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/handlers/organizationapi/dto"
	"media-equipment-tracker/internal/utils/ginerror"

	"github.com/gin-gonic/gin"
)

type OrganizationRouter struct {
	service organizationservice.OrganizationService
}

func NewRouter(router *gin.RouterGroup, service organizationservice.OrganizationService) OrganizationRouter {
	r := OrganizationRouter{
		service: service,
	}
	gr := router.Group("organizations")
	gr.GET("", r.List)
	gr.POST("", r.Create)
	gr.GET("/:id", r.Get)
	gr.PATCH("/:id", r.Update)
	gr.DELETE("/:id", r.Delete)
	gr.GET("/:id/users", r.GetUsers)
	return r
}

func (r *OrganizationRouter) List(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeListOrganizationsRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("query", err.Error())))
		return
	}

	organizations, err := r.service.ListOrganizations(ctx, req.Search)
	if err != nil {
		c.JSON(500, ginerror.ErrJSONBody(err))
		return
	}

	c.JSON(200, dto.SerializeListOrganizationsResponse(c, organizations))
}

func (r *OrganizationRouter) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeCreateOrganizationRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("CreateOrganizationRequest", err.Error())))
		return
	}

	organization, err := r.service.CreateOrganization(ctx, &req)
	if err != nil {
		switch {
		case errs.IsValidationError(err):
			c.JSON(400, ginerror.ErrJSONBody(err))
		case errs.IsRoleAuthError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(201, dto.SerializeCreateOrganizationResponse(c, organization))
}

func (r *OrganizationRouter) Get(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeGetOrganizationRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	organization, err := r.service.GetOrganization(ctx, id)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeGetOrganizationResponse(c, organization))
}

func (r *OrganizationRouter) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, req, err := dto.DeserializeUpdateOrganizationRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("UpdateOrganizationRequest", err.Error())))
		return
	}

	organization, err := r.service.UpdateOrganization(ctx, id, &req)
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

	c.JSON(200, dto.SerializeUpdateOrganizationResponse(c, organization))
}

func (r *OrganizationRouter) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeDeleteOrganizationRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.service.DeleteOrganization(ctx, id)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		case errors.Is(err, organizationservice.ErrOrganizationHasUsers):
			c.JSON(409, ginerror.ErrJSONBody(err))
		case errors.Is(err, organizationservice.ErrOrganizationHasInvocations):
			c.JSON(409, ginerror.ErrJSONBody(err))
		case errs.IsRoleAuthError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(204, nil)
}

func (r *OrganizationRouter) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeGetOrganizationUsersRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("GetOrganizationUsersRequest", err.Error())))
		return
	}

	users, err := r.service.GetOrganizationUsers(ctx, id)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeGetOrganizationUsersResponse(c, users))
}
