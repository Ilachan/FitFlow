package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"my-course-backend/db"
	"my-course-backend/model"
)

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
		} else {
			log.Printf("existing student: %s (id=%d)", email, studentID)
		}
		studentIDs = append(studentIDs, studentID)
	}

	createdEnrollments := 0
	for idx, studentID := range studentIDs {
		course := courses[idx%len(courses)]
		enrollTime := time.Now().AddDate(0, 0, -idx).Truncate(time.Second)

		created, err := ensureEnrollment(studentID, course.ID, enrollTime)
		if err != nil {
			log.Fatalf("failed to ensure enrollment for student_id=%d: %v", studentID, err)
		}
		if created {
			createdEnrollments++
			log.Printf("created enrollment: student_id=%d course_id=%d", studentID, course.ID)
		} else {
			log.Printf("existing enrollment: student_id=%d course_id=%d", studentID, course.ID)
		}
	}

	log.Printf("seed complete: students=10, enrollments_created=%d", createdEnrollments)
	log.Printf("demo login password for all seeded users: SeedPass123!")
}

func ensureStudentRole() (uint, error) {
	var role model.Role
	if err := db.DB.Where("role_name = ?", "Student").First(&role).Error; err == nil {
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
	if err := db.DB.Where("email = ?", email).First(&student).Error; err == nil {
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

func ensureEnrollment(studentID, courseID uint, enrollTime time.Time) (bool, error) {
	var existing model.StudentEnrollment
	if err := db.DB.Where("student_id = ? AND course_id = ?", studentID, courseID).First(&existing).Error; err == nil {
		return false, nil
	}

	enrollment := model.StudentEnrollment{
		StudentID:  studentID,
		CourseID:   courseID,
		Status:     "registered",
		EnrollTime: enrollTime,
	}
	if err := db.DB.Create(&enrollment).Error; err != nil {
		return false, err
	}
	return true, nil
}
