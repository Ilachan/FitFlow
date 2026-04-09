package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"my-course-backend/model"
	"my-course-backend/service"

	"github.com/gin-gonic/gin"
)

// RegisterClass enrolls the authenticated user in a course.
func RegisterClass(c *gin.Context) {
	var input model.EnrollmentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := requireRegisterPermission(c)
	if err != nil {
		return
	}

	if err := service.RegisterClass(userID, input.CourseID); err != nil {
		switch err.Error() {
		case "user not found", "class not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "enrollment already exists", "class is full", "class schedule overlaps with an existing enrolled class":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Class enrolled successfully"})
}

// DropClass removes the authenticated user's enrollment from a course.
func DropClass(c *gin.Context) {
	var input model.EnrollmentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := requireRegisterPermission(c)
	if err != nil {
		return
	}

	if err := service.DropClass(userID, input.CourseID); err != nil {
		if err.Error() == "enrollment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Class unenrolled successfully"})
}

// ListClasses returns paginated courses. Public endpoint.
// GET /classes?page=1
func ListClasses(c *gin.Context) {
	const pageSize = 20

	pageStr := c.Query("page")
	page := 1
	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			page = v
		}
	}

	classes, total, err := service.ListClassesPaged(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
		"classes":   classes,
	})
}

// GetClass returns a single course by ID. Public endpoint.
func GetClass(c *gin.Context) {
	classIDStr := c.Param("id")
	classID, err := strconv.ParseUint(classIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	class, err := service.GetClass(uint(classID))
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

// ListClassEnrollments returns all enrollments for a class.
func ListClassEnrollments(c *gin.Context) {
	if err := requireManagerRole(c); err != nil {
		return
	}

	classIDStr := c.Param("id")
	classID, err := strconv.ParseUint(classIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	enrollments, err := service.ListClassEnrollments(uint(classID))
	if err != nil {
		if err.Error() == "class not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"enrollments": enrollments})
}

// GetUserEnrolledClasses returns all courses a user is enrolled in.
func GetUserEnrolledClasses(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	authenticatedUserID, err := getUserIDFromAuthHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if authenticatedUserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: you can only view your own enrollments"})
		return
	}

	courses, err := service.GetUserEnrolledClasses(uint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

// GetStudentAnalytics returns student dashboard analytics for 7d, 1m, or 3m.
func GetUserAnalytics(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	authUserID, err := getUserIDFromAuthHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if uint(userID) != authUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	rangeKey := strings.TrimSpace(c.DefaultQuery("range", "7d"))

	analytics, err := service.GetUserAnalytics(authUserID, rangeKey)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analytics": analytics})
}

func getUserIDFromAuthHeader(c *gin.Context) (uint, error) {
	authorization := strings.TrimSpace(c.GetHeader("Authorization"))
	if authorization == "" {
		return 0, errors.New("missing authorization header")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authorization, bearerPrefix) {
		return 0, errors.New("invalid authorization header")
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authorization, bearerPrefix))
	if tokenString == "" {
		return 0, errors.New("invalid authorization header")
	}

	return service.ExtractUserIDFromToken(tokenString)
}

func getRoleIDFromAuthHeader(c *gin.Context) (uint, error) {
	authorization := strings.TrimSpace(c.GetHeader("Authorization"))
	if authorization == "" {
		return 0, errors.New("missing authorization header")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authorization, bearerPrefix) {
		return 0, errors.New("invalid authorization header")
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authorization, bearerPrefix))
	if tokenString == "" {
		return 0, errors.New("invalid authorization header")
	}

	return service.GetRoleIDFromToken(tokenString)
}

// register/drop permission check
func requireRegisterPermission(c *gin.Context) (uint, error) {
	userID, err := getUserIDFromAuthHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return 0, err
	}

	roleID, err := getRoleIDFromAuthHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return 0, err
	}

	// role_id 1(Student)/2(SuperManager)/3(Manager) are allowed
	if roleID != 1 && roleID != 2 && roleID != 3 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient role permissions"})
		return 0, errors.New("forbidden")
	}

	return userID, nil
}