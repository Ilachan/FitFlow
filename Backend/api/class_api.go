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
func RegisterClass(c *gin.Context) {
	var input model.StudentEnrollmentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	studentID, err := getStudentIDFromAuthHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
func DropClass(c *gin.Context) {
	var input model.StudentEnrollmentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	studentID, err := getStudentIDFromAuthHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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

// CreateClass adds a new class to the catalog.
// Disabled: classes are imported manually into the DB.
// func CreateClass(c *gin.Context) {
// 	var input model.CreateClassInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	if err := service.CreateClass(input); err != nil {
// 		if err.Error() == "capacity must be at least 1" {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	c.JSON(http.StatusCreated, gin.H{"message": "Class created successfully"})
// }

// ListClasses returns all courses.
func ListClasses(c *gin.Context) {
	classes, err := service.ListClasses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"classes": classes})
}

// GetClass returns a single course by ID.
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

	authStudentID, err := getStudentIDFromAuthHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if uint(studentID) != authStudentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	courses, err := service.GetStudentEnrolledClasses(authStudentID)
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

// GetStudentAnalytics returns student dashboard analytics for 7d, 1m, or 3m.
func GetStudentAnalytics(c *gin.Context) {
	studentIDStr := c.Param("id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	authStudentID, err := getStudentIDFromAuthHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if uint(studentID) != authStudentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	rangeKey := strings.TrimSpace(c.DefaultQuery("range", "7d"))

	analytics, err := service.GetStudentAnalytics(authStudentID, rangeKey)
	if err != nil {
		if err.Error() == "student not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analytics": analytics})
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
