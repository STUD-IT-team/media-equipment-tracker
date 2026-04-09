package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthRouter struct {
}

func NewHealthRouter(router *gin.RouterGroup) HealthRouter {
	r := HealthRouter{}
	router.GET("/health", r.Health)
	return r
}

func (r *HealthRouter) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
