package invocationapi

import (
	"media-equipment-tracker/internal/application/invocationservice"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/handlers/invocationapi/dto"
	"media-equipment-tracker/internal/utils/ginerror"

	"github.com/gin-gonic/gin"
)

type InvocationRouter struct {
	service invocationservice.InvocationService
}

func NewRouter(router *gin.RouterGroup, service invocationservice.InvocationService) InvocationRouter {
	r := InvocationRouter{
		service: service,
	}
	gr := router.Group("equipment-invocations")
	gr.GET("", r.Search)
	gr.POST("", r.Create)
	gr.GET("/:id", r.Get)
	gr.PATCH("/:id", r.Update)
	gr.DELETE("/:id", r.Delete)
	gr.PATCH("/:id/status", r.UpdateStatus)

	gr.POST("/:id/become-curator", r.BecomeCurator)
	gr.POST("/:id/approve", r.Approve)
	gr.POST("/:id/request-changes", r.RequestChanges)
	gr.POST("/:id/issue/:equipment_id", r.Issue)
	gr.POST("/:id/return/:equipment_id", r.Return)
	gr.POST("/:id/complete", r.Complete)
	gr.POST("/:id/cancel", r.Cancel)

	return r
}

func (r *InvocationRouter) Search(c *gin.Context) {
	ctx := c.Request.Context()

	search, err := dto.DeserializeSearchInvocationRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("query", err.Error())))
		return
	}

	items, err := r.service.Search(ctx, &search)
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

	c.JSON(200, dto.SerializeSearchInvocationResponse(c, items))
}

func (r *InvocationRouter) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeCreateInvocationRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("CreateInvocationRequest", err.Error())))
		return
	}

	inv, err := r.service.Create(ctx, req)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		case errs.IsValidationError(err):
			c.JSON(400, ginerror.ErrJSONBody(err))
		case errs.IsEntityAlreadyExistsError(err):
			c.JSON(409, ginerror.ErrJSONBody(err))
		case errs.IsRoleAuthError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		case errs.IsEquipmentAccessError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(201, dto.SerializeCreateInvocationResponse(c, inv))
}

func (r *InvocationRouter) Get(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeGetInvocationRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	inv, err := r.service.Get(ctx, id)
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

	c.JSON(200, dto.SerializeGetInvocationResponse(c, inv))
}

func (r *InvocationRouter) Update(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeUpdateInvocationRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("UpdateInvocationRequest", err.Error())))
		return
	}

	inv, err := r.service.Update(ctx, req)
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
		case errs.IsEquipmentAccessError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeUpdateInvocationResponse(c, inv))
}

func (r *InvocationRouter) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeUpdateEquipmentStatusRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("UpdateEquipmentStatusRequest", err.Error())))
		return
	}

	inv, err := r.service.UpdateStatus(ctx, req)
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

	c.JSON(200, dto.SerializeUpdateInvocationStatusResponse(c, inv))
}

func (r *InvocationRouter) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeDeleteInvocationRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.service.Delete(ctx, id)
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

	c.JSON(204, nil)
}

func (r *InvocationRouter) BecomeCurator(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeStatusWorkflowRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.service.Become(ctx, id)
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

	c.JSON(200, nil)
}

func (r *InvocationRouter) Approve(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeStatusWorkflowRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.service.Approve(ctx, id)
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

	c.JSON(200, nil)
}

func (r *InvocationRouter) RequestChanges(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeStatusWorkflowRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.service.RequestChanges(ctx, id)
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

	c.JSON(200, nil)
}

func (r *InvocationRouter) Issue(c *gin.Context) {
	ctx := c.Request.Context()

	id, equipmentID, err := dto.DeserializeStatusEquipmentRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.service.Issue(ctx, id, equipmentID)
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

	c.JSON(200, nil)
}

func (r *InvocationRouter) Return(c *gin.Context) {
	ctx := c.Request.Context()

	id, equipmentID, err := dto.DeserializeStatusEquipmentRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.service.Return(ctx, id, equipmentID)
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

	c.JSON(200, nil)
}

func (r *InvocationRouter) Complete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeStatusWorkflowRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.service.Complete(ctx, id)
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

	c.JSON(200, nil)
}

func (r *InvocationRouter) Cancel(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeStatusWorkflowRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.service.Cancel(ctx, id)
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

	c.JSON(200, nil)
}
