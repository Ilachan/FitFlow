package dao

import (
	"my-course-backend/model"

	"gorm.io/gorm"
)

func CreateManagerInviteCodeTx(tx *gorm.DB, invite *model.ManagerInviteCode) error {
	return tx.Model(&model.ManagerInviteCode{}).Create(invite).Error
}