package api

import (
	"errors"
	"net/http"
	"strconv"
	"my-course-backend/model"
	"my-course-backend/service"

	"github.com/gin-gonic/gin"
)

// manager role required: role_id == 2 or 3
func requireManagerRole(c *gin.Context) error {
	tokenString, err := getTokenStringFromAuthHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return err
	}

	roleID, err := service.GetRoleIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return err
	}

	if roleID != 2 && roleID != 3 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: manager role required"})
		return errors.New("forbidden")
	}
	return nil
}

// POST /classes (manager only)
func ManagerCreateClass(c *gin.Context) {
	if err := requireManagerRole(c); err != nil {
		return
	}

	var input service.CourseUpsertInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	class, err := service.ManagerCreateCourse(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"class": class})
}

// PUT /classes/:id (manager only)
func ManagerUpdateClass(c *gin.Context) {
	if err := requireManagerRole(c); err != nil {
		return
	}

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	var input service.CourseUpsertInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	class, err := service.ManagerUpdateCourse(uint(id64), input)
	if err != nil {
		if err.Error() == "class not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"class": class})
}

// DELETE /classes/:id (manager only)
func ManagerDeleteClass(c *gin.Context) {
	if err := requireManagerRole(c); err != nil {
		return
	}

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	if err := service.ManagerDeleteCourse(uint(id64)); err != nil {
		if err.Error() == "class not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Class deleted successfully"})
}

// ManagerRegister handles POST /auth/manager/register
func ManagerRegister(c *gin.Context) {
	var input model.ManagerRegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}

	if err := service.RegisterManager(input); err != nil {
		// Map common errors to status codes
		switch err.Error() {
		case "email already exists":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case "invalid invite code":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "invite code is not active", "invite code already used", "invite code expired", "invite code not allowed for this email":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Manager registration successful"})
}

// ✅ GET /manager/users?page=1&limit=20
func ManagerListUsers(c *gin.Context) {
	if err := requireManagerRole(c); err != nil {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, total, page, limit, totalPages, err := service.ManagerListUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
	})
}

// ✅ GET /manager/users/:id/enrollments
func ManagerListUserEnrollments(c *gin.Context) {
	if err := requireManagerRole(c); err != nil {
		return
	}

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	enrollments, err := service.ManagerListUserEnrollments(uint(id64))
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"enrollments": enrollments})
}

type ManagerAddEnrollmentInput struct {
	CourseID uint `json:"course_id" binding:"required"`
}

// ✅ POST /manager/users/:id/enrollments
func ManagerAddUserEnrollment(c *gin.Context) {
	if err := requireManagerRole(c); err != nil {
		return
	}

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input ManagerAddEnrollmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.ManagerAddUserEnrollment(uint(id64), input.CourseID); err != nil {
		switch err.Error() {
		case "user not found", "class not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "enrollment already exists":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "enrollment created"})
}

// ✅ DELETE /manager/users/:id/enrollments/:course_id
func ManagerDeleteUserEnrollment(c *gin.Context) {
	if err := requireManagerRole(c); err != nil {
		return
	}

	userIDStr := c.Param("id")
	userID64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	courseIDStr := c.Param("course_id")
	courseID64, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	if err := service.ManagerDeleteUserEnrollment(uint(userID64), uint(courseID64)); err != nil {
		if err.Error() == "enrollment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enrollment deleted"})
}

