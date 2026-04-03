package db

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dbPath := ResolvedDBPath()
	// connect SQLite
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB.Exec("PRAGMA foreign_keys = ON;")
	migrateUserInfoTable()
	migrateEnrollmentTable()
	ensureUserDailyActivityTable()
	// Normalize TIME values to HH:MM:SS for consistent scanning.
	if DB.Migrator().HasTable("Course") {
		DB.Exec("UPDATE Course SET start_time = start_time || ':00' WHERE start_time IS NOT NULL AND length(start_time) = 5;")
		DB.Exec("UPDATE Course SET end_time = end_time || ':00' WHERE end_time IS NOT NULL AND length(end_time) = 5;")
	}

	log.Printf("Database connected and foreign keys enabled. Using %s", dbPath)
}

func ResolvedDBPath() string {
	if configured := strings.TrimSpace(os.Getenv("FITFLOW_DB_PATH")); configured != "" {
		return absPathOrFallback(configured)
	}

	if backendRoot := backendRootFromSource(); backendRoot != "" {
		return filepath.Join(backendRoot, "homework.db")
	}

	return resolveDBPathFromWorkingDirectory()
}

func backendRootFromSource() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	backendRoot := filepath.Clean(filepath.Join(filepath.Dir(file), ".."))
	return absPathOrFallback(backendRoot)
}

func resolveDBPathFromWorkingDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		return absPathOrFallback("homework.db")
	}

	current := wd
	for {
		goModPath := filepath.Join(current, "go.mod")
		if stat, statErr := os.Stat(goModPath); statErr == nil && !stat.IsDir() {
			return absPathOrFallback(filepath.Join(current, "homework.db"))
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return absPathOrFallback("homework.db")
}

func absPathOrFallback(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return absPath
}

func migrateUserInfoTable() {
	if DB == nil {
		return
	}

	hasStudentInfo := DB.Migrator().HasTable("student_info")
	hasUserInfo := DB.Migrator().HasTable("user_info")

	if hasStudentInfo && !hasUserInfo {
		if err := DB.Migrator().RenameTable("student_info", "user_info"); err != nil {
			log.Printf("Failed to rename table student_info to user_info: %v", err)
		}
	}

	if DB.Migrator().HasTable("user_info") && DB.Migrator().HasColumn("user_info", "student_id") && !DB.Migrator().HasColumn("user_info", "user_id") {
		if err := DB.Migrator().RenameColumn("user_info", "student_id", "user_id"); err != nil {
			log.Printf("Failed to rename column student_id to user_id in user_info: %v", err)
		}
	}
}

func migrateEnrollmentTable() {
	if DB == nil {
		return
	}

	hasStudentEnrollment := DB.Migrator().HasTable("StudentEnrollment")
	hasEnrollment := DB.Migrator().HasTable("Enrollment")

	if hasStudentEnrollment && !hasEnrollment {
		if err := DB.Migrator().RenameTable("StudentEnrollment", "Enrollment"); err != nil {
			log.Printf("Failed to rename table StudentEnrollment to Enrollment: %v", err)
		}
	}

	if DB.Migrator().HasTable("Enrollment") && DB.Migrator().HasColumn("Enrollment", "student_id") && !DB.Migrator().HasColumn("Enrollment", "user_id") {
		if err := DB.Migrator().RenameColumn("Enrollment", "student_id", "user_id"); err != nil {
			log.Printf("Failed to rename column student_id to user_id in Enrollment: %v", err)
		}
	}

	if DB.Migrator().HasTable("Enrollment") && DB.Migrator().HasColumn("Enrollment", "status") {
		if err := DB.Exec(`UPDATE Enrollment SET status = 'enrolled' WHERE status IN ('registered', 'pending')`).Error; err != nil {
			log.Printf("Failed to normalize Enrollment status values: %v", err)
		}
		if err := DB.Exec(`DELETE FROM Enrollment WHERE status = 'dropped'`).Error; err != nil {
			log.Printf("Failed to remove dropped Enrollment rows: %v", err)
		}
	}
}

func ensureUserDailyActivityTable() {
	if DB == nil {
		return
	}

	query := `
		CREATE TABLE IF NOT EXISTS "UserDailyActivity" (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			enrollment_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			course_id INTEGER NOT NULL,
			activity_date DATE NOT NULL,
			created_at DATETIME,
			UNIQUE(enrollment_id, activity_date),
			FOREIGN KEY (enrollment_id) REFERENCES "Enrollment"(id) ON UPDATE CASCADE ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES "User"(id) ON UPDATE CASCADE ON DELETE CASCADE,
			FOREIGN KEY (course_id) REFERENCES "Course"(id) ON UPDATE CASCADE ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_user_daily_activity_user_id ON "UserDailyActivity" (user_id);
		CREATE INDEX IF NOT EXISTS idx_user_daily_activity_course_id ON "UserDailyActivity" (course_id);
		CREATE INDEX IF NOT EXISTS idx_user_daily_activity_activity_date ON "UserDailyActivity" (activity_date);
	`

	if err := DB.Exec(query).Error; err != nil {
		log.Printf("Failed to ensure UserDailyActivity table exists: %v", err)
	}
}
