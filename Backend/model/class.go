package model

import "time"

// Course represents the Course table in SQLite.
type Course struct {
	ID uint `gorm:"primaryKey;autoIncrement;column:id" json:"id"`

	CourseName string `gorm:"column:course_name;not null" json:"name"`
	CourseCode string `gorm:"column:course_code;not null" json:"course_code"`

	Description string `gorm:"column:description" json:"description"`

	StartTime TimeOnly `gorm:"column:start_time;type:time" json:"start_time"`
	EndTime   TimeOnly `gorm:"column:end_time;type:time" json:"end_time"`

	Capacity int `gorm:"column:capacity;not null" json:"capacity"`

	Duration int    `gorm:"column:duration" json:"duration"`
	Category string `gorm:"column:category" json:"category"`
	Weekday  string `gorm:"column:weekday" json:"weekday"`

	Spot int `gorm:"-" json:"spot"`
}

// StudentEnrollment is the join table between User and Course.
type StudentEnrollment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID uint      `gorm:"column:student_id;not null" json:"student_id"`
	CourseID  uint      `gorm:"column:course_id;not null" json:"course_id"`

	// NOTE: keep as your current design (recommended name: User)
	User   User   `gorm:"foreignKey:StudentID" json:"user"`
	Course Course `gorm:"foreignKey:CourseID" json:"course"`

	Status     string    `gorm:"column:status;not null" json:"status"`
	EnrollTime time.Time `gorm:"column:enroll_time;autoCreateTime" json:"enroll_time"`
}

func (Course) TableName() string {
	return "Course"
}

func (StudentEnrollment) TableName() string {
	return "StudentEnrollment"
}

type StudentEnrollmentRequest struct {
	CourseID uint `json:"course_id" binding:"required"`
}