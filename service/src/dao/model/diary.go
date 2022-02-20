package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"sms/service/src/utils"
	"time"

	"github.com/jinzhu/gorm"
)

type UserDiary struct {
	Id         int64     `json:"id" `             //id serial primary key,
	Userid     int64     `json:"userid" `         //userid integer not null ,
	Year       int64     `json:"year" `           //year integer not null ,
	Month      int64     `json:"month" `          //month integer not null ,
	Day        int64     `json:"day" `            //day integer not null ,
	CreateTime time.Time `json:"create_time" `    //create_time timestamp default CURRENT_TIMESTAMP ,
	UpdateTime time.Time `json:"update_time" `    //update_time timestamp default CURRENT_TIMESTAMP,
	Status     int64     `json:"status" `         //status int not null default 0,
	Color      string    `json:"color" `          //color varchar(32) not null default 'orange',
	Diary      string    `json:"diary" enc:"aes"` //diary text
}

//保存前加密数据
func (diary *UserDiary) BeforeSave(db *gorm.DB) (err error) {
	utils.Log.Infof("Call UserDiary BeforeSave [%s]", diary.Diary)
	if diary.Diary != "" {
		keyname := fmt.Sprintf("user_code_%v", diary.Userid)
		key, ok := utils.GetCache(keyname)
		if !ok {
			utils.Log.Errorf("GetCache(%s) failed!", keyname)
			return errors.New("获取用户密钥code失败!")
		}
		rs, err := utils.AesEcrypt([]byte(diary.Diary), []byte(utils.MD5(key.(string))))
		if err != nil {
			utils.Log.Errorf("AesEcrypt(%s) [%d] failed! %v", diary.Diary, diary.Userid, err)
			return err
		}
		diary.Diary = base64.StdEncoding.EncodeToString(rs)
	}
	return
}

//返回前解密数据
func (diary *UserDiary) AfterFind(tx *gorm.DB) (err error) {
	//utils.Log.Info("Call UserDiary AfterFind")
	if diary.Diary != "" {
		b, err := base64.StdEncoding.DecodeString(diary.Diary)
		if err == nil {
			keyname := fmt.Sprintf("user_code_%v", diary.Userid)
			key, ok := utils.GetCache(keyname)
			if !ok {
				utils.Log.Errorf("GetCache(%s) failed!", keyname)
				return errors.New("获取用户密钥code失败!")
			}
			rs, err := utils.AesDeCrypt(b, []byte(utils.MD5(key.(string))))
			if err == nil {
				diary.Diary = string(rs)
			} else {
				utils.Log.Errorf("AesDeCrypt(%s) [%d] failed! %v", diary.Diary, diary.Userid, err)
			}
		} else {
			utils.Log.Errorf("DecodeStringB64(%s) failed! %v", diary.Diary, err)
		}
	}
	return
}

func (diary *UserDiary) DoDec() {
	t := reflect.TypeOf(diary)
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("enc")
		if tag != "" {
			b, err := base64.StdEncoding.DecodeString(diary.Diary)
			if err == nil {
				keyname := fmt.Sprintf("user_code_%v", diary.Userid)
				key, ok := utils.GetCache(keyname)
				if !ok {
					utils.Log.Errorf("GetCache(%s) failed!", keyname)
					return
				}
				rs, err := utils.AesDeCrypt(b, []byte(utils.MD5(key.(string))))
				if err == nil {
					diary.Diary = string(rs)
				} else {
					utils.Log.Errorf("AesDeCrypt(%s) [%d] failed! %v", diary.Diary, diary.Userid, err)
				}
			} else {
				utils.Log.Errorf("DecodeStringB64(%s) failed! %v", diary.Diary, err)
			}
		}
	}
}
