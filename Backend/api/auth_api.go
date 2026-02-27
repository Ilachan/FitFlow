package api

import (
	"net/http"
	"strconv"
	"strings"

	// Ensure these match your go.mod module name
	"my-course-backend/model"
	"my-course-backend/service"

	"github.com/gin-gonic/gin"
)

// Register handles user registration with specific error codes
func Register(c *gin.Context) {
	var input model.RegisterInput

	// 400: Validation error (e.g., missing fields, invalid email format)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}

	if err := service.RegisterUser(input); err != nil {
		// 409: Conflict (Email already exists)
		if err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "This email is already registered"})
		} else {
			// 500: Internal Server Error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

// CHANGED: Login now returns role_id in response.
func Login(c *gin.Context) {
	var input model.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide email and password"})
		return
	}

	token, roleID, err := service.LoginUserWithRole(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"role_id": roleID, // NEW
	})
}

// GetProfile handles GET /auth/profile by manually verifying the JWT
func GetProfile(c *gin.Context) {
	// 1. Check for Authorization Header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		// 401: Client didn't provide credentials
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	// 2. Validate Bearer Token Format
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		// 400: The request was formed incorrectly
		c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed token. Please use 'Bearer <token>' format"})
		return
	}

	// 3. Parse Token and Handle Expiration/Invalidity
	userID, err := service.GetStudentIDFromToken(tokenString)
	if err != nil {
		errorMessage := err.Error()
		if strings.Contains(errorMessage, "expired") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired. Please log in again"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid security token. Authentication failed"})
		}
		return
	}

	// 4. Fetch the Profile from Database
	profile, err := service.GetUserProfile(userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Detailed profile information could not be found for this student"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "A server-side error occurred while retrieving your profile"})
		}
		return
	}

	c.JSON(http.StatusOK, profile)
}

func UpdateProfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := service.GetStudentIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	var input model.StudentProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := service.UpdateUserProfile(userID, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// DeleteUser handles user deletion
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := service.RemoveUser(uint(id)); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User and related data deleted successfully"})
}