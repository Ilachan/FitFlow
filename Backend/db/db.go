package db

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// connect SQLite
	DB, err = gorm.Open(sqlite.Open("homework.db"), &gorm.Config{}) //自己的db文件路径
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB.Exec("PRAGMA foreign_keys = ON;")
	// Normalize TIME values to HH:MM:SS for consistent scanning.
	if DB.Migrator().HasTable("Course") {
		DB.Exec("UPDATE Course SET start_time = start_time || ':00' WHERE start_time IS NOT NULL AND length(start_time) = 5;")
		DB.Exec("UPDATE Course SET end_time = end_time || ':00' WHERE end_time IS NOT NULL AND length(end_time) = 5;")
	}

	log.Println("Database connected and foreign keys enabled.")
}
