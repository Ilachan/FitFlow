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

	// 2. Auto Migrate Table Structures
	// CHANGED: &model.Student{} -> &model.User{}
	db.DB.AutoMigrate(&model.Role{}, &model.User{}, &model.Course{}, &model.StudentEnrollment{})

	// 3. Seed Initial Data
	seedRoles()

	// 4. Initialize Router
	r := routes.SetupRouter()

	// 5. Start Server
	r.Run(":8080")
}

// seedRoles ensures default roles exist in the database
func seedRoles() {
	roles := []string{"Admin", "Teacher", "Student"}
	for _, roleName := range roles {
		var role model.Role
		if err := db.DB.Where("role_name = ?", roleName).First(&role).Error; err != nil {
			db.DB.Create(&model.Role{RoleName: roleName})
			log.Printf("Created role: %s", roleName)
		}
	}
}