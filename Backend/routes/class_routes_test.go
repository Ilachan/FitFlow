package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"my-course-backend/db"
	"my-course-backend/model"
	"my-course-backend/routes"
	"my-course-backend/service"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var routeTestSeq uint64

func setupRouteTestDB(t *testing.T) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	dsn := fmt.Sprintf("file:route_test_%d?mode=memory&cache=shared", time.Now().UnixNano())
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

func seedRouteRole(t *testing.T, roleID uint, roleName string) model.Role {
	t.Helper()

	role := model.Role{ID: roleID, RoleName: roleName}
	if err := db.DB.Create(&role).Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}

	return role
}

func seedRouteUser(t *testing.T, roleID uint, password string) model.User {
	t.Helper()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := model.User{
		Name:     "Route Test User",
		Email:    fmt.Sprintf("route-user-%d-%d@example.com", time.Now().UnixNano(), atomic.AddUint64(&routeTestSeq, 1)),
		Password: string(hashedPassword),
		RoleID:   roleID,
	}
	if err := db.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	return user
}

func seedRouteCourse(t *testing.T, name string, capacity int, category string) model.Course {
	t.Helper()
	start := time.Now().Add(2 * time.Hour)
	end := start.Add(1 * time.Hour)
	return seedRouteCourseWithSchedule(t, name, capacity, category, start.Format("15:04"), end.Format("15:04"), start.Weekday().String())
}

func seedRouteCourseWithSchedule(t *testing.T, name string, capacity int, category string, start string, end string, weekday string) model.Course {
	t.Helper()

	startTime, err := model.ParseTimeOnly(start)
	if err != nil {
		t.Fatalf("failed to parse start time: %v", err)
	}
	endTime, err := model.ParseTimeOnly(end)
	if err != nil {
		t.Fatalf("failed to parse end time: %v", err)
	}

	course := model.Course{
		CourseName: name,
		CourseCode: fmt.Sprintf("R-%d-%d", time.Now().UnixNano(), atomic.AddUint64(&routeTestSeq, 1)),
		Capacity:   capacity,
		Category:   category,
		StartTime:  startTime,
		EndTime:    endTime,
		Weekday:    weekday,
	}
	if err := db.DB.Create(&course).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}

	seedRouteCourseSession(t, course)

	return course
}

func seedRouteCourseSession(t *testing.T, course model.Course) model.ClassSession {
	t.Helper()

	sessionDate := time.Now().AddDate(0, 0, 1)
	session := model.ClassSession{
		CourseID:    course.ID,
		SessionDate: sessionDate.Format("2006-01-02"),
		StartAt:     time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), course.StartTime.Hour(), course.StartTime.Minute(), 0, 0, sessionDate.Location()),
		EndAt:       time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(), course.EndTime.Hour(), course.EndTime.Minute(), 0, 0, sessionDate.Location()),
		Status:      "scheduled",
		Capacity:    course.Capacity,
	}
	if err := db.DB.Create(&session).Error; err != nil {
		t.Fatalf("failed to seed class session: %v", err)
	}

	return session
}

func seedRouteEnrollmentAt(t *testing.T, userID uint, courseID uint, status string, enrollTime time.Time) model.Enrollment {
	t.Helper()

	var session model.ClassSession
	if err := db.DB.Where("course_id = ? AND status = ?", courseID, "scheduled").Order("session_date ASC").First(&session).Error; err != nil {
		t.Fatalf("failed to find scheduled session: %v", err)
	}

	enrollment := model.Enrollment{
		UserID:     userID,
		CourseID:   courseID,
		SessionID:  &session.ID,
		Status:     status,
		EnrollTime: enrollTime,
	}
	if err := db.DB.Create(&enrollment).Error; err != nil {
		t.Fatalf("failed to seed enrollment: %v", err)
	}

	return enrollment
}

func issueRouteToken(t *testing.T, email string, password string) string {
	t.Helper()

	token, _, err := service.LoginUserWithRole(model.LoginInput{
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Fatalf("failed to issue token: %v", err)
	}

	return token
}

func performJSONRequest(t *testing.T, router http.Handler, method string, path string, token string, payload any) *httptest.ResponseRecorder {
	t.Helper()

	var body bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			t.Fatalf("failed to encode request payload: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, &body)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func TestRegisterClassEndpoint_Created(t *testing.T) {
	setupRouteTestDB(t)
	seedRouteRole(t, 1, "Student")
	user := seedRouteUser(t, 1, "secret123")
	course := seedRouteCourse(t, "Yoga", 2, "Wellness")
	token := issueRouteToken(t, user.Email, "secret123")
	router := routes.SetupRouter()

	recorder := performJSONRequest(t, router, http.MethodPost, "/classes/register", token, map[string]uint{"course_id": course.ID})
	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Message != "Class enrolled successfully" {
		t.Fatalf("unexpected message: %s", response.Message)
	}
}

func TestRegisterClassEndpoint_Conflict(t *testing.T) {
	setupRouteTestDB(t)
	seedRouteRole(t, 1, "Student")
	user := seedRouteUser(t, 1, "secret123")
	course := seedRouteCourse(t, "Spin", 3, "Cardio")
	seedRouteEnrollmentAt(t, user.ID, course.ID, model.EnrollmentStatusEnrolled, time.Now())
	token := issueRouteToken(t, user.Email, "secret123")
	router := routes.SetupRouter()

	recorder := performJSONRequest(t, router, http.MethodPost, "/classes/register", token, map[string]uint{"course_id": course.ID})
	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Error != "enrollment already exists" {
		t.Fatalf("unexpected error: %s", response.Error)
	}
}

func TestRegisterClassEndpoint_ConflictOnScheduleOverlap(t *testing.T) {
	setupRouteTestDB(t)
	seedRouteRole(t, 1, "Student")
	user := seedRouteUser(t, 1, "secret123")
	base := time.Now().Add(2 * time.Hour)
	weekday := base.Weekday().String()
	existingCourse := seedRouteCourseWithSchedule(t, "Morning Yoga", 3, "Wellness", base.Format("15:04"), base.Add(1*time.Hour).Format("15:04"), weekday)
	targetCourse := seedRouteCourseWithSchedule(t, "Strength Flow", 3, "Strength", base.Add(30*time.Minute).Format("15:04"), base.Add(90*time.Minute).Format("15:04"), weekday[:3])
	seedRouteEnrollmentAt(t, user.ID, existingCourse.ID, model.EnrollmentStatusEnrolled, time.Now())
	token := issueRouteToken(t, user.Email, "secret123")
	router := routes.SetupRouter()

	recorder := performJSONRequest(t, router, http.MethodPost, "/classes/register", token, map[string]uint{"course_id": targetCourse.ID})
	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Error != "class schedule overlaps with an existing enrolled class" {
		t.Fatalf("unexpected error: %s", response.Error)
	}
}

func TestDropClassEndpoint_OK(t *testing.T) {
	setupRouteTestDB(t)
	seedRouteRole(t, 1, "Student")
	user := seedRouteUser(t, 1, "secret123")
	course := seedRouteCourse(t, "HIIT", 4, "Cardio")
	seedRouteEnrollmentAt(t, user.ID, course.ID, model.EnrollmentStatusEnrolled, time.Now())
	token := issueRouteToken(t, user.Email, "secret123")
	router := routes.SetupRouter()

	recorder := performJSONRequest(t, router, http.MethodPost, "/classes/drop", token, map[string]uint{"course_id": course.ID})
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Message != "Class unenrolled successfully" {
		t.Fatalf("unexpected message: %s", response.Message)
	}
}

func TestDropClassEndpoint_NotFound(t *testing.T) {
	setupRouteTestDB(t)
	seedRouteRole(t, 1, "Student")
	user := seedRouteUser(t, 1, "secret123")
	course := seedRouteCourse(t, "Pilates", 4, "Core")
	token := issueRouteToken(t, user.Email, "secret123")
	router := routes.SetupRouter()

	recorder := performJSONRequest(t, router, http.MethodPost, "/classes/drop", token, map[string]uint{"course_id": course.ID})
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Error != "enrollment not found" {
		t.Fatalf("unexpected error: %s", response.Error)
	}
}

func TestListClassesEndpoint_OK(t *testing.T) {
	setupRouteTestDB(t)
	seedRouteRole(t, 1, "Student")
	user := seedRouteUser(t, 1, "secret123")
	courseA := seedRouteCourse(t, "Course A", 3, "Cardio")
	courseB := seedRouteCourse(t, "Course B", 2, "Strength")
	seedRouteEnrollmentAt(t, user.ID, courseA.ID, model.EnrollmentStatusAttended, time.Now())
	router := routes.SetupRouter()

	recorder := performJSONRequest(t, router, http.MethodGet, "/classes?page=1", "", nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Page     int            `json:"page"`
		PageSize int            `json:"page_size"`
		Total    int64          `json:"total"`
		Classes  []model.Course `json:"classes"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Page != 1 {
		t.Fatalf("expected page 1, got %d", response.Page)
	}
	if response.PageSize != 20 {
		t.Fatalf("expected page size 20, got %d", response.PageSize)
	}
	if response.Total != 2 {
		t.Fatalf("expected total 2, got %d", response.Total)
	}
	if len(response.Classes) != 2 {
		t.Fatalf("expected 2 classes, got %d", len(response.Classes))
	}

	spots := map[uint]int{}
	for _, class := range response.Classes {
		spots[class.ID] = class.Spot
	}

	if spots[courseA.ID] != 2 {
		t.Fatalf("expected Course A spot 2, got %d", spots[courseA.ID])
	}
	if spots[courseB.ID] != 2 {
		t.Fatalf("expected Course B spot 2, got %d", spots[courseB.ID])
	}
}

func TestGetUserAnalyticsEndpoint_OK(t *testing.T) {
	setupRouteTestDB(t)
	seedRouteRole(t, 1, "Student")
	user := seedRouteUser(t, 1, "secret123")
	courseA := seedRouteCourse(t, "Run", 10, "Cardio")
	courseB := seedRouteCourse(t, "Stretch", 10, "")
	if err := db.DB.Model(&model.Course{}).Where("id = ?", courseA.ID).Update("duration", 45).Error; err != nil {
		t.Fatalf("failed to set courseA duration: %v", err)
	}
	if err := db.DB.Model(&model.Course{}).Where("id = ?", courseB.ID).Update("duration", 30).Error; err != nil {
		t.Fatalf("failed to set courseB duration: %v", err)
	}
	now := time.Now()
	seedRouteEnrollmentAt(t, user.ID, courseA.ID, model.EnrollmentStatusAttended, now.AddDate(0, 0, -1))
	seedRouteEnrollmentAt(t, user.ID, courseB.ID, model.EnrollmentStatusMissed, now.AddDate(0, 0, -2))
	token := issueRouteToken(t, user.Email, "secret123")
	router := routes.SetupRouter()

	path := fmt.Sprintf("/users/%d/analytics?range=7d", user.ID)
	recorder := performJSONRequest(t, router, http.MethodGet, path, token, nil)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Analytics model.UserAnalyticsResponse `json:"analytics"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Analytics.UserID != user.ID {
		t.Fatalf("expected user id %d, got %d", user.ID, response.Analytics.UserID)
	}
	if response.Analytics.TotalClasses != 1 {
		t.Fatalf("expected total classes 1, got %d", response.Analytics.TotalClasses)
	}
	if response.Analytics.TotalTime != 45 {
		t.Fatalf("expected total time 45, got %d", response.Analytics.TotalTime)
	}
	if response.Analytics.ActiveDays != 1 {
		t.Fatalf("expected active days 1, got %d", response.Analytics.ActiveDays)
	}
	if response.Analytics.Range != "7d" {
		t.Fatalf("expected range 7d, got %s", response.Analytics.Range)
	}
	if len(response.Analytics.Categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(response.Analytics.Categories))
	}
}

func TestGetUserAnalyticsEndpoint_Forbidden(t *testing.T) {
	setupRouteTestDB(t)
	seedRouteRole(t, 1, "Student")
	user := seedRouteUser(t, 1, "secret123")
	otherUser := seedRouteUser(t, 1, "secret456")
	token := issueRouteToken(t, otherUser.Email, "secret456")
	router := routes.SetupRouter()

	path := fmt.Sprintf("/users/%d/analytics?range=7d", user.ID)
	recorder := performJSONRequest(t, router, http.MethodGet, path, token, nil)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Error != "forbidden" {
		t.Fatalf("unexpected error: %s", response.Error)
	}
}
