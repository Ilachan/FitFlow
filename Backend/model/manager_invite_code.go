package model

import "time"

// ManagerInviteCode maps to the Manager_Invite_Code table in SQLite.
type ManagerInviteCode struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Code         string     `gorm:"column:code" json:"code"`
	InviterID    *uint      `gorm:"column:inviter_id" json:"inviter_id"`
	InviteeEmail *string    `gorm:"column:invitee_email" json:"invitee_email"`
	Status       *string    `gorm:"column:status" json:"status"`
	CreatedAt    *time.Time `gorm:"column:created_at" json:"created_at"`
	ExpiredAt    *time.Time `gorm:"column:expired_at;not null" json:"expired_at"`
	UsedAt       *time.Time `gorm:"column:used_at" json:"used_at"`
}

func (ManagerInviteCode) TableName() string {
	return "Manager_Invite_Code"
}

// ManagerRegisterInput is the request payload for manager registration.
type ManagerRegisterInput struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	InviteCode string `json:"invite_code" binding:"required"`
}

// CHANGED/NEW: input for creating invite codes (SuperManager only)
type CreateManagerInviteInput struct {
	// optional: bind invite to a specific email (recommended)
	InviteeEmail string `json:"invitee_email" binding:"omitempty,email"`

	// optional: inviter_id can be set from token, so frontend doesn't need to pass it
	// but if you want to allow explicit, keep it; otherwise remove it.
	// InviterID uint `json:"inviter_id"`
	
	// required: expiry time in hours from now, e.g., 24, 72, 168
	ExpireHours int `json:"expire_hours" binding:"required,min=1,max=720"` // up to 30 days
}