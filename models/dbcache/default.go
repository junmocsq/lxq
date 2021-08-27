package dbcache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/junmocsq/jlib/jredis"
	"io"
	"math/rand"
	"time"
)

var redisModule = "sql"
var expire = 300
var emptyString = "MNIL"

func init() {
	RegisterCacheAccessor("127.0.0.1", "6379", "")
}

func SETModule(module string) {
	redisModule = module
}
func SETExpire(e int) {
	expire = e
}

func RegisterCacheAccessor(host, port, auth string, debug ...bool) {
	jredis.RegisterRedisPool(host, port, jredis.ModuleConf(redisModule), jredis.AuthConf(auth), jredis.PrefixConf(redisModule))
	if len(debug) > 0 {
		jredis.SetDebug(debug[0])
	}
}

func Tag(tag string, sql string) (result string, needSelectDb bool) {
	r := jredis.NewRedis(redisModule)
	tagCache := r.GET(tag)
	if tagCache == "" {
		tagCache = time.Now().String() + fmt.Sprintf("%d", rand.Int63())
		r.SETEX(tag, tagCache, expire)
	} else {
		r.EXPIRE(tag, expire)
	}
	key := md5Hash(tagCache + sql)
	result = r.GET(key)
	if result == "" {
		return "", true
	}
	if result == emptyString {
		result = ""
	}
	r.EXPIRE(key, expire)
	return
}

func SetTag(tag string, sql string, data interface{}) bool {
	r := jredis.NewRedis(redisModule)
	tagCache := r.GET(tag)
	if tagCache == "" {
		tagCache = time.Now().String() + fmt.Sprintf("%d", rand.Int63())
		r.SETEX(tag, tagCache, expire)
	}else {
		r.EXPIRE(tag, expire)
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
	return r.SETEX(md5Hash(tagCache+sql), s, expire)
}

func ClearCache(tagOrKey string) int {
	r := jredis.NewRedis(redisModule)
	return r.DEL(tagOrKey)
}

func Key(key string) (result string, needSelectDb bool) {
	r := jredis.NewRedis(redisModule)
	result = r.GET(key)
	if result == "" {
		return "", true
	}
	if result == emptyString {
		result = ""
	}
	r.EXPIRE(key, expire)
	return
}

func SetKey(key string, data interface{}) bool {
	r := jredis.NewRedis(redisModule)
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
	return r.SETEX(fmt.Sprintf("%x", key), s, expire)
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
