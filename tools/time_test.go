package tools

import (
	"testing"
	"time"
)

func TestTime2Read(t *testing.T) {
	if Time2Read(1629796823) != "2021-08-24 17:20:23" {
		t.Error("Time2Read 时间戳转换失败")
	}
	t.Log(Time2Read2(time.Now().Unix()))
	t.Log(Time2Read2(time.Now().Unix() - 361))
	t.Log(Time2Read2(time.Now().Unix() - 3599))
	t.Log(Time2Read2(time.Now().Unix() - 10800))
	t.Log(Time2Read2(time.Now().Unix() - 864000))
	t.Log(Time2Read2(time.Now().Unix() - 365*86400))
	r, err := Read2UnixTime("2021-08-24 17:20:23")
	if err != nil || r.Unix() != 1629796823 {
		t.Error("Read2UnixTime 时间转换时间戳失败", err)
	}
	t.Log(Read2UnixTime("2021-08-24 17:20"))

}
