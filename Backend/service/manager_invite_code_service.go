package service

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"strings"
	"time"

	"my-course-backend/dao"
	"my-course-backend/model"

	"gorm.io/gorm"
)

func generateInviteCode() (string, error) {
	// 8 bytes -> base32 about 13 chars, remove padding -> ~13 chars
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	code := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	code = strings.ToUpper(code)
	return code, nil
}

// CreateManagerInviteCode creates an invite code row and returns the generated code.
// Only SuperManager (role_id=2) should be allowed by API layer.
func CreateManagerInviteCode(inviterID uint, input model.CreateManagerInviteInput) (string, error) {
	if input.ExpireHours <= 0 {
		return "", errors.New("expire_hours must be greater than 0")
	}

	now := time.Now()
	expiredAt := now.Add(time.Duration(input.ExpireHours) * time.Hour)

	// Create within tx, retry a few times if code collision happens (unique constraint not shown, but safe to retry)
	var finalCode string
	err := dao.WithTx(func(tx *gorm.DB) error {
		for i := 0; i < 5; i++ {
			code, err := generateInviteCode()
			if err != nil {
				return err
			}

			status := "active"
			inviter := inviterID

			var inviteeEmailPtr *string
			if strings.TrimSpace(input.InviteeEmail) != "" {
				email := strings.ToLower(strings.TrimSpace(input.InviteeEmail))
				inviteeEmailPtr = &email
			}

			invite := model.ManagerInviteCode{
				Code:         code,
				InviterID:    &inviter,
				InviteeEmail: inviteeEmailPtr,
				Status:       &status,
				CreatedAt:    &now,
				ExpiredAt:    &expiredAt,
				UsedAt:       nil,
			}

			if err := dao.CreateManagerInviteCodeTx(tx, &invite); err != nil {
				// if collision, retry; otherwise fail fast
				// SQLite error text varies; simplest: retry on any error a few times
				continue
			}

			finalCode = code
			return nil
		}
		return errors.New("failed to generate invite code")
	})

	if err != nil {
		return "", err
	}
	return finalCode, nil
}