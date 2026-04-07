package main

import (
	"log"

	"my-course-backend/db"
	"my-course-backend/model"
	"my-course-backend/routes"
)

func main() {
	// 1. Initialize Database
	db.InitDB()

	// 2. Auto migrate only missing tables to avoid SQLite temp-table rebuild issues.
	ensure := []struct {
		tableName string
		model     any
	}{
		{tableName: model.Role{}.TableName(), model: &model.Role{}},
		{tableName: model.User{}.TableName(), model: &model.User{}},
		{tableName: model.UserInfo{}.TableName(), model: &model.UserInfo{}},
		{tableName: model.Course{}.TableName(), model: &model.Course{}},
		{tableName: model.Enrollment{}.TableName(), model: &model.Enrollment{}},
		{tableName: model.UserDailyActivity{}.TableName(), model: &model.UserDailyActivity{}},
	}

	for _, item := range ensure {
		if db.DB.Migrator().HasTable(item.tableName) {
			continue
		}
		if err := db.DB.AutoMigrate(item.model); err != nil {
			log.Fatalf("auto-migrate failed for %s: %v", item.tableName, err)
		}
	}

	// 3. Seed Initial Data
	seedRoles()

	// 4. Initialize Router
	r := routes.SetupRouter()

	// 5. Start Server
	r.Run(":8080")
}

func seedRoles() {
	type fixedRole struct {
		ID   uint
		Name string
	}

	roles := []fixedRole{
		{ID: 1, Name: "Student"},
		{ID: 2, Name: "SuperManager"},
		{ID: 3, Name: "Manager"},
		{ID: 4, Name: "Instructor"},
	}

	for _, r := range roles {
		var role model.Role

		// 1) If there is already a row with this ID, make sure its name is correct.
		if err := db.DB.First(&role, r.ID).Error; err == nil {
			if role.RoleName != r.Name {
				role.RoleName = r.Name
				if err := db.DB.Save(&role).Error; err != nil {
					log.Printf("Failed to update role id=%d name=%s: %v", r.ID, r.Name, err)
				} else {
					log.Printf("Updated role id=%d to name=%s", r.ID, r.Name)
				}
			}
			continue
		}

		// 2) Otherwise, if there is a row with this name but different ID, move it to the fixed ID.
		if err := db.DB.Where("role_name = ?", r.Name).First(&role).Error; err == nil {
			// try to update primary key to fixed ID (works in SQLite; may vary by DB settings)
			if err := db.DB.Model(&model.Role{}).Where("id = ?", role.ID).Update("id", r.ID).Error; err != nil {
				// fallback: create the fixed row, leave old row (but unique role_name may block create)
				log.Printf("Failed to move role name=%s to id=%d (existing id=%d): %v", r.Name, r.ID, role.ID, err)
			} else {
				log.Printf("Moved role name=%s from id=%d to id=%d", r.Name, role.ID, r.ID)
			}
			continue
		}

		// 3) Create fixed row with specified ID.
		if err := db.DB.Create(&model.Role{ID: r.ID, RoleName: r.Name}).Error; err != nil {
			log.Printf("Failed to create role id=%d name=%s: %v", r.ID, r.Name, err)
		} else {
			log.Printf("Created role id=%d name=%s", r.ID, r.Name)
		}
	}
}
