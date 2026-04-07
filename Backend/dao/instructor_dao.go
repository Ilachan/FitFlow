package dao

import (
	"my-course-backend/db"
	"my-course-backend/model"
)

func ListCoursesByInstructor(instructorID uint) ([]model.Course, error) {
	var courses []model.Course
	if err := db.DB.
		Where("instructor_id = ?", instructorID).
		Order("start_time ASC").
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func ListEnrollmentsByInstructorCourse(courseID uint) ([]model.Enrollment, error) {
	var enrollments []model.Enrollment
	if err := db.DB.Where("course_id = ?", courseID).
		Preload("User").
		Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}