package main

import (
	"fmt"

	fitdb "my-course-backend/db"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type tableRow struct{ Name string }
type colRow struct{ Name string }

func hasTable(db *gorm.DB, t string) bool {
	var c int64
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", t).Scan(&c)
	return c > 0
}

func hasColumn(db *gorm.DB, t, c string) bool {
	var cols []colRow
	db.Raw(fmt.Sprintf("PRAGMA table_info(%s)", t)).Scan(&cols)
	for _, col := range cols {
		if col.Name == c {
			return true
		}
	}
	return false
}

func main() {
	dbPath := fitdb.ResolvedDBPath()
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if hasTable(db, "StudentEnrollment") && !hasTable(db, "Enrollment") {
		if err := db.Migrator().RenameTable("StudentEnrollment", "Enrollment"); err != nil {
			fmt.Printf("rename table failed: %v\n", err)
		} else {
			fmt.Println("renamed table: StudentEnrollment -> Enrollment")
		}
	}

	if hasTable(db, "Enrollment") && hasColumn(db, "Enrollment", "student_id") && !hasColumn(db, "Enrollment", "user_id") {
		if err := db.Migrator().RenameColumn("Enrollment", "student_id", "user_id"); err != nil {
			fmt.Printf("rename column failed: %v\n", err)
		} else {
			fmt.Println("renamed column: Enrollment.student_id -> Enrollment.user_id")
		}
	}

	var tables []tableRow
	db.Raw("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name").Scan(&tables)
	fmt.Println("tables:")
	for _, t := range tables {
		fmt.Println("-", t.Name)
	}

	if hasTable(db, "Enrollment") {
		var cols []colRow
		db.Raw("PRAGMA table_info(Enrollment)").Scan(&cols)
		fmt.Println("Enrollment columns:")
		for _, c := range cols {
			fmt.Println("-", c.Name)
		}
	}
}
