package dao

import (
	"errors"

	"my-course-backend/db"
	"my-course-backend/model"
)

// GetCourseByID retrieves a course by ID.
func GetCourseByID(id uint) (*model.Course, error) {
	var class model.Course
	if err := db.DB.First(&class, id).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

// ListClassesPaged returns paginated courses and total count.
func ListClassesPaged(limit int, offset int) ([]model.Course, int64, error) {
	var classes []model.Course
	var total int64

	if err := db.DB.Model(&model.Course{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.DB.Order("start_time ASC").Limit(limit).Offset(offset).Find(&classes).Error; err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}

// ListClasses retrieves all courses.
func ListClasses() ([]model.Course, error) {
	var classes []model.Course
	if err := db.DB.Order("start_time ASC").Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

// GetStudentByID now reads from User table/model.
// Note: function name kept as GetStudentByID to reduce cross-file changes.
func GetStudentByID(id uint) (*model.User, error) {
	var user model.User
	if err := db.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CheckRegistrationExists checks if a student is already registered for a course.
func CheckRegistrationExists(studentID uint, courseID uint) (bool, error) {
	var count int64
	if err := db.DB.Model(&model.StudentEnrollment{}).
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountRegistrationsByClass returns the number of registrations for a course.
func CountRegistrationsByClass(courseID uint) (int64, error) {
	var count int64
	if err := db.DB.Model(&model.StudentEnrollment{}).
		Where("course_id = ? AND status = ?", courseID, "registered").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CreateRegistration inserts a new class registration.
func CreateRegistration(registration *model.StudentEnrollment) error {
	return db.DB.Create(registration).Error
}

// DeleteRegistration removes a student enrollment by student and course IDs.
func DeleteRegistration(studentID uint, courseID uint) error {
	result := db.DB.Where("student_id = ? AND course_id = ?", studentID, courseID).
		Delete(&model.StudentEnrollment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("registration not found")
	}
	return nil
}

// ListRegistrationsByClass returns registrations for a course with user info.
func ListRegistrationsByClass(courseID uint) ([]model.StudentEnrollment, error) {
	var registrations []model.StudentEnrollment
	if err := db.DB.Where("course_id = ?", courseID).
		Preload("User").
		Find(&registrations).Error; err != nil {
		return nil, err
	}
	return registrations, nil
}

// ListEnrolledCoursesByStudent returns all courses a student is enrolled in.
func ListEnrolledCoursesByStudent(studentID uint) ([]model.Course, error) {
	var courses []model.Course
	if err := db.DB.Joins("INNER JOIN StudentEnrollment ON StudentEnrollment.course_id = Course.id").
		Where("StudentEnrollment.student_id = ? AND StudentEnrollment.status = ?", studentID, "registered").
		Order("Course.start_time ASC").
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

// NEW: CreateCourse inserts a new course.
func CreateCourse(course *model.Course) error {
	return db.DB.Create(course).Error
}

// NEW: UpdateCourse updates an existing course (all fields).
func UpdateCourse(course *model.Course) error {
	return db.DB.Save(course).Error
}

// NEW: DeleteCourseByID deletes a course by ID.
func DeleteCourseByID(id uint) error {
	return db.DB.Delete(&model.Course{}, id).Error
}