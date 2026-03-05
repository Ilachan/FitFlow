package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"

	"my-course-backend/db"
	"my-course-backend/model"
)

type enrollmentRef struct {
	ID       uint
	CourseID uint
}

func main() {
	db.InitDB()
	db.DB.AutoMigrate(&model.Role{}, &model.Student{}, &model.Course{}, &model.StudentEnrollment{}, &model.StudentDailyActivity{})

	roleID, err := ensureStudentRole()
	if err != nil {
		log.Fatalf("failed to ensure Student role: %v", err)
	}

	courses, err := listCourses()
	if err != nil {
		log.Fatalf("failed to query courses: %v", err)
	}
	if len(courses) == 0 {
		log.Fatalf("no courses found in Course table. import courses first")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("SeedPass123!"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	studentIDs := make([]uint, 0, 10)
	for i := 1; i <= 10; i++ {
		name := fmt.Sprintf("Demo Student %02d", i)
		email := fmt.Sprintf("demo.student%02d@seed.local", i)

		studentID, created, err := ensureStudent(name, email, string(hash), roleID)
		if err != nil {
			log.Fatalf("failed to ensure student %s: %v", email, err)
		}
		if created {
			log.Printf("created student: %s (id=%d)", email, studentID)
		}
		studentIDs = append(studentIDs, studentID)
	}

	enrollmentsByStudent := make(map[uint][]enrollmentRef, len(studentIDs))
	createdEnrollments := 0
	for studentIndex, studentID := range studentIDs {
		for slot := 0; slot < 3; slot++ {
			course := courses[(studentIndex+slot)%len(courses)]
			enrollTime := time.Now().AddDate(0, 0, -(14 + studentIndex + slot)).Truncate(time.Second)

			enrollmentID, created, err := ensureEnrollment(studentID, course.ID, enrollTime)
			if err != nil {
				log.Fatalf("failed to ensure enrollment for student_id=%d course_id=%d: %v", studentID, course.ID, err)
			}
			if created {
				createdEnrollments++
			}

			enrollmentsByStudent[studentID] = append(enrollmentsByStudent[studentID], enrollmentRef{
				ID:       enrollmentID,
				CourseID: course.ID,
			})
		}
	}

	createdActivities := 0
	today := truncateToDate(time.Now())
	for studentIndex, studentID := range studentIDs {
		enrollments := enrollmentsByStudent[studentID]
		if len(enrollments) == 0 {
			continue
		}

		for dayOffset := 0; dayOffset < 7; dayOffset++ {
			activityDate := truncateToDate(today.AddDate(0, 0, -dayOffset))

			primary := enrollments[(studentIndex+dayOffset)%len(enrollments)]
			created, err := ensureDailyActivity(primary.ID, studentID, primary.CourseID, activityDate)
			if err != nil {
				log.Fatalf("failed to ensure daily activity (primary) for student_id=%d: %v", studentID, err)
			}
			if created {
				createdActivities++
			}

			if (studentIndex+dayOffset)%2 == 0 {
				secondary := enrollments[(studentIndex+dayOffset+1)%len(enrollments)]
				created, err := ensureDailyActivity(secondary.ID, studentID, secondary.CourseID, activityDate)
				if err != nil {
					log.Fatalf("failed to ensure daily activity (secondary) for student_id=%d: %v", studentID, err)
				}
				if created {
					createdActivities++
				}
			}
		}
	}

	log.Printf("fuller analytics seed complete: students=10, enrollments_created=%d, daily_activities_created=%d", createdEnrollments, createdActivities)
	log.Printf("demo login password for all seeded users: SeedPass123!")
}

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, t.Location())
}

func ensureStudentRole() (uint, error) {
	var role model.Role
	err := db.DB.Where("role_name = ?", "Student").Find(&role).Error
	if err != nil {
		return 0, err
	}
	if role.ID != 0 {
		return role.ID, nil
	}

	role = model.Role{RoleName: "Student"}
	if err := db.DB.Create(&role).Error; err != nil {
		return 0, err
	}
	return role.ID, nil
}

func listCourses() ([]model.Course, error) {
	var courses []model.Course
	err := db.DB.Order("id ASC").Find(&courses).Error
	return courses, err
}

func ensureStudent(name, email, password string, roleID uint) (uint, bool, error) {
	var student model.Student
	err := db.DB.Select("id").Where("email = ?", email).Find(&student).Error
	if err != nil {
		return 0, false, err
	}
	if student.ID != 0 {
		return student.ID, false, nil
	}

	student = model.Student{
		Name:     name,
		Email:    email,
		Password: password,
		RoleID:   roleID,
	}

	if err := db.DB.Create(&student).Error; err != nil {
		return 0, false, err
	}
	return student.ID, true, nil
}

func ensureEnrollment(studentID, courseID uint, enrollTime time.Time) (uint, bool, error) {
	var existing model.StudentEnrollment
	err := db.DB.Where("student_id = ? AND course_id = ?", studentID, courseID).Find(&existing).Error
	if err != nil {
		return 0, false, err
	}
	if existing.ID != 0 {
		if existing.Status != "registered" {
			existing.Status = "registered"
			if saveErr := db.DB.Save(&existing).Error; saveErr != nil {
				return 0, false, saveErr
			}
		}
		return existing.ID, false, nil
	}

	enrollment := model.StudentEnrollment{
		StudentID:  studentID,
		CourseID:   courseID,
		Status:     "registered",
		EnrollTime: enrollTime,
	}
	if err := db.DB.Create(&enrollment).Error; err != nil {
		return 0, false, err
	}
	return enrollment.ID, true, nil
}

func ensureDailyActivity(enrollmentID, studentID, courseID uint, activityDate time.Time) (bool, error) {
	activity := model.StudentDailyActivity{
		EnrollmentID: enrollmentID,
		StudentID:    studentID,
		CourseID:     courseID,
		ActivityDate: activityDate,
	}
	result := db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&activity)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
