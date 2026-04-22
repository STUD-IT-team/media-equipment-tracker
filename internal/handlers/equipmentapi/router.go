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
	return r
}

func (r *EquipmentRouter) Search(c *gin.Context) {
	ctx := c.Request.Context()

	search, err := dto.DeserializeSearchEquipmentRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJsonBody(errs.ValidationError{Field: "query", Message: err.Error()}))
		return
	}

	items, err := r.service.Search(ctx, search)
	if err != nil {
		if errors.Is(err, errs.EntityNotFoundError{}) {
			c.JSON(404, ginerror.ErrJsonBody(err))
		} else if errors.Is(err, errs.ValidationError{}) {
			c.JSON(400, ginerror.ErrJsonBody(err))
		} else {
			c.JSON(500, ginerror.ErrJsonBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeSearchEquipmentResponse(c, items))
}
