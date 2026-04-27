package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeserializeDeleteStudioRequest(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
