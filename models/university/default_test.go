package university

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

// db-book.com 素材地址
func TestCourse_TableName(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/lxq?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	//db = db.Debug()
	if err != nil {
		panic("failed to connect database")
	}
	// drop table classroom,teaches,section,instructor,department ,course;
	db.AutoMigrate(&Classroom{}, &Department{}, &Course{}, &Instructor{}, &Section{}, &Teaches{})
	// drop table student,takes,advisor ,time_slot,prereq;
	db.AutoMigrate(&Student{}, &Takes{}, &Advisor{}, &TimeSlot{}, &Prereq{})
	// drop table student,takes,advisor ,time_slot,prereq,classroom,teaches,section,instructor,department ,course;
	// source /Users/junmo/go/src/lxq/models/university/smallsql.sql

}
