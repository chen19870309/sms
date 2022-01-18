package model

import (
	"time"
)

type SmsUser struct {
	Id         int64
	Code       string
	CreateTime time.Time
	UpdateTime time.Time
	Username   string
	Nickname   string
	Secret     string
	Icon       string
	Remark     string
	Level      uint
	Status     uint
	LoginIp    string
}
