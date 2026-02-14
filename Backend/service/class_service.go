package service

import (
	"errors"
	"my-course-backend/dao"
	"my-course-backend/model"
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
	return dao.CreateRegistration(&registration)
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
