package service

import (
	"errors"
	"my-course-backend/dao"
	"my-course-backend/model"
)

func ListInstructorCourses(instructorID uint) ([]model.Course, error) {
	return dao.ListCoursesByInstructor(instructorID)
}

func ListInstructorCourseEnrollments(instructorID uint, courseID uint) ([]model.Enrollment, error) {
	course, err := dao.GetCourseByID(courseID)
	if err != nil {
		return nil, errors.New("class not found")
	}
	if course.InstructorID != instructorID {
		return nil, errors.New("forbidden")
	}
	return dao.ListEnrollmentsByInstructorCourse(courseID)
}

func UpdateEnrollmentStatusByInstructor(instructorID, courseID, userID uint, status string) error {
	if status != "attended" && status != "missed" && status != "enrolled" {
		return errors.New("invalid status")
	}

	course, err := dao.GetCourseByID(courseID)
	if err != nil {
		return errors.New("class not found")
	}

	if course.InstructorID != instructorID {
		return errors.New("forbidden")
	}

	// ✅ 校验用户是否选过课
	if _, err := dao.GetEnrollment(userID, courseID); err != nil {
		return errors.New("enrollment not found")
	}

	ok, err := dao.UpdateEnrollmentStatus(userID, courseID, status)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("enrollment not found")
	}
	return nil
}