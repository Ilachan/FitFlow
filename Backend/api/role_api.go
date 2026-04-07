package api

import (
    "errors"
    "net/http"

    "my-course-backend/service"

    "github.com/gin-gonic/gin"
)

type AssignRoleInput struct {
    UserID  uint   `json:"user_id" binding:"required"`
    RoleName string `json:"role_name" binding:"required"`
}

func requireSuperManagerRole(c *gin.Context) error {
    tokenString, err := getTokenStringFromAuthHeader(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return err
    }

    roleID, err := service.GetRoleIDFromToken(tokenString)
    if err != nil || roleID != 2 {
        c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: SuperManager role required"})
        return errors.New("forbidden")
    }
    return nil
}

func AssignUserRole(c *gin.Context) {
    if err := requireSuperManagerRole(c); err != nil {
        return
    }

    var input AssignRoleInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := service.AssignUserRole(input.UserID, input.RoleName); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "role updated"})
}