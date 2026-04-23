package equipmentapi

import (
	"errors"

	"media-equipment-tracker/internal/application/equipmentservice"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/handlers/equipmentapi/dto"
	"media-equipment-tracker/internal/utils/ginerror"

	"github.com/gin-gonic/gin"
)

type EquipmentRouter struct {
	service equipmentservice.EquipmentService
}

func NewRouter(router *gin.RouterGroup, service equipmentservice.EquipmentService) EquipmentRouter {
	r := EquipmentRouter{
		service: service,
	}
	gr := router.Group("equipment")
	gr.GET("/list", r.Search)
	gr.POST("", r.Create)
	gr.PATCH("/:id", r.Update)
	gr.DELETE("/:id", r.Delete)
	gr.GET("/:id", r.Get)
	gr.GET("/:id/availability", r.Availability)
	gr.GET("/inventory/:inv", r.GetByInventoryNumber)
	return r
}

func (r *EquipmentRouter) Search(c *gin.Context) {
	ctx := c.Request.Context()

	search, err := dto.DeserializeSearchEquipmentRequest(c)
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

	c.JSON(200, dto.SerializeSearchEquipmentResponse(c, items))
}

func (r *EquipmentRouter) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeCreateEquipmentRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("CreateEquipmentRequest", err.Error())))
		return
	}

	equipment, err := r.service.CreateEquipment(ctx, &req)
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

	c.JSON(201, dto.SerializeCreateEquipmentResponse(c, equipment))
}

func (r *EquipmentRouter) Update(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeUpdateEquipmentRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("UpdateEquipmentRequest", err.Error())))
		return
	}

	equipment, err := r.service.UpdateEquipment(ctx, &req)
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

	c.JSON(200, dto.SerializeUpdateEquipmentResponse(c, equipment))
}

func (r *EquipmentRouter) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeDeleteEquipmentRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.service.Delete(ctx, id)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		case errors.Is(err, equipmentservice.ErrEquipmentHasInvocations):
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

	c.JSON(204, nil)
}

func (r *EquipmentRouter) Get(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeGetEquipmentRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	equipment, err := r.service.Get(ctx, id)
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

	c.JSON(200, dto.SerializeGetEquipmentResponse(c, equipment))
}

func (r *EquipmentRouter) GetByInventoryNumber(c *gin.Context) {
	ctx := c.Request.Context()

	inv, err := dto.DeserializeGetEquipmentByInventoryNumberRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("inv", err.Error())))
		return
	}

	equipment, err := r.service.GetByInventoryNumber(ctx, inv)
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

	c.JSON(200, dto.SerializeGetEquipmentResponse(c, equipment))
}

func (r *EquipmentRouter) Availability(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeAvailabilityEquipmentRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("AvailabilityEquipmentRequest", err.Error())))
		return
	}

	resp, err := r.service.Availability(ctx, req)
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

	c.JSON(200, dto.SerializeAvailabilityEquipmentResponse(c, *resp))
}
