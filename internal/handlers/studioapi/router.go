package studioapi

import (
	"media-equipment-tracker/internal/application/studioservice"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/handlers/studioapi/dto"
	"media-equipment-tracker/internal/utils/ginerror"

	"github.com/gin-gonic/gin"
)

type StudioRouter struct {
	studioService studioservice.StudioService
}

func NewRouter(router *gin.RouterGroup, studioService studioservice.StudioService) *StudioRouter {
	r := &StudioRouter{
		studioService: studioService,
	}
	gr := router.Group("studio-invocations")
	gr.GET("", r.Search)
	gr.GET("/schedule", r.Schedule)
	gr.POST("", r.Create)
	gr.GET("/:id", r.Get)
	gr.PATCH("/:id", r.Update)
	gr.PATCH("/:id/status", r.UpdateStatus)
	gr.DELETE("/:id", r.Delete)

	gr.POST("/:id/become-curator", r.BecomeCurator)
	gr.POST("/:id/approve", r.Approve)
	gr.POST("/:id/request-changes", r.RequestChanges)
	gr.POST("/:id/complete", r.Complete)
	gr.POST("/:id/cancel", r.Cancel)

	return r
}

func (r *StudioRouter) Search(c *gin.Context) {
	ctx := c.Request.Context()

	search, err := dto.DeserializeSearchStudioRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("query", err.Error())))
		return
	}

	items, err := r.studioService.Search(ctx, search)
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

	c.JSON(200, dto.SerializeSearchStudioResponse(c, items))
}

func (r *StudioRouter) Schedule(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeScheduleStudioRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("query", err.Error())))
		return
	}

	items, err := r.studioService.Schedule(ctx, req)
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

	c.JSON(200, dto.SerializeScheduleStudioResponse(c, items))
}

func (r *StudioRouter) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeCreateStudioRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("CreateStudioRequest", err.Error())))
		return
	}

	inv, err := r.studioService.Create(ctx, req)
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
		case errs.IsStudioAccessError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(201, dto.SerializeCreateStudioResponse(c, inv))
}

func (r *StudioRouter) Update(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeUpdateStudioRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("UpdateStudioRequest", err.Error())))
		return
	}

	inv, err := r.studioService.Update(ctx, req)
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
		case errs.IsStudioAccessError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeUpdateStudioResponse(c, inv))
}

func (r *StudioRouter) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeUpdateStudioStatusRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("UpdateStudioStatusRequest", err.Error())))
		return
	}

	inv, err := r.studioService.UpdateStatus(ctx, req)
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

	c.JSON(200, dto.SerializeUpdateStudioResponse(c, inv))
}

func (r *StudioRouter) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeDeleteStudioRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.studioService.Delete(ctx, id)
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

func (r *StudioRouter) Get(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeGetStudioRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	inv, err := r.studioService.Get(ctx, id)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		case errs.IsRoleAuthError(err):
			c.JSON(403, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeGetStudioResponse(c, inv))
}

func (r *StudioRouter) BecomeCurator(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeStatusWorkflowRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.studioService.Become(ctx, id)
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
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, nil)
}

func (r *StudioRouter) Approve(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeStatusWorkflowRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.studioService.Approve(ctx, id)
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

func (r *StudioRouter) RequestChanges(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeStatusWorkflowRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.studioService.RequestChanges(ctx, id)
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

func (r *StudioRouter) Complete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeStatusWorkflowRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.studioService.Complete(ctx, id)
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

func (r *StudioRouter) Cancel(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeStatusWorkflowRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.studioService.Cancel(ctx, id)
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
