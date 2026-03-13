package dao

import (
	"errors"

	"my-course-backend/db"
	"my-course-backend/model"

	"gorm.io/gorm"
)

// GetRoleByName retrieves the role ID based on the role name.
func GetRoleByName(name string) (uint, error) {
	var role model.Role
	err := db.DB.Where("role_name = ?", name).First(&role).Error
	if err != nil {
		return 0, err
	}
	return role.ID, nil
}

// CheckEmailExist checks if the given email already exists.
func CheckEmailExist(email string) bool {
	var count int64
	db.DB.Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// CreateUser creates a new user.
func CreateUser(user *model.User) error {
	return db.DB.Create(user).Error
}

// GetUserByEmail retrieves user by email.
func GetUserByEmail(email string) (*model.User, error) {
	var user model.User

	err := db.DB.Where("email = ?", email).
		Preload("Role").
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserProfileByID fetches combined data from User and user_info.
func GetUserProfileByID(id uint) (*model.UserProfile, error) {
	var profile model.UserProfile

	err := db.DB.Table("User").
		Select("User.name, User.email, User.avatar_url, user_info.date_of_birth, user_info.gender, user_info.phone_number, user_info.address").
		Joins("left join user_info on user_info.user_id = User.id").
		Where("User.id = ?", id).
		Scan(&profile).Error

	return &profile, err
}

// UpdateUserProfile updates profile fields partially. (旧版：无法区分 undefined vs null)
func UpdateUserProfile(id uint, p model.UserProfile) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {

		userUpdates := map[string]interface{}{}

		if p.Name != nil {
			userUpdates["name"] = *p.Name
		}

		if p.AvatarURL != nil {
			userUpdates["avatar_url"] = *p.AvatarURL
		}

		if len(userUpdates) > 0 {
			if err := tx.Model(&model.User{}).
				Where("id = ?", id).
				Updates(userUpdates).Error; err != nil {
				return err
			}
		}

		var info model.UserInfo
		err := tx.Where("user_id = ?", id).First(&info).Error

		infoUpdates := map[string]interface{}{}

		if p.DateOfBirth != nil {
			infoUpdates["date_of_birth"] = *p.DateOfBirth
		}
		if p.Gender != nil {
			infoUpdates["gender"] = *p.Gender
		}
		if p.PhoneNumber != nil {
			infoUpdates["phone_number"] = *p.PhoneNumber
		}
		if p.Address != nil {
			infoUpdates["address"] = *p.Address
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if len(infoUpdates) == 0 {
				return nil
			}

			newInfo := model.UserInfo{UserID: id}
			if p.DateOfBirth != nil {
				newInfo.DateOfBirth = *p.DateOfBirth
			}
			if p.Gender != nil {
				newInfo.Gender = *p.Gender
			}
			if p.PhoneNumber != nil {
				newInfo.PhoneNumber = *p.PhoneNumber
			}
			if p.Address != nil {
				newInfo.Address = *p.Address
			}

			return tx.Create(&newInfo).Error
		}

		if len(infoUpdates) > 0 {
			return tx.Model(&info).Updates(infoUpdates).Error
		}

		return nil
	})
}

// New: UpdateUserProfilePatch distinguishes undefined vs null vs value.
func UpdateUserProfilePatch(id uint, p model.UserProfilePatch) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {

		/* -------------------- User table -------------------- */
		userUpdates := map[string]interface{}{}

		if p.Name.Set {
			if p.Name.Valid {
				userUpdates["name"] = p.Name.Value
			} else {
				userUpdates["name"] = nil // explicit null
			}
		}

		if p.AvatarURL.Set {
			if p.AvatarURL.Valid {
				userUpdates["avatar_url"] = p.AvatarURL.Value
			} else {
				userUpdates["avatar_url"] = nil
			}
		}

		if len(userUpdates) > 0 {
			if err := tx.Model(&model.User{}).
				Where("id = ?", id).
				Updates(userUpdates).Error; err != nil {
				return err
			}
		}

		/* -------------------- user_info table -------------------- */
		var info model.UserInfo
		err := tx.Where("user_id = ?", id).First(&info).Error

		infoUpdates := map[string]interface{}{}

		if p.DateOfBirth.Set {
			if p.DateOfBirth.Valid {
				infoUpdates["date_of_birth"] = p.DateOfBirth.Value
			} else {
				infoUpdates["date_of_birth"] = nil
			}
		}
		if p.Gender.Set {
			if p.Gender.Valid {
				infoUpdates["gender"] = p.Gender.Value
			} else {
				infoUpdates["gender"] = nil
			}
		}
		if p.PhoneNumber.Set {
			if p.PhoneNumber.Valid {
				infoUpdates["phone_number"] = p.PhoneNumber.Value
			} else {
				infoUpdates["phone_number"] = nil
			}
		}
		if p.Address.Set {
			if p.Address.Valid {
				infoUpdates["address"] = p.Address.Value
			} else {
				infoUpdates["address"] = nil
			}
		}

		// If user_info doesn't exist, create it when there is at least one field to set.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If all updates are "set to null" only, you have 2 choices:
			// - create a row with NULLs (allowed if columns nullable)
			// - OR do nothing
			// Here: if any field is Set (even null), we create the row.
			shouldCreate := p.DateOfBirth.Set || p.Gender.Set || p.PhoneNumber.Set || p.Address.Set
			if !shouldCreate {
				return nil
			}

			newInfo := model.UserInfo{UserID: id}

			// We can only assign concrete string values to struct fields; for null we rely on Updates(map) normally.
			// So: if user_info missing, we create it with empty strings for Valid=false (or keep defaults).
			// Better approach is to use map insert; but keep it simple:
			if p.DateOfBirth.Set && p.DateOfBirth.Valid {
				newInfo.DateOfBirth = p.DateOfBirth.Value
			}
			if p.Gender.Set && p.Gender.Valid {
				newInfo.Gender = p.Gender.Value
			}
			if p.PhoneNumber.Set && p.PhoneNumber.Valid {
				newInfo.PhoneNumber = p.PhoneNumber.Value
			}
			if p.Address.Set && p.Address.Valid {
				newInfo.Address = p.Address.Value
			}

			return tx.Create(&newInfo).Error
		}

		// Normal update
		if len(infoUpdates) > 0 {
			return tx.Model(&info).Updates(infoUpdates).Error
		}

		return nil
	})
}

// DeleteUserByID deletes a user.
func DeleteUserByID(id uint) error {
	result := db.DB.Delete(&model.User{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}