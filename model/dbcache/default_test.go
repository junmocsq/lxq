package dbcache

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"testing"
	"time"
)

type UserTest struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func TestCache(t *testing.T) {
	db := NewDb()
	db.DB().AutoMigrate(&UserTest{})
}

func ExampleDao_Exec() {
	db := NewDb()
	var user1 = UserTest{
		ID:       3,
		Name:     "junmocsq",
		Email:    nil,
		Age:      100,
		Birthday: nil,
	}
	db.DB().AutoMigrate(&UserTest{})

	tag := "users000"
	stmt := db.DryRun().Delete(&user1).Statement
	db.SetTag(tag).PrepareSql(stmt.SQL.String(), stmt.Vars...).EXEC()
	stmt = db.DryRun().Create(&user1).Statement
	n, err := db.SetTag(tag).PrepareSql(stmt.SQL.String(), stmt.Vars...).EXEC()
	fmt.Println(n, err)
	// OutPut: 1 <nil>
}

func ExampleDao_Fetch() {
	db := NewDb()
	var user UserTest
	db.DB().AutoMigrate(&UserTest{})
	stmt := db.DryRun().Find(&user, 1).Statement
	tag := "users000"
	err := db.SetTag(tag).PrepareSql(stmt.SQL.String(), stmt.Vars...).Fetch(&user)
	fmt.Println(user.ID, err)

	var users []UserTest
	// 会话模式
	stmt = db.DryRun().Find(&users).Statement
	db.SetTag(tag).PrepareSql(stmt.SQL.String(), stmt.Vars...).Fetch(&users)
	fmt.Println("result", len(users))
	// OutPut:
	// 	1 <nil>
	//result 3
}
func BenchmarkPrepare1(b *testing.B) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/lxq?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		var users []UserTest
		var user UserTest
		tx := db.Session(&gorm.Session{PrepareStmt: true})
		tx.First(&user, 1)
		tx.Where(fmt.Sprintf("name='junmo-%d' AND age>?", rand.Int()), 10).Find(&users)
	}

}

func BenchmarkSelect(b *testing.B) {
	db := NewDb()
	var user UserTest
	stmt := db.DryRun().Find(&user, 1).Statement
	tag := "users000"
	for i := 0; i < b.N; i++ {
		db.SetTag(tag).PrepareSql(stmt.SQL.String(), stmt.Vars...).Fetch(&user)
	}
}
func BenchmarkPrepare2(b *testing.B) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/lxq?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		var users []UserTest
		var user UserTest
		db.First(&user, 1)
		db.Where(fmt.Sprintf("name='junmo-%d' AND age>?", rand.Int()), 10).Find(&users)
	}
}
func BenchmarkPrepare3(b *testing.B) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/lxq?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		var users []UserTest
		var user UserTest
		tx := db.Session(&gorm.Session{PrepareStmt: true})
		tx.First(&user, 1)
		tx.Where("name=? AND age>?", fmt.Sprintf("junmo-%d", rand.Int()), 10).Find(&users)
	}
}

func BenchmarkHashMd5(b *testing.B) {
	str := "junmocsqjunmocsqjunmocsqjunmocsqjunmocsqjunmocsqjunmocsq"
	b.Run("hash-1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			md5Hash(str)
		}
	})
	b.Run("hash-2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			md5Hash2(str)
		}
	})
	b.Run("hash-3", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			md5Hash3(str)
		}
	})
}
