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

// RegisterClass registers a student for a course.
// Auth required, and role_id must be 1/2/3.
func RegisterClass(c *gin.Context) {
	var input model.StudentEnrollmentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	studentID, err := requireRegisterPermission(c)
	if err != nil {
		return
	}

	if err := service.RegisterClass(studentID, input.CourseID); err != nil {
		switch err.Error() {
		case "student not found", "class not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "registration already exists", "class is full":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Class registered successfully"})
}

// DropClass removes a student from a course.
// Auth required, and role_id must be 1/2/3.
func DropClass(c *gin.Context) {
	var input model.StudentEnrollmentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	studentID, err := requireRegisterPermission(c)
	if err != nil {
		return
	}

	if err := service.DropClass(studentID, input.CourseID); err != nil {
		if err.Error() == "registration not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Class dropped successfully"})
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

// ListClassRegistrations returns all registrations for a class.
func ListClassRegistrations(c *gin.Context) {
	classIDStr := c.Param("id")
	classID, err := strconv.ParseUint(classIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	registrations, err := service.ListClassRegistrations(uint(classID))
	if err != nil {
		if err.Error() == "class not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"registrations": registrations})
}

// GetStudentEnrolledClasses returns all courses a student is enrolled in.
func GetStudentEnrolledClasses(c *gin.Context) {
	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	courses, err := service.GetStudentEnrolledClasses(uint(studentID))
	if err != nil {
		if err.Error() == "student not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func getStudentIDFromAuthHeader(c *gin.Context) (uint, error) {
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

	return service.GetStudentIDFromToken(tokenString)
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
	studentID, err := getStudentIDFromAuthHeader(c)
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

	return studentID, nil
}