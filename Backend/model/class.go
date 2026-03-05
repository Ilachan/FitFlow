package model

import "time"

// Course represents a class that students can register for.
type Course struct {
	ID          uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseName  string   `gorm:"not null" json:"name"`
	CourseCode  string   `gorm:"not null" json:"course_code"`
	Description string   `json:"description"`
	Duration    int      `gorm:"default:0" json:"duration"`
	Weekday     string   `gorm:"default:''" json:"weekday"`
	Category    string   `gorm:"default:''" json:"category"`
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

// StudentDailyActivity stores one activity row per student/class/day for analytics.
type StudentDailyActivity struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EnrollmentID uint      `gorm:"column:enrollment_id;not null;index" json:"enrollment_id"`
	StudentID    uint      `gorm:"column:student_id;not null;index:idx_activity_student_date;index:idx_activity_unique,unique,priority:1" json:"student_id"`
	CourseID     uint      `gorm:"column:course_id;not null;index:idx_activity_unique,unique,priority:2" json:"course_id"`
	ActivityDate time.Time `gorm:"column:activity_date;type:date;not null;index:idx_activity_student_date;index:idx_activity_unique,unique,priority:3" json:"activity_date"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type DailyActivitySummary struct {
	Date    string `json:"date"`
	Classes int64  `json:"classes"`
}

type CategoryActivitySummary struct {
	Category   string  `json:"category"`
	Classes    int64   `json:"classes"`
	Percentage float64 `json:"percentage"`
}

type StudentAnalyticsResponse struct {
	StudentID    uint                      `json:"student_id"`
	Range        string                    `json:"range"`
	FromDate     string                    `json:"from_date"`
	ToDate       string                    `json:"to_date"`
	TotalClasses int64                     `json:"total_classes"`
	ActiveDays   int64                     `json:"active_days"`
	Daily        []DailyActivitySummary    `json:"daily"`
	Categories   []CategoryActivitySummary `json:"categories"`
}

// TableName overrides the default table name to "Course".
func (Course) TableName() string {
	return "Course"
}

// TableName overrides the default table name to "StudentEnrollment".
func (StudentEnrollment) TableName() string {
	return "StudentEnrollment"
}

// TableName overrides the default table name to "StudentDailyActivity".
func (StudentDailyActivity) TableName() string {
	return "StudentDailyActivity"
}

// StudentEnrollmentRequest captures parameters to register or drop a course.
type StudentEnrollmentRequest struct {
	CourseID uint `json:"course_id" binding:"required"`
}
