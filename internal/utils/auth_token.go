package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"media-equipment-tracker/cmd/app/config"
)

func TokenFromHeader(c *gin.Context) (string, error) {
	authorizationHeader := c.GetHeader(config.AuthorizationHeaderKey)
	if authorizationHeader == "" {
		return "", errors.New("authorization header is not provided")
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		return "", errors.New("invalid authorization header format")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != config.AuthorizationTypeBearer {
		return "", fmt.Errorf("unsupported authorization type %s", authorizationType)
	}
	accessToken := fields[1]
	return accessToken, nil
}
