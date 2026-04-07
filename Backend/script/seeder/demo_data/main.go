package main

import (
	"log"

	"my-course-backend/db"
	"my-course-backend/model"
)

func main() {
	// Keep schema creation together with seeding for local demo setup.
	db.InitDB()
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
		{tableName: model.ClassSession{}.TableName(), model: &model.ClassSession{}},
	}

	for _, item := range ensure {
		if db.DB.Migrator().HasTable(item.tableName) {
			continue
		}
		if err := db.DB.AutoMigrate(item.model); err != nil {
			log.Fatalf("auto-migrate failed for %s: %v", item.tableName, err)
		}
	}

	seedRoles()
	seedUsers(3, 8, 25)
	seedCourses(12)
	seedClassSessions()
	seedEnrollments(420)
	backfillEnrollmentSessionIDs()
}
