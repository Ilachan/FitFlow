package dao

import (
	"errors"
	"my-course-backend/db"
	"my-course-backend/model"
	"time"

	"gorm.io/gorm/clause"
)

// GetCourseByID retrieves a course by ID.
func GetCourseByID(id uint) (*model.Course, error) {
	var class model.Course
	if err := db.DB.First(&class, id).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

// ListClasses retrieves all courses.
func ListClasses() ([]model.Course, error) {
	var classes []model.Course
	if err := db.DB.Order("start_time ASC").Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

// CreateClass inserts a new course.
// Disabled: classes are imported manually into the DB.
// func CreateClass(class *model.Course) error {
// 	return db.DB.Create(class).Error
// }

// GetStudentByID retrieves a student by ID.
func GetStudentByID(id uint) (*model.Student, error) {
	var student model.Student
	if err := db.DB.First(&student, id).Error; err != nil {
		return nil, err
	}
	return &student, nil
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

// ListRegistrationsByClass returns registrations for a course with student info.
func ListRegistrationsByClass(courseID uint) ([]model.StudentEnrollment, error) {
	var registrations []model.StudentEnrollment
	if err := db.DB.Where("course_id = ?", courseID).
		Preload("Student").
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

// CreateDailyActivity inserts a daily activity row.
func CreateDailyActivity(activity *model.StudentDailyActivity) error {
	return db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(activity).Error
}

// BackfillStudentDailyActivityFromEnrollments syncs missing daily rows from StudentEnrollment.
func BackfillStudentDailyActivityFromEnrollments(studentID uint) error {
	query := `
		INSERT INTO StudentDailyActivity (enrollment_id, student_id, course_id, activity_date, created_at)
		SELECT se.id, se.student_id, se.course_id, DATE(se.enroll_time), CURRENT_TIMESTAMP
		FROM StudentEnrollment se
		WHERE se.student_id = ?
		AND NOT EXISTS (
			SELECT 1
			FROM StudentDailyActivity sda
			WHERE sda.enrollment_id = se.id
			  AND sda.activity_date = DATE(se.enroll_time)
		)
	`

	return db.DB.Exec(query, studentID).Error
}

// GetStudentActivityStats returns total activity stats in a date range.
func GetStudentActivityStats(studentID uint, fromDate time.Time, toDate time.Time) (int64, int64, error) {
	type statsResult struct {
		TotalClasses int64
		ActiveDays   int64
	}

	var result statsResult
	err := db.DB.Model(&model.StudentDailyActivity{}).
		Select("COUNT(*) as total_classes, COUNT(DISTINCT activity_date) as active_days").
		Where("student_id = ? AND activity_date BETWEEN ? AND ?", studentID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).
		Scan(&result).Error
	if err != nil {
		return 0, 0, err
	}

	return result.TotalClasses, result.ActiveDays, nil
}

// GetStudentDailyActivitySummary returns grouped daily analytics for the student.
func GetStudentDailyActivitySummary(studentID uint, fromDate time.Time, toDate time.Time) ([]model.DailyActivitySummary, error) {
	var daily []model.DailyActivitySummary

	err := db.DB.Table("StudentDailyActivity AS sda").
		Select(`DATE(sda.activity_date) AS date,
			COUNT(*) AS classes`).
		Where("sda.student_id = ? AND sda.activity_date BETWEEN ? AND ?", studentID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).
		Group("DATE(sda.activity_date)").
		Order("DATE(sda.activity_date) ASC").
		Scan(&daily).Error
	if err != nil {
		return nil, err
	}

	return daily, nil
}

// GetStudentCategoryActivitySummary returns grouped category analytics for the student.
func GetStudentCategoryActivitySummary(studentID uint, fromDate time.Time, toDate time.Time) ([]model.CategoryActivitySummary, error) {
	var categories []model.CategoryActivitySummary

	err := db.DB.Table("StudentDailyActivity AS sda").
		Select(`COALESCE(NULLIF(TRIM(c.category), ''), 'Uncategorized') AS category,
			COUNT(*) AS classes`).
		Joins("INNER JOIN Course c ON c.id = sda.course_id").
		Where("sda.student_id = ? AND sda.activity_date BETWEEN ? AND ?", studentID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).
		Group("COALESCE(NULLIF(TRIM(c.category), ''), 'Uncategorized')").
		Order("classes DESC, category ASC").
		Scan(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}
