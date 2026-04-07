package api

import (
	"errors"
	"net/http"
	"strconv"

	"my-course-backend/service"

	"github.com/gin-gonic/gin"
)

func requireInstructorRole(c *gin.Context) (uint, error) {
	tokenString, err := getTokenStringFromAuthHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return 0, err
	}

	roleID, err := service.GetRoleIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return 0, err
	}

	if roleID != 4 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: instructor role required"})
		return 0, errors.New("forbidden")
	}

	userID, err := service.ExtractUserIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return 0, err
	}

	return userID, nil
}

func InstructorListCourses(c *gin.Context) {
	instructorID, err := requireInstructorRole(c)
	if err != nil {
		return
	}

	courses, err := service.ListInstructorCourses(instructorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func InstructorListCourseEnrollments(c *gin.Context) {
	instructorID, err := requireInstructorRole(c)
	if err != nil {
		return
	}

	courseIDStr := c.Param("id")
	courseID64, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	enrollments, err := service.ListInstructorCourseEnrollments(instructorID, uint(courseID64))
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if err.Error() == "class not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "class not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"enrollments": enrollments})
}

type InstructorUpdateStatusInput struct {
	UserID uint   `json:"user_id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

func InstructorUpdateEnrollmentStatus(c *gin.Context) {
	instructorID, err := requireInstructorRole(c)
	if err != nil {
		return
	}

	courseIDStr := c.Param("id")
	courseID64, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	var input InstructorUpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.UpdateEnrollmentStatusByInstructor(instructorID, uint(courseID64), input.UserID, input.Status); err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if err.Error() == "enrollment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "enrollment not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}