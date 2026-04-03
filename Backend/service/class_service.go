package service

import (
	"errors"
	"math"
	"my-course-backend/dao"
	"my-course-backend/model"
	"time"
)

// RegisterClass enrolls a user in a course.
func RegisterClass(userID uint, courseID uint) error {
	if _, err := dao.GetUserByID(userID); err != nil {
		return errors.New("user not found")
	}

	class, err := dao.GetCourseByID(courseID)
	if err != nil {
		return errors.New("class not found")
	}

	exists, err := dao.CheckEnrollmentExists(userID, courseID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("enrollment already exists")
	}

	count, err := dao.CountEnrollmentsByClass(courseID)
	if err != nil {
		return err
	}
	if int(count) >= class.Capacity {
		return errors.New("class is full")
	}

	enrollment := model.Enrollment{
		UserID:   userID,
		CourseID: courseID,
		Status:   model.EnrollmentStatusEnrolled,
	}
	if err := dao.CreateEnrollment(&enrollment); err != nil {
		return err
	}

	return dao.BackfillUserDailyActivityFromEnrollments(userID)
}

// DropClass removes a user's enrollment from a course.
func DropClass(userID uint, courseID uint) error {
	if err := dao.DeleteEnrollment(userID, courseID); err != nil {
		return err
	}

	return nil
}

// ListClassEnrollments returns all enrollments for a course.
func ListClassEnrollments(courseID uint) ([]model.Enrollment, error) {
	if _, err := dao.GetCourseByID(courseID); err != nil {
		return nil, errors.New("class not found")
	}
	return dao.ListEnrollmentsByClass(courseID)
}

// CHANGED/NEW: ListClassesPaged returns paged result with spot filled.
func ListClassesPaged(page int, pageSize int) ([]model.Course, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	classes, total, err := dao.ListClassesPaged(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	for i := range classes {
		if err := fillCourseSpot(&classes[i]); err != nil {
			return nil, 0, err
		}
	}

	return classes, total, nil
}

func fillCourseSpot(class *model.Course) error {
	count, err := dao.CountEnrollmentsByClass(class.ID)
	if err != nil {
		return err
	}
	spot := class.Capacity - int(count)
	if spot < 0 {
		spot = 0
	}
	class.Spot = spot
	return nil
}

// ListClasses returns all courses with spot populated.
func ListClasses() ([]model.Course, error) {
	classes, err := dao.ListClasses()
	if err != nil {
		return nil, err
	}
	for i := range classes {
		if err := fillCourseSpot(&classes[i]); err != nil {
			return nil, err
		}
	}
	return classes, nil
}

// GetClass returns a single class by ID with spot populated.
func GetClass(courseID uint) (*model.Course, error) {
	class, err := dao.GetCourseByID(courseID)
	if err != nil {
		return nil, errors.New("class not found")
	}
	if err := fillCourseSpot(class); err != nil {
		return nil, err
	}
	return class, nil
}

// GetUserEnrolledClasses returns all courses a user is enrolled in with spot populated.
func GetUserEnrolledClasses(userID uint) ([]model.Course, error) {
	if _, err := dao.GetUserByID(userID); err != nil {
		return nil, errors.New("user not found")
	}

	courses, err := dao.ListEnrolledCoursesByUser(userID)
	if err != nil {
		return nil, err
	}

	for i := range courses {
		if err := fillCourseSpot(&courses[i]); err != nil {
			return nil, err
		}
	}

	return courses, nil
}

// GetStudentAnalytics returns dashboard analytics for a date range.
func GetUserAnalytics(userID uint, rangeKey string) (*model.UserAnalyticsResponse, error) {
	if _, err := dao.GetUserByID(userID); err != nil {
		return nil, errors.New("user not found")
	}

	if err := dao.BackfillUserDailyActivityFromEnrollments(userID); err != nil {
		return nil, err
	}

	toDate := time.Now()
	fromDate := resolveRangeStart(rangeKey, toDate)

	totalClasses, activeDays, err := dao.GetUserActivityStats(userID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	totalTime, err := dao.GetUserTotalTime(userID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	daily, err := dao.GetUserDailyActivitySummary(userID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	categories, err := dao.GetUserCategoryActivitySummary(userID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	for i := range categories {
		if totalClasses <= 0 {
			categories[i].Percentage = 0
			continue
		}
		percentage := (float64(categories[i].Classes) / float64(totalClasses)) * 100
		categories[i].Percentage = math.Round(percentage*100) / 100
	}

	response := &model.UserAnalyticsResponse{
		UserID:       userID,
		Range:        normalizeRangeKey(rangeKey),
		FromDate:     fromDate.Format("2006-01-02"),
		ToDate:       toDate.Format("2006-01-02"),
		TotalClasses: totalClasses,
		TotalTime:    totalTime,
		ActiveDays:   activeDays,
		Daily:        daily,
		Categories:   categories,
	}

	return response, nil
}

func resolveRangeStart(rangeKey string, now time.Time) time.Time {
	key := normalizeRangeKey(rangeKey)

	switch key {
	case "1m":
		return now.AddDate(0, -1, 0)
	case "3m":
		return now.AddDate(0, -3, 0)
	default:
		return now.AddDate(0, 0, -7)
	}
}

func normalizeRangeKey(rangeKey string) string {
	switch rangeKey {
	case "1m", "3m":
		return rangeKey
	default:
		return "7d"
	}
}
