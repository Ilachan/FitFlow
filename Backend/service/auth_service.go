package service

import (
	"errors"
	"my-course-backend/dao"
	"my-course-backend/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Secret key for JWT signing.
// WARNING: In a production environment, store this in environment variables (e.g., os.Getenv("JWT_SECRET")), DO NOT hardcode it here.
var jwtSecret = []byte("my_super_secret_key_2026")

// ==========================================
// 1. Registration Logic (RegisterUser)
// ==========================================
func RegisterUser(input model.RegisterInput) error {
	if dao.CheckEmailExist(input.Email) {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	roleID, err := dao.GetRoleByName("Student")
	if err != nil {
		return errors.New("role 'Student' not found in database")
	}

	// NOTE: assumes you've already renamed Student model to User in your codebase
	user := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:   roleID,
	}

	return dao.CreateUser(&user)
}

// ==========================================
// 2. Login Logic (LoginUser)
// ==========================================
func LoginUser(input model.LoginInput) (string, error) {
	user, err := dao.GetUserByEmail(input.Email)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"role_id": user.RoleID,
		"exp":     time.Now().Add(time.Hour * 48).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetStudentIDFromToken extracts the user ID from a JWT token string.
// (function name kept to avoid touching other callers)
func GetStudentIDFromToken(tokenString string) (uint, error) {
	claims := jwt.MapClaims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !parsed.Valid {
		return 0, errors.New("invalid token")
	}

	idValue, ok := claims["id"]
	if !ok {
		return 0, errors.New("invalid token")
	}

	switch typed := idValue.(type) {
	case float64:
		return uint(typed), nil
	case int:
		return uint(typed), nil
	case uint:
		return typed, nil
	default:
		return 0, errors.New("invalid token")
	}
}

// CHANGED: new helper to extract role_id from token
func GetRoleIDFromToken(tokenString string) (uint, error) {
	claims := jwt.MapClaims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !parsed.Valid {
		return 0, errors.New("invalid token")
	}

	roleValue, ok := claims["role_id"]
	if !ok {
		return 0, errors.New("invalid token")
	}

	switch typed := roleValue.(type) {
	case float64:
		return uint(typed), nil
	case int:
		return uint(typed), nil
	case uint:
		return typed, nil
	default:
		return 0, errors.New("invalid token")
	}
}

func RemoveUser(id uint) error {
	return dao.DeleteUserByID(id)
}

func GetUserProfile(id uint) (*model.StudentProfile, error) {
	return dao.GetProfileByID(id)
}

func UpdateUserProfile(id uint, input model.StudentProfile) error {
	return dao.UpdateProfile(id, input)
}

// CHANGED: LoginUserWithRole returns both token and role_id for frontend convenience.
func LoginUserWithRole(input model.LoginInput) (string, uint, error) {
	user, err := dao.GetUserByEmail(input.Email)
	if err != nil {
		return "", 0, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", 0, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"role_id": user.RoleID,
		"exp":     time.Now().Add(time.Hour * 70).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, user.RoleID, nil
}