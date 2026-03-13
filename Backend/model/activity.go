package model

import "time"

// UserDailyActivity stores one row per user's class activity on a date.
type UserDailyActivity struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	EnrollmentID uint `gorm:"column:enrollment_id;not null;index:idx_sda_enrollment_date,unique" json:"enrollment_id"`
	UserID       uint `gorm:"column:user_id;not null;index" json:"user_id"`
	CourseID     uint `gorm:"column:course_id;not null;index" json:"course_id"`

	ActivityDate time.Time `gorm:"column:activity_date;type:date;not null;index:idx_sda_enrollment_date,unique;index" json:"activity_date"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (UserDailyActivity) TableName() string {
	return "UserDailyActivity"
}

// DailyActivitySummary is grouped daily analytics for dashboard charts.
type DailyActivitySummary struct {
	Date    string `json:"date"`
	Classes int64  `json:"classes"`
}

// CategoryActivitySummary is grouped category analytics for dashboard charts.
type CategoryActivitySummary struct {
	Category   string  `json:"category"`
	Classes    int64   `json:"classes"`
	Percentage float64 `json:"percentage"`
}

// UserAnalyticsResponse is the API response for user dashboard analytics.
type UserAnalyticsResponse struct {
	UserID       uint                      `json:"user_id"`
	Range        string                    `json:"range"`
	FromDate     string                    `json:"from_date"`
	ToDate       string                    `json:"to_date"`
	TotalClasses int64                     `json:"total_classes"`
	TotalTime    int64                     `json:"total_time"`
	ActiveDays   int64                     `json:"active_days"`
	Daily        []DailyActivitySummary    `json:"daily"`
	Categories   []CategoryActivitySummary `json:"categories"`
}
