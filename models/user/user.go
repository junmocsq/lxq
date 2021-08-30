package user

type User struct {
	Id         int64
	Nickname   string
	Avatar     string
	Passwd     string
	Email      string
	Phone      int64
	NationCode int64
	Salt       string
	Gender     uint8
	WhatIsUp   string
	LoginTime  int64
	CreatedAt  int64 `gorm:"autoUpdateTime"`
	UpdatedAt  int64 `gorm:"autoCreateTime"`
}
