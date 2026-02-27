package dao

import (
	"errors"

	"my-course-backend/db"
	"my-course-backend/model"

	"gorm.io/gorm"
)

// GetRoleByName retrieves the role ID based on the role name.
func GetRoleByName(name string) (uint, error) {
	var role model.Role
	// Query the database and return the error directly if it fails
	err := db.DB.Where("role_name = ?", name).First(&role).Error
	if err != nil {
		return 0, err
	}
	return role.ID, nil
}

// CheckEmailExist checks if the given email already exists in the database.
func CheckEmailExist(email string) bool {
	var count int64
	// CHANGED: model.Student{} -> model.User{}
	db.DB.Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// CHANGED: CreateStudent -> CreateUser
func CreateUser(user *model.User) error {
	return db.DB.Create(user).Error
}

// CHANGED: GetStudentByEmail -> GetUserByEmail
func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	// Preload("Role") eagerly loads the associated Role data, which is often needed by the frontend.
	err := db.DB.Where("email = ?", email).Preload("Role").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetProfileByID fetches combined data from User and student_info
func GetProfileByID(id uint) (*model.StudentProfile, error) {
	var profile model.StudentProfile

	// CHANGED: all "Student" table refs -> "User"
	// NOTE: student_info table remains unchanged per your request.
	err := db.DB.Table("User").
		Select("User.name, User.email, User.avatar_url, student_info.date_of_birth, student_info.gender, student_info.phone_number, student_info.address").
		Joins("left join student_info on student_info.student_id = User.id").
		Where("User.id = ?", id).
		Scan(&profile).Error

	return &profile, err
}

// UpdateProfile updates both tables in a single transaction
func UpdateProfile(id uint, p model.StudentProfile) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// CHANGED: tx.Model(&model.Student{}) -> tx.Model(&model.User{})
		if err := tx.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
			"name":       p.Name,
			"avatar_url": p.AvatarURL,
		}).Error; err != nil {
			return err
		}

		// 2. Update or Create student_info table (unchanged)
		var info model.StudentInfo
		err := tx.Where("student_id = ?", id).First(&info).Error

		infoData := model.StudentInfo{
			StudentID:   id,
			DateOfBirth: p.DateOfBirth,
			Gender:      p.Gender,
			PhoneNumber: p.PhoneNumber,
			Address:     p.Address,
		}

		if err != nil { // If record doesn't exist, create it
			return tx.Create(&infoData).Error
		}
		// If exists, update it
		return tx.Model(&info).Updates(infoData).Error
	})
}

// CHANGED: DeleteStudentByID -> DeleteUserByID
func DeleteUserByID(id uint) error {
	// CHANGED: model.Student{} -> model.User{}
	result := db.DB.Delete(&model.User{}, id)

	if result.Error != nil {
		return result.Error
	}

	// If RowsAffected is 0, it means the ID does not exist in the database.
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}