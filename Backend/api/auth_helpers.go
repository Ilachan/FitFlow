package api

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func getTokenStringFromAuthHeader(c *gin.Context) (string, error) {
	auth := strings.TrimSpace(c.GetHeader("Authorization"))
	if auth == "" {
		return "", errors.New("missing authorization header")
	}
	const bearer = "Bearer "
	if !strings.HasPrefix(auth, bearer) {
		return "", errors.New("invalid authorization header")
	}
	token := strings.TrimSpace(strings.TrimPrefix(auth, bearer))
	if token == "" {
		return "", errors.New("invalid authorization header")
	}
	return token, nil
}