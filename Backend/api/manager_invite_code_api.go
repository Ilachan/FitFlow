package api

import (
	"net/http"

	"my-course-backend/model"
	"my-course-backend/service"

	"github.com/gin-gonic/gin"
)

// CreateManagerInviteCode handles POST /auth/manager/invite-codes (SuperManager only)
func CreateManagerInviteCode(c *gin.Context) {
	tokenString, err := getTokenStringFromAuthHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, err := service.GetStudentIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	roleID, err := service.GetRoleIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	// Only role_id=2 (SuperManager) can create invite codes
	if roleID != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: SuperManager role required"})
		return
	}

	var input model.CreateManagerInviteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}

	code, err := service.CreateManagerInviteCode(userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Invite code created",
		"code":    code,
	})
}