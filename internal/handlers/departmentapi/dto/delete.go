package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeleteDepartmentRequest struct {
	ID string `uri:"id" binding:"required"`
}

func DeserializeDeleteDepartmentRequest(c *gin.Context) (uuid.UUID, error) {
	var req DeleteDepartmentRequest
	if err := c.ShouldBindUri(&req); err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(req.ID)
}
