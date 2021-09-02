package university

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestPart3(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/lxq?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键
	})
	db = db.Debug()
	if err != nil {
		panic("failed to connect database")
	}
	//instructor := Instructor{
	//	Id:         "10211",
	//	Name:       "Smith",
	//	DeptName:   "Biologys",
	//	Salary:     66000,
	//}
	//db.Delete(&instructor)
	//db.Create(&instructor)
	var res []Instructor
	db.Model(&Instructor{}).Joins("Department").Scan(&res)
	t.Log(res)

	var res1 []Instructor
	db.Model(&Instructor{}).Joins("Department").Select("name", "instructor.dept_name", "building").Where("14365").Or("15347").Find(&res1)
	t.Log(res1)

	// with max_budget (value) as (select max(budget) from department) select *,budget from department,max_budget where department.budget=max_budget.value;

	// WITH dept_total(dept_name,value) AS (SELECT dept_name,SUM(salary) from instructor GROUP BY dept_name),
	// dept_total_avg(value) as (SELECT AVG(value) FROM dept_total)
	// SELECT dept_name,dept_total.value FROM dept_total,dept_total_avg WHERE dept_total.value>=dept_total_avg.value
}

type User struct {
	gorm.Model
	Name         string `gorm:"<-:create;->:false"`
	Email        *string
	Age          uint8
	Birthday     uint32
	MemintNumber uint32
	ActivatedAt  uint32
	Admin        bool
}

func UserTable(user User) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if user.Admin {
			return tx.Table("admin_users")
		}

		return tx.Table("users")
	}
}

func TestUser(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/lxq?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键
	})
	db = db.Debug()
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	user := User{Name: "Jinzhu", Age: 18, Birthday: uint32(time.Now().Unix()), Admin: true}
	result := db.Create(&user)
	t.Log(user.ID, result.Error, result.RowsAffected)
	var user1 User = User{
		Admin: true,
	}
	db.Scopes(UserTable(user1)).Create(&user1)
	//var res []User
	db.Scopes(UserTable(user1)).Where("name=?", "Jinzhu").Scan(&user1)
	t.Logf("%#v", user1)
}
