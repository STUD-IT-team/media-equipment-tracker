package utils

import (
	"errors"
	"fmt"
	"media-equipment-tracker/cmd/app/config"
	"strings"

	"github.com/gin-gonic/gin"
)

func TokenFromHeader(c *gin.Context) (string, error) {
	authorizationHeader := c.GetHeader(config.AuthorizationHeaderKey)
	if len(authorizationHeader) == 0 {
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
