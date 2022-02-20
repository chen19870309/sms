package model

import (
	"fmt"
	"sms/service/src/utils"
	"time"

	"github.com/jinzhu/gorm"
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
	Email      string
}

// 查询用户信息时缓存用户的code
func (user *SmsUser) AfterFind(tx *gorm.DB) (err error) {
	//utils.Log.Info("Call SmsUser AfterFind")
	if user.Code != "" {
		keyname := fmt.Sprintf("user_code_%v", user.Id)
		utils.SetCache(keyname, user.Code, 2*time.Hour)
	}
	return
}
