package model

import "time"

// Role represents the role table in the database.
type Role struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleName string `gorm:"unique;not null" json:"role_name"`
}

// CHANGED: Student -> User (because DB table renamed to "User")
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // json:"-" ensures the password is never returned in the API response
	AvatarURL string    `json:"avatar_url"`
	RoleID    uint      `json:"role_id"`                          // Foreign Key
	Role      Role      `gorm:"foreignKey:RoleID" json:"role"`    // Association (Belongs To)
	CreatedAt time.Time `json:"created_at"`
}

// StudentProfile represents the combined data from Student and student_info tables
// NOTE: unchanged on purpose - only Student table renamed per your request.
type StudentProfile struct {
	Name        string `json:"name"`
	Email       string `json:"email"` // Usually read-only in profile updates
	AvatarURL   string `json:"avatar_url"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

// StudentInfo matches your SQL student_info table structure
// NOTE: unchanged on purpose - only Student table renamed per your request.
type StudentInfo struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID   uint   `gorm:"unique;not null" json:"student_id"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

func (StudentInfo) TableName() string {
	return "student_info"
}

// CHANGED: TableName now points to "User"
func (User) TableName() string {
	return "User"
}

// TableName overrides the default table name to "Role".
func (Role) TableName() string {
	return "Role"
}

// RegisterInput captures the parameters sent by the frontend for registration.
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"` // Password must be at least 6 characters
}

// LoginInput captures the parameters sent by the frontend for login.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}