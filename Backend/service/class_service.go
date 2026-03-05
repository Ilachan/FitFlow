package service

import (
	"errors"
	"math"
	"my-course-backend/dao"
	"my-course-backend/model"
	"time"
)

// RegisterClass registers a student for a course.
func RegisterClass(studentID uint, courseID uint) error {
	if _, err := dao.GetStudentByID(studentID); err != nil {
		return errors.New("student not found")
	}

	class, err := dao.GetCourseByID(courseID)
	if err != nil {
		return errors.New("class not found")
	}

	exists, err := dao.CheckRegistrationExists(studentID, courseID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("registration already exists")
	}

	count, err := dao.CountRegistrationsByClass(courseID)
	if err != nil {
		return err
	}
	if int(count) >= class.Capacity {
		return errors.New("class is full")
	}

	registration := model.StudentEnrollment{
		StudentID: studentID,
		CourseID:  courseID,
		Status:    "registered",
	}
	if err := dao.CreateRegistration(&registration); err != nil {
		return err
	}

	activity := model.StudentDailyActivity{
		EnrollmentID: registration.ID,
		StudentID:    studentID,
		CourseID:     courseID,
		ActivityDate: time.Now(),
	}

	return dao.CreateDailyActivity(&activity)
}

// DropClass removes a student's registration from a course.
func DropClass(studentID uint, courseID uint) error {
	return dao.DeleteRegistration(studentID, courseID)
}

// ListClassRegistrations returns all registrations for a course.
func ListClassRegistrations(courseID uint) ([]model.StudentEnrollment, error) {
	if _, err := dao.GetCourseByID(courseID); err != nil {
		return nil, errors.New("class not found")
	}
	return dao.ListRegistrationsByClass(courseID)
}

func fillCourseSpot(class *model.Course) error {
	count, err := dao.CountRegistrationsByClass(class.ID)
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

// GetStudentEnrolledClasses returns all courses a student is enrolled in with spot populated.
func GetStudentEnrolledClasses(studentID uint) ([]model.Course, error) {
	// Verify student exists
	if _, err := dao.GetStudentByID(studentID); err != nil {
		return nil, errors.New("student not found")
	}

	courses, err := dao.ListEnrolledCoursesByStudent(studentID)
	if err != nil {
		return nil, err
	}

	// Populate spot for each course
	for i := range courses {
		if err := fillCourseSpot(&courses[i]); err != nil {
			return nil, err
		}
	}

	return courses, nil
}

// GetStudentAnalytics returns dashboard analytics for a date range.
func GetStudentAnalytics(studentID uint, rangeKey string) (*model.StudentAnalyticsResponse, error) {
	if _, err := dao.GetStudentByID(studentID); err != nil {
		return nil, errors.New("student not found")
	}

	if err := dao.BackfillStudentDailyActivityFromEnrollments(studentID); err != nil {
		return nil, err
	}

	toDate := time.Now()
	fromDate := resolveRangeStart(rangeKey, toDate)

	totalClasses, activeDays, err := dao.GetStudentActivityStats(studentID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	daily, err := dao.GetStudentDailyActivitySummary(studentID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	categories, err := dao.GetStudentCategoryActivitySummary(studentID, fromDate, toDate)
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

	response := &model.StudentAnalyticsResponse{
		StudentID:    studentID,
		Range:        normalizeRangeKey(rangeKey),
		FromDate:     fromDate.Format("2006-01-02"),
		ToDate:       toDate.Format("2006-01-02"),
		TotalClasses: totalClasses,
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
