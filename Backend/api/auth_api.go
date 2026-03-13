package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"my-course-backend/model"
	"my-course-backend/service"

	"github.com/gin-gonic/gin"
)

// Register handles user registration with specific error codes
func Register(c *gin.Context) {
	var input model.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}

	if err := service.RegisterUser(input); err != nil {
		if err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "This email is already registered"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

// Login now returns role_id in response.
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
		"role_id": roleID,
	})
}

// GetProfile handles GET /auth/profile by manually verifying the JWT
func GetProfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed token. Please use 'Bearer <token>' format"})
		return
	}

	userID, err := service.ExtractUserIDFromToken(tokenString)
	if err != nil {
		errorMessage := err.Error()
		if strings.Contains(errorMessage, "expired") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired. Please log in again"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid security token. Authentication failed"})
		}
		return
	}

	profile, err := service.GetUserProfile(userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Detailed profile information could not be found for this user"})
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

	userID, err := service.ExtractUserIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Use RawMessage map to distinguish:
	// - missing key (undefined)
	// - key with null
	// - key with actual string
	var body map[string]*json.RawMessage
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	parsePatchString := func(key string) (model.PatchString, error) {
		raw, ok := body[key]
		if !ok {
			return model.PatchString{Set: false}, nil // undefined => don't update
		}

		// key exists
		if raw == nil || string(*raw) == "null" {
			return model.PatchString{Set: true, Valid: false}, nil // explicit null
		}

		var val string
		if err := json.Unmarshal(*raw, &val); err != nil {
			return model.PatchString{}, err
		}
		return model.PatchString{Set: true, Valid: true, Value: val}, nil
	}

	var patch model.UserProfilePatch

	if patch.Name, err = parsePatchString("name"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field: name"})
		return
	}
	if patch.AvatarURL, err = parsePatchString("avatar_url"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field: avatar_url"})
		return
	}
	if patch.Gender, err = parsePatchString("gender"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field: gender"})
		return
	}
	if patch.PhoneNumber, err = parsePatchString("phone_number"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field: phone_number"})
		return
	}
	if patch.Address, err = parsePatchString("address"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field: address"})
		return
	}
	if patch.DateOfBirth, err = parsePatchString("date_of_birth"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field: date_of_birth"})
		return
	}

	if err := service.UpdateUserProfilePatch(userID, patch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database update failed: " + err.Error()})
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