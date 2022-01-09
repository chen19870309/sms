package model

import "time"

type BlogCtx struct {
	Id         int64
	Code       string
	CreateTime time.Time
	UpdateTime time.Time
	AuthorId   uint
	Title      string
	Tags       string
	Sum        string
	Content    string
	Status     uint
}
