package dao

import (
	"errors"
	"fmt"
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
		user.Secret, _ = utils.EnPwdCode("smsblog")
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

func QueryUserEmail(email string) *model.SmsUser {
	user := &model.SmsUser{}
	result := database.Table(TB_USER).Where("email = ?", email).First(user)
	if result.Error != nil && result.Error.Error() == "record not found" {
		return nil
	}
	return user
}

func QueryUser(id int64, code string) *model.SmsUser {
	if id == 0 && code == "" {
		return nil
	}
	user := &model.SmsUser{}
	key := fmt.Sprintf("%d:%s", id, code)
	t, b := utils.GetCache(key)
	if b {
		utils.Log.Infof("GetCache(%s)=%v", key, t)
		user = t.(*model.SmsUser)
		return user
	}
	var result *gorm.DB
	if id != 0 {
		result = database.Table(TB_USER).Where("id = ?", id).First(user)
	} else {
		result = database.Table(TB_USER).Where("code = ?", code).First(user)
	}
	if result.Error != nil {
		return nil
	}
	utils.SetCache(key, user, 2*time.Minute)
	return user
}

func AuthUser(username, password, ip string) (*model.SmsUser, error) {
	if username == "" || password == "" {
		return nil, errors.New("username & password cant be nil!")
	}
	user := &model.SmsUser{}
	result := database.Table(TB_USER).Where("username = ? and secret = ? and status in (0,1,2)", username, password).First(user)
	if result.Error != nil {
		utils.Log.Errorf("AuthUser(%v,%v,%v) failed! ", username, password, ip)
		return nil, result.Error
	} else {
		user.UpdateTime = time.Now()
		user.LoginIp = ip
		database.Table(TB_USER).Save(user)
		user.Secret = utils.SHA1(utils.Gen8RCode())
	}
	return user, nil
}

func ExchangeUserPwd(op, np string, userid int) (*model.SmsUser, error) {
	if op == "" || np == "" {
		return nil, errors.New(" password cant be nil!")
	}
	user := &model.SmsUser{}
	result := database.Table(TB_USER).Where("id = ? and secret = ? and status in (0,1,2)", userid, op).First(user)
	if result.Error != nil {
		utils.Log.Errorf("ExchangeUserPwd(%v,%v,%v) failed! ", userid, op, np)
		return nil, errors.New("原密码输入错误")
	} else {
		user.UpdateTime = time.Now()
		user.Secret = np
		result = database.Table(TB_USER).Save(user)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return user, nil
}

func UpdateUserInfo(icon, nickname, remark string, userid int) (*model.SmsUser, error) {
	user := &model.SmsUser{}
	result := database.Table(TB_USER).Where("id = ? and status in (0,1,2)", userid).First(user)
	if result.Error != nil {
		utils.Log.Errorf("UpdateUserInfo(%v,%v,%v) failed! ", userid, nickname, remark)
		return nil, result.Error
	} else {
		user.UpdateTime = time.Now()
		user.Nickname = nickname
		user.Remark = remark
		if icon != "" {
			user.Icon = icon
		}
		result = database.Table(TB_USER).Save(user)
		if result.Error != nil {
			return nil, result.Error
		}
	}
	return user, nil
}

func RegistUser(username, password, email, code, ip string) (*model.SmsUser, error) {
	// 查验email和code
	err := CheckEmailCode(email, code)
	if err != nil {
		return nil, err
	}
	// 查验账号密码
	user := &model.SmsUser{}
	result := database.Table(TB_USER).Where("username = ? and secret = ? and status in (0,1,2)", username, password).First(user)
	if result.Error != nil && result.Error.Error() != "record not found" {
		return nil, result.Error
	}
	if user.Id > 0 {
		return nil, errors.New("用户已存在")
	}
	// 创建账号密码
	user = &model.SmsUser{
		Code:       utils.Gen8RCode(),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Level:      2,
		Status:     1,
		Icon:       "",
		Nickname:   "新来的:" + username,
		Username:   username,
		Secret:     password,
		Email:      email,
		LoginIp:    ip,
	}
	result = database.Debug().Table(TB_USER).Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
