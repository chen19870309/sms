package dao

import (
	"errors"
	"fmt"
	"regexp"
	"sms/service/src/dao/model"
	"sms/service/src/utils"
	"strconv"
	"time"
)

const TB_RESOURCE = "tb_resource"

func SaveResource(ctx *model.Resource) error {
	result := database.Table(TB_RESOURCE).Save(ctx)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		result = database.Table(TB_RESOURCE).Create(ctx)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

//记录用户授权码
func NewAuthCode(userid int64, code string) error {
	d, _ := time.ParseDuration("24h")
	uri := fmt.Sprintf("userid:%d", userid)
	res := &model.Resource{}
	result := database.Table(TB_RESOURCE).Where("res_type = 'account' and uri = ?", uri).First(res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			res = &model.Resource{
				ResType:    "account",
				Uri:        uri,
				ResVal:     code,
				CreateTime: time.Now(),
				ExpireTime: time.Now().Add(d),
			}
		} else {
			return result.Error
		}
	} else {
		res.CreateTime = time.Now()
		res.ExpireTime = time.Now().Add(d)
		res.ResVal = code
	}
	return SaveResource(res)
}

func CheckAuthCode(code string) *model.SmsUser {
	res := &model.Resource{}
	result := database.Table(TB_RESOURCE).Where("res_type = 'account' and res_val = ?", code).First(res)
	if result.Error != nil {
		return nil
	}
	if res.ExpireTime.Before(time.Now()) {
		return nil
	}
	userid, err := strconv.ParseInt(res.Uri[7:], 10, 32)
	if err != nil {
		return nil
	} else {
		return QueryUser(userid, "")
	}
}

func NewEmailCode(id int64, email, code string) error {
	ok, _ := regexp.Match("^\\w+@\\w+\\.\\w+$", []byte(email))
	if !ok {
		return errors.New("无效的Email地址:" + email)
	}
	d, _ := time.ParseDuration("10m")
	res := &model.Resource{
		ResType:    "email",
		Uri:        email,
		ResVal:     code,
		CreateTime: time.Now(),
		ExpireTime: time.Now().Add(d),
	}
	if id > 0 {
		res.Id = id
	}
	return SaveResource(res)
}

func GenEmailCode(email string) (string, error) {
	res := &model.Resource{}
	code := utils.GetRandNum(6)
	result := database.Table(TB_RESOURCE).Where("res_type = 'email' and uri = ?", email).First(res)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			return "", result.Error
		}
	}
	return code, NewEmailCode(res.Id, email, code)
}

func CheckEmailCode(email, code string) error {
	res := &model.Resource{}
	result := database.Table(TB_RESOURCE).Where("res_type = 'email' and uri = ? and res_val = ?", email, code).First(res)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return errors.New("无效的校验码")
		} else {
			return result.Error
		}
	}
	if res.ExpireTime.Before(time.Now()) {
		return errors.New("校验码已过期")
	}
	return nil
}
