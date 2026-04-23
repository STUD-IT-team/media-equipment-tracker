package ginerror

import "github.com/gin-gonic/gin"

func ErrJSONBody(err error) any {
	return gin.H{"error": err.Error()}
}
