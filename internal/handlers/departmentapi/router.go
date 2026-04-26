package departmentapi

import (
	"errors"

	"media-equipment-tracker/internal/application/departmentservice"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/handlers/departmentapi/dto"
	"media-equipment-tracker/internal/utils/ginerror"

	"github.com/gin-gonic/gin"
)

type DepartmentRouter struct {
	service departmentservice.DepartmentService
}

func NewRouter(router *gin.RouterGroup, service departmentservice.DepartmentService) DepartmentRouter {
	r := DepartmentRouter{
		service: service,
	}
	gr := router.Group("departments")
	gr.GET("", r.List)
	gr.POST("", r.Create)
	gr.GET("/:id", r.Get)
	gr.PATCH("/:id", r.Update)
	gr.DELETE("/:id", r.Delete)
	gr.GET("/:id/users", r.GetUsers)
	gr.GET("/:id/equipment", r.GetEquipment)
	return r
}

func (r *DepartmentRouter) List(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeListDepartmentsRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("query", err.Error())))
		return
	}

	departments, err := r.service.ListDepartments(ctx, req.Search)
	if err != nil {
		c.JSON(500, ginerror.ErrJSONBody(err))
		return
	}

	c.JSON(200, dto.SerializeListDepartmentsResponse(c, departments))
}

func (r *DepartmentRouter) Create(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := dto.DeserializeCreateDepartmentRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("CreateDepartmentRequest", err.Error())))
		return
	}

	department, err := r.service.CreateDepartment(ctx, &req)
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

	c.JSON(201, dto.SerializeCreateDepartmentResponse(c, department))
}

func (r *DepartmentRouter) Get(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeGetDepartmentRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	department, err := r.service.GetDepartment(ctx, id)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeGetDepartmentResponse(c, department))
}

func (r *DepartmentRouter) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, req, err := dto.DeserializeUpdateDepartmentRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("UpdateDepartmentRequest", err.Error())))
		return
	}

	department, err := r.service.UpdateDepartment(ctx, id, &req)
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

	c.JSON(200, dto.SerializeUpdateDepartmentResponse(c, department))
}

func (r *DepartmentRouter) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := dto.DeserializeDeleteDepartmentRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("id", err.Error())))
		return
	}

	err = r.service.DeleteDepartment(ctx, id)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		case errors.Is(err, departmentservice.ErrDepartmentHasUsers):
			c.JSON(409, ginerror.ErrJSONBody(err))
		case errors.Is(err, departmentservice.ErrDepartmentHasInvocations):
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

func (r *DepartmentRouter) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()

	id, role, err := dto.DeserializeGetDepartmentUsersRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("GetDepartmentUsersRequest", err.Error())))
		return
	}

	users, err := r.service.GetDepartmentUsers(ctx, id, role)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeGetDepartmentUsersResponse(c, users))
}

func (r *DepartmentRouter) GetEquipment(c *gin.Context) {
	ctx := c.Request.Context()

	id, availableOnly, err := dto.DeserializeGetDepartmentEquipmentRequest(c)
	if err != nil {
		c.JSON(400, ginerror.ErrJSONBody(errs.NewValidationError("GetDepartmentEquipmentRequest", err.Error())))
		return
	}

	equipment, err := r.service.GetDepartmentEquipment(ctx, id, availableOnly)
	if err != nil {
		switch {
		case errs.IsEntityNotFoundError(err):
			c.JSON(404, ginerror.ErrJSONBody(err))
		default:
			c.JSON(500, ginerror.ErrJSONBody(err))
		}
		return
	}

	c.JSON(200, dto.SerializeGetDepartmentEquipmentResponse(c, equipment))
}
