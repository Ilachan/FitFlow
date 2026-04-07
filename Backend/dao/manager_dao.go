package dao

import (
	"time"

	"my-course-backend/db"
	"my-course-backend/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetManagerInviteCodeForUpdate locks and returns an invite code row.
// Using a transaction + FOR UPDATE style locking (SQLite supports it via transaction locking semantics).
func GetManagerInviteCodeForUpdate(tx *gorm.DB, code string) (*model.ManagerInviteCode, error) {
	var invite model.ManagerInviteCode
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("code = ?", code).
		First(&invite).Error; err != nil {
		return nil, err
	}
	return &invite, nil
}

func MarkInviteCodeUsed(tx *gorm.DB, id uint, inviteeEmail string) error {
	now := time.Now()
	return tx.Model(&model.ManagerInviteCode{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        "used",
			"used_at":       now,
			"invitee_email": inviteeEmail,
		}).Error
}

// Expose DB for other DAOs if needed
func WithTx(fn func(tx *gorm.DB) error) error {
	return db.DB.Transaction(fn)
}

// ✅ Manager: 获取所有用户（分页）
func ListUsersPaged(limit int, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if err := db.DB.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.DB.
		Preload("Role").
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ✅ Manager: 查看用户已选课程
func ListEnrollmentsByUser(userID uint) ([]model.Enrollment, error) {
	var enrollments []model.Enrollment
	if err := db.DB.
		Where("user_id = ?", userID).
		Preload("Course").
		Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}
