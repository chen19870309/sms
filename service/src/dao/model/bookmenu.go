package model

import (
	"time"
)

type BookMenu struct {
	Id         int64
	Pid        int64
	CreateTime time.Time
	UpdateTime time.Time
	Name       string
	Sum        string
	Code       string
	Pic        string
	Remark     string
	Day        string
	Status     uint
}
