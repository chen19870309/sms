package dao

import (
	"errors"
	"sms/service/src/dao/model"
	"sms/service/src/utils"
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

func NewEmailCode(id int64, email, code string) error {
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
	result := database.Table(TB_RESOURCE).Where("res_type = 'email' uri = ? and res_val = ?", email, code).First(res)
	if result != nil {
		if result.Error.Error() != "record not found" {
			return "", result.Error
		}
	}
	return code, NewEmailCode(res.Id, email, code)
}

func CheckEmailCode(email, code string) error {
	res := &model.Resource{}
	result := database.Table(TB_RESOURCE).Where("res_type = 'email' uri = ? and res_val = ?", email, code).First(res)
	if result != nil {
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
