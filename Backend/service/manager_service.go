package service

import (
	"errors"
	"strings"
	"time"

	"my-course-backend/dao"
	"my-course-backend/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CourseUpsertInput kept in manager_service to avoid creating new files.
type CourseUpsertInput struct {
	CourseName  string `json:"name" binding:"required"`
	CourseCode  string `json:"course_code" binding:"required"`
	Description string `json:"description"`

	StartTime string `json:"start_time" binding:"required"` // "08:00" or "08:00:00"
	EndTime   string `json:"end_time" binding:"required"`   // "09:00" or "09:00:00"

	Capacity int `json:"capacity" binding:"required,min=1"`

	// NEW fields
	Duration int    `json:"duration" binding:"omitempty,min=0"`
	Category string `json:"category"`
	Weekday  string `json:"weekday"`
}

// ManagerCreateCourse creates a course (manager role required at API layer).
// ManagerCreateCourse creates a course (manager role required at API layer).
func ManagerCreateCourse(input CourseUpsertInput) (*model.Course, error) {
	start, err := model.ParseTimeOnly(input.StartTime)
	if err != nil {
		return nil, errors.New("invalid start_time, expected HH:MM or HH:MM:SS")
	}
	end, err := model.ParseTimeOnly(input.EndTime)
	if err != nil {
		return nil, errors.New("invalid end_time, expected HH:MM or HH:MM:SS")
	}

	course := &model.Course{
		CourseName:  input.CourseName,
		CourseCode:  input.CourseCode,
		Description: input.Description,
		StartTime:   start,
		EndTime:     end,
		Capacity:    input.Capacity,
		Duration:    input.Duration,
		Category:    input.Category,
		Weekday:     input.Weekday,
	}

	if err := dao.CreateCourse(course); err != nil {
		return nil, err
	}

	_ = fillCourseSpot(course)
	return course, nil
}

// ManagerUpdateCourse updates a course by ID.
func ManagerUpdateCourse(id uint, input CourseUpsertInput) (*model.Course, error) {
	course, err := dao.GetCourseByID(id)
	if err != nil {
		return nil, errors.New("class not found")
	}

	start, err := model.ParseTimeOnly(input.StartTime)
	if err != nil {
		return nil, errors.New("invalid start_time, expected HH:MM or HH:MM:SS")
	}
	end, err := model.ParseTimeOnly(input.EndTime)
	if err != nil {
		return nil, errors.New("invalid end_time, expected HH:MM or HH:MM:SS")
	}

	course.CourseName = input.CourseName
	course.CourseCode = input.CourseCode
	course.Description = input.Description
	course.StartTime = start
	course.EndTime = end
	course.Capacity = input.Capacity
	course.Duration = input.Duration
	course.Category = input.Category
	course.Weekday = input.Weekday

	if err := dao.UpdateCourse(course); err != nil {
		return nil, err
	}

	_ = fillCourseSpot(course)
	return course, nil
}

// ManagerDeleteCourse deletes a course by ID.
func ManagerDeleteCourse(id uint) error {
	if _, err := dao.GetCourseByID(id); err != nil {
		return errors.New("class not found")
	}
	return dao.DeleteCourseByID(id)
}

func RegisterManager(input model.ManagerRegisterInput) error {
	// email uniqueness
	if dao.CheckEmailExist(input.Email) {
		return errors.New("email already exists")
	}

	// Normalize email
	email := strings.TrimSpace(strings.ToLower(input.Email))

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Manager role_id = 3 (per your definition)
	const managerRoleID uint = 3

	// Transaction: validate invite -> create user -> mark invite used
	return dao.WithTx(func(tx *gorm.DB) error {
		invite, err := dao.GetManagerInviteCodeForUpdate(tx, strings.TrimSpace(input.InviteCode))
		if err != nil {
			return errors.New("invalid invite code")
		}

		// Validate status
		if invite.Status == nil || *invite.Status != "active" {
			return errors.New("invite code is not active")
		}

		// Validate not used
		if invite.UsedAt != nil {
			return errors.New("invite code already used")
		}

		// Validate expiry
		if invite.ExpiredAt == nil || invite.ExpiredAt.Before(time.Now()) {
			return errors.New("invite code expired")
		}

		// If invitee_email is set, must match
		if invite.InviteeEmail != nil && strings.TrimSpace(strings.ToLower(*invite.InviteeEmail)) != email {
			return errors.New("invite code not allowed for this email")
		}

		// Create manager user in "User" table
		user := model.User{
			Name:     input.Name,
			Email:    email,
			Password: string(hashedPassword),
			RoleID:   managerRoleID,
		}

		// Use tx to create user (avoid partial commits)
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// Mark invite as used
		if err := dao.MarkInviteCodeUsed(tx, invite.ID, email); err != nil {
			return err
		}

		return nil
	})
}