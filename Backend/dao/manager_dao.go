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