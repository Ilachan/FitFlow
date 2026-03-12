package service

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"my-course-backend/db"
	"my-course-backend/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var testSeq uint64

func setupClassServiceTestDB(t *testing.T) {
	t.Helper()

	dsn := fmt.Sprintf("file:test_%d?mode=memory&cache=shared", time.Now().UnixNano())
	testDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := testDB.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	if err := testDB.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.UserInfo{},
		&model.Course{},
		&model.Enrollment{},
		&model.UserDailyActivity{},
	); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	db.DB = testDB
}

func seedRoleAndUser(t *testing.T, roleID uint) model.User {
	t.Helper()

	role := model.Role{ID: roleID, RoleName: fmt.Sprintf("Role-%d", roleID)}
	if err := db.DB.Create(&role).Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}

	user := model.User{
		Name:     "User One",
		Email:    fmt.Sprintf("user-%d-%d@example.com", time.Now().UnixNano(), atomic.AddUint64(&testSeq, 1)),
		Password: "secret",
		RoleID:   roleID,
	}
	if err := db.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	return user
}

func seedCourse(t *testing.T, name string, capacity int, category string) model.Course {
	t.Helper()

	course := model.Course{
		CourseName: name,
		CourseCode: fmt.Sprintf("C-%d-%d", time.Now().UnixNano(), atomic.AddUint64(&testSeq, 1)),
		Capacity:   capacity,
		Category:   category,
	}
	if err := db.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	return course
}

func seedEnrollmentAt(t *testing.T, userID uint, courseID uint, status string, enrollTime time.Time) model.Enrollment {
	t.Helper()

	enrollment := model.Enrollment{
		UserID:     userID,
		CourseID:   courseID,
		Status:     status,
		EnrollTime: enrollTime,
	}
	if err := db.DB.Create(&enrollment).Error; err != nil {
		t.Fatalf("failed to seed enrollment: %v", err)
	}

	return enrollment
}

func TestRegisterClass_Success(t *testing.T) {
	setupClassServiceTestDB(t)

	user := seedRoleAndUser(t, 1)
	course := seedCourse(t, "Yoga", 3, "Wellness")

	if err := RegisterClass(user.ID, course.ID); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	var enrollments int64
	if err := db.DB.Model(&model.Enrollment{}).
		Where("user_id = ? AND course_id = ?", user.ID, course.ID).
		Count(&enrollments).Error; err != nil {
		t.Fatalf("failed to verify enrollment: %v", err)
	}
	if enrollments != 1 {
		t.Fatalf("expected 1 enrollment, got %d", enrollments)
	}

	var activities int64
	if err := db.DB.Model(&model.UserDailyActivity{}).
		Where("user_id = ?", user.ID).
		Count(&activities).Error; err != nil {
		t.Fatalf("failed to verify daily activity: %v", err)
	}
	if activities != 1 {
		t.Fatalf("expected 1 daily activity row, got %d", activities)
	}
}

func TestRegisterClass_UserNotFound(t *testing.T) {
	setupClassServiceTestDB(t)

	course := seedCourse(t, "Pilates", 2, "Core")

	err := RegisterClass(9999, course.ID)
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected user not found, got: %v", err)
	}
}

func TestRegisterClass_ClassNotFound(t *testing.T) {
	setupClassServiceTestDB(t)

	user := seedRoleAndUser(t, 1)

	err := RegisterClass(user.ID, 9999)
	if err == nil || err.Error() != "class not found" {
		t.Fatalf("expected class not found, got: %v", err)
	}
}

func TestRegisterClass_AlreadyExists(t *testing.T) {
	setupClassServiceTestDB(t)

	user := seedRoleAndUser(t, 1)
	course := seedCourse(t, "Spin", 5, "Cardio")
	seedEnrollmentAt(t, user.ID, course.ID, "registered", time.Now())

	err := RegisterClass(user.ID, course.ID)
	if err == nil || err.Error() != "enrollment already exists" {
		t.Fatalf("expected enrollment already exists, got: %v", err)
	}
}

func TestRegisterClass_ClassFull(t *testing.T) {
	setupClassServiceTestDB(t)

	user1 := seedRoleAndUser(t, 1)
	user2 := seedRoleAndUser(t, 2)
	course := seedCourse(t, "Boxing", 1, "Combat")
	seedEnrollmentAt(t, user1.ID, course.ID, "registered", time.Now())

	err := RegisterClass(user2.ID, course.ID)
	if err == nil || err.Error() != "class is full" {
		t.Fatalf("expected class is full, got: %v", err)
	}
}

func TestDropClass_Success(t *testing.T) {
	setupClassServiceTestDB(t)

	user := seedRoleAndUser(t, 1)
	course := seedCourse(t, "HIIT", 4, "Cardio")
	seedEnrollmentAt(t, user.ID, course.ID, "registered", time.Now())

	if err := DropClass(user.ID, course.ID); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	var enrollments int64
	if err := db.DB.Model(&model.Enrollment{}).
		Where("user_id = ? AND course_id = ?", user.ID, course.ID).
		Count(&enrollments).Error; err != nil {
		t.Fatalf("failed to verify enrollment deletion: %v", err)
	}
	if enrollments != 0 {
		t.Fatalf("expected 0 enrollments, got %d", enrollments)
	}
}

func TestDropClass_NotFound(t *testing.T) {
	setupClassServiceTestDB(t)

	err := DropClass(1, 1)
	if err == nil || err.Error() != "enrollment not found" {
		t.Fatalf("expected enrollment not found, got: %v", err)
	}
}

func TestListClassesPaged_ReturnsSpotAndPagination(t *testing.T) {
	setupClassServiceTestDB(t)

	user1 := seedRoleAndUser(t, 1)
	user2 := seedRoleAndUser(t, 2)
	courseA := seedCourse(t, "Course A", 3, "Cardio")
	courseB := seedCourse(t, "Course B", 2, "Strength")

	seedEnrollmentAt(t, user1.ID, courseA.ID, "registered", time.Now())
	seedEnrollmentAt(t, user2.ID, courseA.ID, "registered", time.Now())

	classes, total, err := ListClassesPaged(1, 10)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total classes 2, got %d", total)
	}
	if len(classes) != 2 {
		t.Fatalf("expected 2 classes in page, got %d", len(classes))
	}

	spots := map[uint]int{}
	for _, c := range classes {
		spots[c.ID] = c.Spot
	}

	if spots[courseA.ID] != 1 {
		t.Fatalf("expected course A spot 1, got %d", spots[courseA.ID])
	}
	if spots[courseB.ID] != 2 {
		t.Fatalf("expected course B spot 2, got %d", spots[courseB.ID])
	}
}

func TestGetUserAnalytics_SuccessWithPercentages(t *testing.T) {
	setupClassServiceTestDB(t)

	user := seedRoleAndUser(t, 1)
	courseCardio1 := seedCourse(t, "Run", 10, "Cardio")
	courseCardio2 := seedCourse(t, "Bike", 10, "Cardio")
	courseNoCategory := seedCourse(t, "Stretch", 10, "")

	now := time.Now()
	seedEnrollmentAt(t, user.ID, courseCardio1.ID, "registered", now.AddDate(0, 0, -1))
	seedEnrollmentAt(t, user.ID, courseCardio2.ID, "registered", now.AddDate(0, 0, -2))
	seedEnrollmentAt(t, user.ID, courseNoCategory.ID, "registered", now.AddDate(0, 0, -3))

	analytics, err := GetUserAnalytics(user.ID, "7d")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if analytics.TotalClasses != 3 {
		t.Fatalf("expected total classes 3, got %d", analytics.TotalClasses)
	}
	if analytics.ActiveDays != 3 {
		t.Fatalf("expected active days 3, got %d", analytics.ActiveDays)
	}
	if analytics.Range != "7d" {
		t.Fatalf("expected range 7d, got %s", analytics.Range)
	}
	if len(analytics.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(analytics.Categories))
	}

	categoryPct := map[string]float64{}
	categoryClasses := map[string]int64{}
	for _, c := range analytics.Categories {
		categoryPct[c.Category] = c.Percentage
		categoryClasses[c.Category] = c.Classes
	}

	if categoryClasses["Cardio"] != 2 {
		t.Fatalf("expected Cardio classes 2, got %d", categoryClasses["Cardio"])
	}
	if categoryPct["Cardio"] != 66.67 {
		t.Fatalf("expected Cardio percentage 66.67, got %.2f", categoryPct["Cardio"])
	}
	if categoryClasses["Uncategorized"] != 1 {
		t.Fatalf("expected Uncategorized classes 1, got %d", categoryClasses["Uncategorized"])
	}
	if categoryPct["Uncategorized"] != 33.33 {
		t.Fatalf("expected Uncategorized percentage 33.33, got %.2f", categoryPct["Uncategorized"])
	}
	if len(analytics.Daily) != 3 {
		t.Fatalf("expected 3 daily summary rows, got %d", len(analytics.Daily))
	}
}

func TestGetUserAnalytics_UserNotFound(t *testing.T) {
	setupClassServiceTestDB(t)

	analytics, err := GetUserAnalytics(9999, "7d")
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected user not found, got analytics=%v err=%v", analytics, err)
	}
}
