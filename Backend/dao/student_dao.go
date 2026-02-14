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
	db.DB.Model(&model.Student{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// CreateStudent creates a new student record in the database.
func CreateStudent(student *model.Student) error {
	return db.DB.Create(student).Error
}

// GetStudentByEmail finds a student by email and preloads their Role information.
func GetStudentByEmail(email string) (*model.Student, error) {
	var student model.Student
	// Preload("Role") eagerly loads the associated Role data, which is often needed by the frontend.
	err := db.DB.Where("email = ?", email).Preload("Role").First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// GetProfileByID fetches combined data from Student and student_info
func GetProfileByID(id uint) (*model.StudentProfile, error) {
	var profile model.StudentProfile
	// Join Student table with student_info table using student_id 
	err := db.DB.Table("Student").
		Select("Student.name, Student.email, Student.avatar_url, student_info.date_of_birth, student_info.gender, student_info.phone_number, student_info.address").
		Joins("left join student_info on student_info.student_id = Student.id").
		Where("Student.id = ?", id).
		Scan(&profile).Error
	return &profile, err
}

// UpdateProfile updates both tables in a single transaction
func UpdateProfile(id uint, p model.StudentProfile) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Update Student table [cite: 7]
		if err := tx.Model(&model.Student{}).Where("id = ?", id).Updates(map[string]interface{}{
			"name":       p.Name,
			"avatar_url": p.AvatarURL,
		}).Error; err != nil {
			return err
		}

		// 2. Update or Create student_info table 
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

// DeleteStudentByID deletes a student record by their ID.
func DeleteStudentByID(id uint) error {
	// Perform the delete operation
	result := db.DB.Delete(&model.Student{}, id)
	
	if result.Error != nil {
		return result.Error
	}
	
	// If RowsAffected is 0, it means the ID does not exist in the database.
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	
	return nil
}