package ginerror

import "github.com/gin-gonic/gin"

func ErrJsonBody(err error) any {
	return gin.H{"error": err.Error()}
}
