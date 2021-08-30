package dbcache

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	SetTag("junmo", "csq", nil)
	dsn := "root:123456@tcp(127.0.0.1:3306)/lxq?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var user1 = UserTest{
		ID:       1,
		Name:     "junmo",
		Email:    nil,
		Age:      100,
		Birthday: nil,
	}
	stmt := db.Session(&gorm.Session{DryRun: true}).Create(&user1).Statement
	stmt.SQL.String() //=> SELECT * FROM `users` WHERE `id` = $1 ORDER BY `id`
	// 注意：SQL 并不总是能安全地执行，GORM 仅将其用于日志，它可能导致会 SQL 注入
	sql := db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
	// SELECT * FROM `users` WHERE `id` = 1
	r := db.Create(&user1)
	t.Log(user1, r.Error)

	var user UserTest
	// 新建会话模式
	stmt = db.Session(&gorm.Session{DryRun: true}).First(&user, 2).Statement
	sql = db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
	result, need := Tag("junmo", sql)
	if need {
		db.First(&user, 2)
		t.Log("result11", user)
		SetTag("junmo", sql, user)
	}
	t.Log("result", result)

	var users []UserTest
	// 会话模式
	tx := db.Session(&gorm.Session{PrepareStmt: true})
	tx.First(&user, 1)
	tx.Find(&users)
	tx.Model(&user).Update("Age", 18)

	// returns prepared statements manager
	stmtManger, ok := tx.ConnPool.(*gorm.PreparedStmtDB)

	// 关闭 *当前会话* 的预编译模式
	//stmtManger.Close()

	// 为 *当前会话* 预编译 SQL
	t.Log(ok, stmtManger.PreparedSQL, stmtManger.Stmts) // => []string{}

	// 为当前数据库连接池的（所有会话）开启预编译模式
	//stmtManger.Stmts // map[string]*sql.Stmt

	for sql, stmt := range stmtManger.Stmts {
		//sql  // 预编译 SQL
		//stmt // 预编译模式
		t.Log(sql, stmt)
		stmt.Close() // 关闭预编译模式
	}

}

func BenchmarkPrepare1(b *testing.B) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/lxq?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	for i := 0; i < b.N; i++ {
		var users []UserTest
		var user UserTest
		tx := db.Session(&gorm.Session{PrepareStmt: true})
		tx.First(&user, 1)
		tx.Where("age>?", 10).Find(&users)
	}

}
func BenchmarkPrepare2(b *testing.B) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/lxq?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	for i := 0; i < b.N; i++ {
		var users []UserTest
		var user UserTest
		db.First(&user, 1)
		db.Where("age>?", 10).Find(&users)
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
