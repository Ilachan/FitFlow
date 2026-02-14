package model

import "time"

// Course represents a class that students can register for.
type Course struct {
	ID          uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseName  string   `gorm:"not null" json:"name"`
	CourseCode  string   `gorm:"not null" json:"course_code"`
	Description string   `json:"description"`
	StartTime   TimeOnly `gorm:"type:time" json:"start_time"`
	EndTime     TimeOnly `gorm:"type:time" json:"end_time"`
	Spot        int      `gorm:"-" json:"spot"` // Transient field to show remaining spots (calculated in service)
	Capacity    int      `gorm:"not null" json:"capacity"`
}

// StudentEnrollment is the join table between Student and Course.
type StudentEnrollment struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID  uint      `gorm:"column:student_id;not null" json:"student_id"`
	CourseID   uint      `gorm:"column:course_id;not null" json:"course_id"`
	Student    Student   `gorm:"foreignKey:StudentID" json:"student"`
	Course     Course    `gorm:"foreignKey:CourseID" json:"course"`
	Status     string    `gorm:"column:status;not null" json:"status"` // "registered" or "dropped"
	EnrollTime time.Time `gorm:"column:enroll_time;autoCreateTime" json:"enroll_time"`
}

// TableName overrides the default table name to "Course".
func (Course) TableName() string {
	return "Course"
}

// TableName overrides the default table name to "StudentEnrollment".
func (StudentEnrollment) TableName() string {
	return "StudentEnrollment"
}

// StudentEnrollmentRequest captures parameters to register or drop a course.
type StudentEnrollmentRequest struct {
	CourseID uint `json:"course_id" binding:"required"`
}
