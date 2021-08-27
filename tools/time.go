package tools

import (
	"fmt"
	"time"
)

func Time2Read(unixtime int64) string {
	if unixtime == 0 {
		unixtime = time.Now().Unix()
	}
	return time.Unix(unixtime, 0).Format("2006-01-02 15:04:05")
}

func Time2Read2(unixtime int64) string {
	now := time.Now().Unix()
	if unixtime == 0 {
		unixtime = now
	}
	m := now - unixtime
	t := time.Unix(unixtime, 0)
	if m < 300 {
		return "刚刚"
	} else if m < 3600 {
		return fmt.Sprintf("%d 分钟前", m/60)
	} else if m < 86400 {
		return fmt.Sprintf("%d 小时前", m/3600)
	} else if t.Year() == time.Now().Year() {
		return t.Format("2006-01-02 15:04:05")
	} else if t.Year() < time.Now().Year() {
		return t.Format("2006-01-02")
	} else {
		return t.Format("2006-01-02 15:04:05")
	}
}

func Read2UnixTime(str string) (time.Time, error) {
	//时间 to 时间戳
	loc, _ := time.LoadLocation("Asia/Shanghai")                 //设置时区
	return time.ParseInLocation("2006-01-02 15:04:05", str, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
}
