package dao

import (
	"errors"
	"time"

	"my-course-backend/db"
	"my-course-backend/model"

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

// GetUserByID retrieves a user by ID.
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := db.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CheckEnrollmentExists checks if a user is already enrolled in a course.
func CheckEnrollmentExists(userID uint, courseID uint) (bool, error) {
	var count int64
	if err := db.DB.Model(&model.Enrollment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountEnrollmentsByClass returns the number of enrollments for a course.
func CountEnrollmentsByClass(courseID uint) (int64, error) {
	var count int64
	if err := db.DB.Model(&model.Enrollment{}).
		Where("course_id = ? AND status = ?", courseID, "registered").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CreateEnrollment inserts a new class enrollment.
func CreateEnrollment(enrollment *model.Enrollment) error {
	return db.DB.Create(enrollment).Error
}

// DeleteEnrollment removes a user enrollment by user and course IDs.
func DeleteEnrollment(userID uint, courseID uint) error {
	result := db.DB.Where("user_id = ? AND course_id = ?", userID, courseID).
		Delete(&model.Enrollment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("enrollment not found")
	}
	return nil
}

// ListEnrollmentsByClass returns enrollments for a course with user info.
func ListEnrollmentsByClass(courseID uint) ([]model.Enrollment, error) {
	var enrollments []model.Enrollment
	if err := db.DB.Where("course_id = ?", courseID).
		Preload("User").
		Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

// ListEnrolledCoursesByUser returns all courses a user is enrolled in.
func ListEnrolledCoursesByUser(userID uint) ([]model.Course, error) {
	var courses []model.Course
	if err := db.DB.Joins("INNER JOIN Enrollment ON Enrollment.course_id = Course.id").
		Joins("INNER JOIN User ON User.id = Enrollment.user_id").
		Where("User.id = ? AND Enrollment.status = ?", userID, "registered").
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

// CreateDailyActivity inserts a daily activity row.
func CreateDailyActivity(activity *model.UserDailyActivity) error {
	return db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(activity).Error
}

// BackfillUserDailyActivityFromEnrollments syncs missing daily rows from Enrollment.
func BackfillUserDailyActivityFromEnrollments(userID uint) error {
	query := `
		INSERT INTO UserDailyActivity (enrollment_id, user_id, course_id, activity_date, created_at)
		SELECT e.id, e.user_id, e.course_id, DATE(e.enroll_time), CURRENT_TIMESTAMP
		FROM Enrollment e
		WHERE e.user_id = ?
		AND NOT EXISTS (
			SELECT 1
			FROM UserDailyActivity uda
			WHERE uda.enrollment_id = e.id
			  AND uda.activity_date = DATE(e.enroll_time)
		)
	`

	return db.DB.Exec(query, userID).Error
}

// GetUserActivityStats returns total activity stats in a date range.
func GetUserActivityStats(userID uint, fromDate time.Time, toDate time.Time) (int64, int64, error) {
	type statsResult struct {
		TotalClasses int64
		ActiveDays   int64
	}

	var result statsResult
	err := db.DB.Model(&model.UserDailyActivity{}).
		Select("COUNT(*) as total_classes, COUNT(DISTINCT activity_date) as active_days").
		Where("user_id = ? AND activity_date BETWEEN ? AND ?", userID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).
		Scan(&result).Error
	if err != nil {
		return 0, 0, err
	}

	return result.TotalClasses, result.ActiveDays, nil
}

// GetUserTotalTime returns the summed class duration (minutes) in a date range.
func GetUserTotalTime(userID uint, fromDate time.Time, toDate time.Time) (int64, error) {
	type totalTimeResult struct {
		TotalTime int64
	}

	var result totalTimeResult
	err := db.DB.Table("UserDailyActivity AS uda").
		Select("COALESCE(SUM(COALESCE(c.duration, 0)), 0) AS total_time").
		Joins("INNER JOIN Course c ON c.id = uda.course_id").
		Where("uda.user_id = ? AND uda.activity_date BETWEEN ? AND ?", userID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).
		Scan(&result).Error
	if err != nil {
		return 0, err
	}

	return result.TotalTime, nil
}

// GetUserDailyActivitySummary returns grouped daily analytics for the user.
func GetUserDailyActivitySummary(userID uint, fromDate time.Time, toDate time.Time) ([]model.DailyActivitySummary, error) {
	var daily []model.DailyActivitySummary

	err := db.DB.Table("UserDailyActivity AS uda").
		Select(`DATE(uda.activity_date) AS date,
			COUNT(*) AS classes`).
		Where("uda.user_id = ? AND uda.activity_date BETWEEN ? AND ?", userID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).
		Group("DATE(uda.activity_date)").
		Order("DATE(uda.activity_date) ASC").
		Scan(&daily).Error
	if err != nil {
		return nil, err
	}

	return daily, nil
}

// GetUserCategoryActivitySummary returns grouped category analytics for the user.
func GetUserCategoryActivitySummary(userID uint, fromDate time.Time, toDate time.Time) ([]model.CategoryActivitySummary, error) {
	var categories []model.CategoryActivitySummary

	err := db.DB.Table("UserDailyActivity AS uda").
		Select(`COALESCE(NULLIF(TRIM(c.category), ''), 'Uncategorized') AS category,
			COUNT(*) AS classes`).
		Joins("INNER JOIN Course c ON c.id = uda.course_id").
		Where("uda.user_id = ? AND uda.activity_date BETWEEN ? AND ?", userID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).
		Group("COALESCE(NULLIF(TRIM(c.category), ''), 'Uncategorized')").
		Order("classes DESC, category ASC").
		Scan(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}
