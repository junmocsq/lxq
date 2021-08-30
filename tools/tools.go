package tools

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"time"
)

// RandStr 获取随机字符串，默认为4
func RandStr(length ...int) string {
	l := 4
	if len(length) > 0 {
		l = length[0]
	}
	str := "1234567890abcdefghijklmnopqrstuvw"
	totalLength := len(str)
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, l)
	for i := 0; i < l; i++ {
		index := rand.Intn(totalLength)
		result[i] = str[index]
	}
	return string(result)
}

func Password(passwd string, salt ...string) string {
	ss := ""
	if len(salt) > 0 {
		ss = salt[0]
	}
	s := sha1.New()
	io.WriteString(s, passwd)
	io.WriteString(s, ss)
	return fmt.Sprintf("%x", s.Sum(nil))
}
