package dbcache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/junmocsq/jlib/jredis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"time"
)

var (
	redisModule = "sql"
	expire      = 300
	emptyString = "MNIL"
	dbs         map[string]*gorm.DB
	defaultDb   = "default"
	dbDebug     = false
	redis       = jredis.NewRedis(redisModule)
)

func init() {
	RegisterCacheAccessor("127.0.0.1", "6379", "")
	dsn := "root:123456@tcp(127.0.0.1:3306)/lxq?charset=utf8mb4&parseTime=True&loc=Local"
	dbs = make(map[string]*gorm.DB)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//db = db.Debug()
	if err != nil {
		panic("failed to connect database")
	}
	dbs[defaultDb] = db
}

func SETModule(module string) {
	redisModule = module
}
func SETExpire(e int) {
	expire = e
}
func Debug(debug ...bool) {
	if len(debug) > 0 {
		dbDebug = debug[0]
	} else {
		dbDebug = true
	}
}
func RegisterCacheAccessor(host, port, auth string, debug ...bool) {
	jredis.RegisterRedisPool(host, port, jredis.ModuleConf(redisModule), jredis.AuthConf(auth), jredis.PrefixConf(redisModule))
	if len(debug) > 0 {
		jredis.SetDebug(debug[0])
	}
}

type Daoer interface {
	DB() *gorm.DB
	DryRun() *gorm.DB
	SetTag(tag string) Daoer
	SetKey(key string) Daoer
	PrepareSql(sql string, params ...interface{}) Daoer
	Fetch(result interface{}) error
	EXEC() (int64, error)
	ClearCache()
}
type Dao struct {
	db     *gorm.DB
	tag    string
	key    string
	sql    string
	params []interface{}
}

func NewDb(dbname ...string) *Dao {
	dbStr := defaultDb
	if len(dbname) > 0 {
		dbStr = dbname[0]
	}
	db := dbs[dbStr]
	if db == nil {
		panic(fmt.Sprintf("%s 不存在,请配置数据库", dbStr))
	}
	if dbDebug {
		db = db.Debug()
	}
	return &Dao{db: db}
}
func (d *Dao) DB() *gorm.DB {
	return d.db
}
func (d *Dao) DryRun() *gorm.DB {
	return d.db.Session(&gorm.Session{DryRun: true})
}
func (d *Dao) SetTag(tag string) Daoer {
	d.tag = tag
	return d
}
func (d *Dao) SetKey(key string) Daoer {
	d.key = d.tag + key
	return d
}
func (d *Dao) PrepareSql(sql string, params ...interface{}) Daoer {
	d.sql = sql
	d.params = params
	return d
}
func (d *Dao) ClearCache() {
	if d.key != "" {
		redis.DEL(d.key)
	} else if d.tag != "" {
		redis.DEL(d.tag)
	}
}
func (d *Dao) clear() {
	d.tag = ""
	d.key = ""
	d.sql = ""
	d.params = nil
}
func (d *Dao) Fetch(result interface{}) error {
	defer d.clear()
	strJson, need := d.cache()
	if need {
		res := d.db.Raw(d.sql, d.params...).Scan(result)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			result = nil
		}
		d.setCache(result)
	} else {
		return json.Unmarshal([]byte(strJson), result)
	}
	return nil
}

func (d *Dao) EXEC() (int64, error) {
	defer d.clear()
	res := d.db.Exec(d.sql, d.params...)
	if res.Error != nil && res.RowsAffected > 0 {
		d.ClearCache()
	}
	return res.RowsAffected, res.Error
}

func (d *Dao) cache() (result string, needSelectDb bool) {
	if d.tag == "" {
		return "", true
	}
	key := d.key
	if key == "" {
		key = d.buildKey()
	}
	result = redis.GET(key)
	if result == "" {
		return "", true
	}
	if result == emptyString {
		result = ""
	}
	redis.EXPIRE(key, expire)
	return
}

func (d *Dao) setCache(data interface{}) bool {
	// 无缓存
	if d.tag == "" {
		return true
	}
	key := d.key
	// 没有指定key，使用tag二级缓存
	if key == "" {
		key = d.buildKey()
	}
	s := emptyString
	if data != nil {
		jsonRes, err := json.Marshal(data)
		if err != nil {
			return false
		}
		if len(jsonRes) != 0 {
			s = string(jsonRes)
		}
	}
	return redis.SETEX(key, s, expire)
}

func (d *Dao) buildKey() string {
	tagCache := redis.GET(d.tag)
	if tagCache == "" {
		tagCache = fmt.Sprintf("%s-%d", d.tag, time.Now().UnixNano())
		redis.SETEX(d.tag, tagCache, expire)
	} else {
		redis.EXPIRE(d.tag, expire)
	}
	return md5Hash(tagCache + d.db.Dialector.Explain(d.sql, d.params...))
}

// 326.9 ns/op
func md5Hash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 512.7 ns/op
func md5Hash2(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

// 467.5 ns/op
func md5Hash3(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}
