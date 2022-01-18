package dao

import (
	"errors"
	"sms/service/src/dao/model"
	"sms/service/src/utils"
	"time"

	"github.com/jinzhu/gorm"
)

const TB_USER = "tb_sms_user"

func TestUser() error {
	user := &model.SmsUser{}
	result := database.Debug().Table(TB_USER).First(user)
	if result.Error != nil && result.Error.Error() != "record not found" {
		return result.Error
	}
	if user.Code == "" {
		user.Code = utils.Gen8RCode()
		user.CreateTime = time.Now()
		user.UpdateTime = time.Now()
		user.Level = 1
		user.Status = 1
		user.Nickname = "管理员"
		user.Username = "admin"
		user.Secret, _ = utils.EnPwdCode([]byte("smsblog"))
		result = database.Debug().Table(TB_USER).Create(user)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func SaveUser(ctx *model.SmsUser) error {
	result := database.Table(TB_USER).Save(ctx)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		result = database.Table(TB_USER).Create(ctx)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func QueryUser(id int64, code string) *model.SmsUser {
	if id == 0 && code == "" {
		return nil
	}
	user := &model.SmsUser{}
	var result *gorm.DB
	if id != 0 {
		result = database.Table(TB_USER).Where("id = ?", id).First(user)
	} else {
		result = database.Table(TB_USER).Where("code = ?", code).First(user)
	}
	if result.Error != nil {
		return nil
	}
	return user
}

func AuthUser(username, password, ip string) (*model.SmsUser, error) {
	if username == "" || password == "" {
		return nil, errors.New("username & password cant be nil!")
	}
	user := &model.SmsUser{}
	result := database.Table(TB_USER).Where("username = ? and secret = ? and status in (0,1,2)", username, password).First(user)
	if result.Error != nil {
		return nil, result.Error
	} else {
		user.UpdateTime = time.Now()
		user.LoginIp = ip
		database.Table(TB_USER).Save(user)
		user.Secret = utils.SHA1(utils.Gen8RCode())
	}
	return user, nil
}
