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

// Login handles user authentication
func Login(c *gin.Context) {
	var input model.LoginInput

	// 400: Bad Request
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide email and password"})
		return
	}

	token, err := service.LoginUser(input)
	if err != nil {
		// 401: Unauthorized (Wrong credentials or user not found)
		// We use a generic message for security reasons 
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
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
	// Standard format is: "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader { 
		// 400: The request was formed incorrectly
		c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed token. Please use 'Bearer <token>' format"})
		return
	}

	// 3. Parse Token and Handle Expiration/Invalidity
	userID, err := service.GetStudentIDFromToken(tokenString)
	if err != nil {
		// Differentiate between an expired token and a fake/invalid one
		errorMessage := err.Error()
		if strings.Contains(errorMessage, "expired") {
			// 401: Token was real but is now too old
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired. Please log in again"})
		} else {
			// 401: Token failed verification (security risk)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid security token. Authentication failed"})
		}
		return
	}

	// 4. Fetch the Profile from Database
	profile, err := service.GetUserProfile(userID)
	if err != nil {
		// If the error message from DAO says "user not found"
		if strings.Contains(err.Error(), "not found") {
			// 404: The student exists but their profile info is missing
			c.JSON(http.StatusNotFound, gin.H{"error": "Detailed profile information could not be found for this student"})
		} else {
			// 500: Database connection issue or other server-side failures
			c.JSON(http.StatusInternalServerError, gin.H{"error": "A server-side error occurred while retrieving your profile"})
		}
		return
	}

	// 200: Success
	c.JSON(http.StatusOK, profile)
}

func UpdateProfile(c *gin.Context) {
	// 1. Manually get the Authorization header (Same as GetProfile)
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	
	// 2. Parse the userID from the token manually
	userID, err := service.GetStudentIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// 3. Now it is safe to proceed with the rest of the logic
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
	// 1. Get 'id' parameter from URL (string type)
	idStr := c.Param("id")

	// 2. Convert string to uint
	// ParseUint(string, base, bitSize)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 3. Call Service
	if err := service.RemoveUser(uint(id)); err != nil {
		// Return 404 if user not found, otherwise 500
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 4. Return success
	c.JSON(http.StatusOK, gin.H{"message": "User and related data deleted successfully"})
}