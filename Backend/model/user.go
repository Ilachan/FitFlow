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
	Password  string    `gorm:"not null" json:"-"` // never returned
	AvatarURL string    `json:"avatar_url"`
	RoleID    uint      `json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type UserProfile struct {
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	AvatarURL   *string `json:"avatar_url"`
	DateOfBirth *string `json:"date_of_birth"`
	Gender      *string `json:"gender"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
}

// PatchString distinguishes: missing vs null vs value
type PatchString struct {
	Set   bool   // field present in request JSON?
	Valid bool   // if Set==true, was it non-null?
	Value string // if Valid==true, the actual string value
}

type UserProfilePatch struct {
	Name        PatchString `json:"name"`
	AvatarURL   PatchString `json:"avatar_url"`
	DateOfBirth PatchString `json:"date_of_birth"`
	Gender      PatchString `json:"gender"`
	PhoneNumber PatchString `json:"phone_number"`
	Address     PatchString `json:"address"`
}

// UserInfo matches the user_info table structure.
type UserInfo struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint   `gorm:"unique;not null" json:"user_id"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

func (UserInfo) TableName() string { return "user_info" }
func (User) TableName() string     { return "User" }
func (Role) TableName() string     { return "Role" }

// RegisterInput captures the parameters sent by the frontend for registration.
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput captures the parameters sent by the frontend for login.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}