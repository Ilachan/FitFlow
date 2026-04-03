package main

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"my-course-backend/db"
	"my-course-backend/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// seedRoles ensures default roles exist in the database.
func seedRoles() {
	roles := []string{"Admin", "Teacher", "Student"}
	for _, roleName := range roles {
		var role model.Role
		err := db.DB.Where("role_name = ?", roleName).First(&role).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Skip role seed for %s: query failed (%v)", roleName, err)
			continue
		}
		if err := db.DB.Create(&model.Role{RoleName: roleName}).Error; err != nil {
			log.Printf("Role seed warning: failed to create %s (%v)", roleName, err)
			continue
		}
		log.Printf("Created role: %s", roleName)
	}
}

// seedUsers creates fake users by role up to the provided targets.
func seedUsers(adminTarget int, teacherTarget int, studentTarget int) {
	if adminTarget < 0 || teacherTarget < 0 || studentTarget < 0 {
		return
	}

	type roleTarget struct {
		name   string
		target int
	}

	roles := []roleTarget{
		{name: "Admin", target: adminTarget},
		{name: "Teacher", target: teacherTarget},
		{name: "Student", target: studentTarget},
	}

	defaultPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Skip User seed: failed to hash default password (%v)", err)
		return
	}

	for _, rt := range roles {
		if rt.target == 0 {
			continue
		}

		// Pick one matching role row for new inserts (older rows may also exist).
		var role model.Role
		if err := db.DB.Where("role_name = ?", rt.name).First(&role).Error; err != nil {
			log.Printf("Skip %s user seed: role not found (%v)", rt.name, err)
			continue
		}

		var existing int64
		if err := db.DB.Table("User").
			Joins("JOIN Role ON Role.id = User.role_id").
			Where("Role.role_name = ?", rt.name).
			Count(&existing).Error; err != nil {
			log.Printf("Skip %s user seed: failed to count rows (%v)", rt.name, err)
			continue
		}
		if existing >= int64(rt.target) {
			log.Printf("%s user seed skipped: already has %d rows (target=%d)", rt.name, existing, rt.target)
			continue
		}

		toCreate := rt.target - int(existing)
		created := 0

		attempts := 0
		maxAttempts := toCreate * 20
		nextIndex := 1
		for created < toCreate && attempts < maxAttempts {
			attempts++
			email := strings.ToLower(rt.name) + "." + time.Now().Format("20060102") + "." + formatIndex(nextIndex) + "@fitflow.local"
			name := rt.name + " User " + formatIndex(nextIndex)
			avatarSeed := formatIndex(nextIndex)
			nextIndex++

			var user model.User
			tx := db.DB.Where("email = ?", email).
				Attrs(model.User{
					Name:      name,
					Email:     email,
					Password:  string(defaultPassword),
					AvatarURL: "https://picsum.photos/seed/" + strings.ToLower(rt.name) + avatarSeed + "/200/200",
					RoleID:    role.ID,
				}).
				FirstOrCreate(&user)
			if tx.Error != nil {
				log.Printf("%s user seed warning: failed to create (%v)", rt.name, tx.Error)
				continue
			}
			if tx.RowsAffected > 0 {
				created++
			}
		}
		if created < toCreate {
			log.Printf("%s user seed partial: created %d/%d rows", rt.name, created, toCreate)
		}

		log.Printf("Seeded %s users: created %d rows", rt.name, created)
	}
}

// seedCourses creates fake classes up to targetCount.
func seedCourses(targetCount int) {
	if targetCount <= 0 {
		return
	}

	var existing int64
	if err := db.DB.Model(&model.Course{}).Count(&existing).Error; err != nil {
		log.Printf("Skip Course seed: failed to count rows (%v)", err)
		return
	}
	if existing >= int64(targetCount) {
		log.Printf("Course seed skipped: already has %d rows (target=%d)", existing, targetCount)
		return
	}

	type courseTemplate struct {
		name     string
		category string
		weekday  string
		start    string
		end      string
		capacity int
		duration int
	}

	templates := []courseTemplate{
		{name: "Morning Yoga", category: "Mind & Body", weekday: "Monday", start: "07:00", end: "08:00", capacity: 20, duration: 60},
		{name: "HIIT Express", category: "Cardio", weekday: "Tuesday", start: "18:00", end: "18:45", capacity: 16, duration: 45},
		{name: "Pilates Core", category: "Core", weekday: "Wednesday", start: "17:30", end: "18:30", capacity: 18, duration: 60},
		{name: "Strength Basics", category: "Strength", weekday: "Thursday", start: "19:00", end: "20:00", capacity: 14, duration: 60},
		{name: "Spin Burn", category: "Cardio", weekday: "Friday", start: "06:30", end: "07:15", capacity: 24, duration: 45},
		{name: "Zumba Beats", category: "Dance", weekday: "Saturday", start: "10:00", end: "11:00", capacity: 25, duration: 60},
		{name: "Mobility Flow", category: "Recovery", weekday: "Sunday", start: "09:00", end: "09:45", capacity: 22, duration: 45},
		{name: "Boxing Fit", category: "Conditioning", weekday: "Monday", start: "20:00", end: "21:00", capacity: 15, duration: 60},
		{name: "Power Circuit", category: "Strength", weekday: "Tuesday", start: "12:00", end: "12:50", capacity: 16, duration: 50},
		{name: "Sunset Stretch", category: "Recovery", weekday: "Wednesday", start: "20:15", end: "21:00", capacity: 20, duration: 45},
		{name: "Beginner Run Club", category: "Cardio", weekday: "Thursday", start: "06:00", end: "06:50", capacity: 30, duration: 50},
		{name: "Functional Training", category: "Conditioning", weekday: "Friday", start: "17:00", end: "18:00", capacity: 18, duration: 60},
	}

	toCreate := targetCount - int(existing)
	created := 0

	for i := 0; i < len(templates) && created < toCreate; i++ {
		t := templates[i]
		start, err := model.ParseTimeOnly(t.start)
		if err != nil {
			log.Printf("Course seed warning: invalid start time %q (%v)", t.start, err)
			continue
		}
		end, err := model.ParseTimeOnly(t.end)
		if err != nil {
			log.Printf("Course seed warning: invalid end time %q (%v)", t.end, err)
			continue
		}

		code := "FIT" + formatIndex(i+1)
		var course model.Course
		tx := db.DB.Where("course_code = ?", code).
			Attrs(model.Course{
				CourseName:  t.name,
				CourseCode:  code,
				Description: t.name + " class for all levels.",
				StartTime:   start,
				EndTime:     end,
				Capacity:    t.capacity,
				Duration:    t.duration,
				Category:    t.category,
				Weekday:     t.weekday,
			}).
			FirstOrCreate(&course)
		if tx.Error != nil {
			log.Printf("Course seed warning: failed to create %s (%v)", code, tx.Error)
			continue
		}
		if tx.RowsAffected > 0 {
			created++
		}
	}

	log.Printf("Seeded Course: created %d rows", created)
}

func formatIndex(i int) string {
	if i < 10 {
		return "0" + strconv.Itoa(i)
	}
	return strconv.Itoa(i)
}

// seedEnrollments generates fake data in Enrollment.
// It is idempotent and only creates missing rows until targetCount is reached.
func seedEnrollments(targetCount int) {
	if targetCount <= 0 {
		return
	}

	var users []model.User
	if err := db.DB.Order("id ASC").Find(&users).Error; err != nil {
		log.Printf("Skip Enrollment seed: failed to load users (%v)", err)
		return
	}

	var courses []model.Course
	if err := db.DB.Order("id ASC").Find(&courses).Error; err != nil {
		log.Printf("Skip Enrollment seed: failed to load courses (%v)", err)
		return
	}

	if len(users) == 0 || len(courses) == 0 {
		log.Printf("Skip Enrollment seed: requires at least 1 user and 1 course (users=%d, courses=%d)", len(users), len(courses))
		return
	}

	if err := ensureAnalyticsCoverageEnrollments(courses); err != nil {
		log.Printf("Enrollment seed warning: analytics coverage setup failed (%v)", err)
	}

	var existing int64
	if err := db.DB.Model(&model.Enrollment{}).Count(&existing).Error; err != nil {
		log.Printf("Skip Enrollment seed: failed to count existing rows (%v)", err)
		return
	}
	if existing >= int64(targetCount) {
		log.Printf("Enrollment seed skipped: already has %d rows (target=%d)", existing, targetCount)
		return
	}

	toCreate := targetCount - int(existing)
	statuses := []string{
		model.EnrollmentStatusEnrolled,
		model.EnrollmentStatusEnrolled,
		model.EnrollmentStatusAttended,
		model.EnrollmentStatusMissed,
		model.EnrollmentStatusEnrolled,
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	lookbackDays := 35

	created := 0
	attempts := 0
	maxAttempts := toCreate * 20

	for created < toCreate && attempts < maxAttempts {
		attempts++

		user := users[rng.Intn(len(users))]
		course := courses[rng.Intn(len(courses))]
		status := statuses[rng.Intn(len(statuses))]
		dayOffset := rng.Intn(lookbackDays)
		minuteOfDay := rng.Intn(24 * 60)
		enrollTime := time.Now().
			AddDate(0, 0, -dayOffset).
			Truncate(24 * time.Hour).
			Add(time.Duration(minuteOfDay) * time.Minute)

		var enrollment model.Enrollment
		tx := db.DB.Where(&model.Enrollment{UserID: user.ID, CourseID: course.ID}).
			Attrs(model.Enrollment{Status: status, EnrollTime: enrollTime}).
			FirstOrCreate(&enrollment)
		if tx.Error != nil {
			log.Printf("Enrollment seed warning: failed to create (%v)", tx.Error)
			continue
		}

		if tx.RowsAffected > 0 {
			created++
		}
	}

	if created < toCreate {
		log.Printf("Enrollment seed partial: created %d/%d rows (likely reached unique pair limit)", created, toCreate)
		return
	}

	log.Printf("Seeded Enrollment: created %d fake rows", created)
}

// ensureAnalyticsCoverageEnrollments guarantees one student has rows across 7d and 1m windows.
func ensureAnalyticsCoverageEnrollments(courses []model.Course) error {
	var student model.User
	err := db.DB.Table("User").
		Joins("JOIN Role ON Role.id = User.role_id").
		Where("Role.role_name = ?", "Student").
		Order("User.id ASC").
		First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	dayOffsets := []int{0, 1, 2, 3, 5, 7, 9, 12, 15, 18, 21, 24, 27, 30}
	limit := len(courses)
	if limit > len(dayOffsets) {
		limit = len(dayOffsets)
	}

	created := 0
	updated := 0
	for i := 0; i < limit; i++ {
		course := courses[i]
		enrollTime := time.Now().
			AddDate(0, 0, -dayOffsets[i]).
			Truncate(24 * time.Hour).
			Add(time.Duration(8+i%10) * time.Hour).
			Add(time.Duration(i*7) * time.Minute)

		var enrollment model.Enrollment
		findErr := db.DB.Where("user_id = ? AND course_id = ?", student.ID, course.ID).First(&enrollment).Error
		if findErr != nil {
			if !errors.Is(findErr, gorm.ErrRecordNotFound) {
				return findErr
			}

			enrollment = model.Enrollment{
				UserID:     student.ID,
				CourseID:   course.ID,
				Status:     model.EnrollmentStatusEnrolled,
				EnrollTime: enrollTime,
			}
			if err := db.DB.Create(&enrollment).Error; err != nil {
				return err
			}
			created++
			continue
		}

		if err := db.DB.Model(&enrollment).Updates(map[string]any{
			"status":      model.EnrollmentStatusEnrolled,
			"enroll_time": enrollTime,
		}).Error; err != nil {
			return err
		}
		updated++
	}

	log.Printf("Enrollment analytics coverage: ensured %d rows for student_id=%d (created=%d, updated=%d)", limit, student.ID, created, updated)
	return nil
}
