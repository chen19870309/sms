package dao

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"sms/service/src/dao/model"
	"sms/service/src/utils"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

const TB_USER_DIARY = "tb_user_diary"

var Colors = []string{"green", "red", "blue", "orange"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func SaveUserDiary(ctx *model.UserDiary) error {
	result := database.Table(TB_USER_DIARY).Save(ctx)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		result = database.Table(TB_USER_DIARY).Create(ctx)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func NewUserDiary(id, userid, year, month, day int, remark string) error {
	t := time.Now()
	user := QueryUser(int64(userid), "")
	if user == nil || user.Id <= 0 {
		return errors.New("无效的用户")
	}
	//d0, err := GetUserDayDiary(id, userid, year, month, day)
	diary := model.UserDiary{
		Userid:     int64(userid),
		Year:       int64(year),
		Month:      int64(month),
		Day:        int64(day),
		Color:      Colors[rand.Intn(4)],
		Status:     1,
		CreateTime: t,
		UpdateTime: t,
	}
	// if err == nil && d0.Id > 0 {
	// 	diary.Id = d0.Id
	// 	diary.CreateTime = d0.CreateTime
	// 	remark = remark + "\n" + d0.Diary
	// }
	rs, err := utils.AesEcrypt([]byte(remark), []byte(utils.MD5(user.Code)))
	if err != nil {
		return err
	}
	diary.Diary = base64.StdEncoding.EncodeToString(rs)
	return SaveUserDiary(&diary)
}

func DecDiary(userid int, diary string) string {
	user := QueryUser(int64(userid), "")
	if user == nil || user.Id <= 0 {
		utils.Log.Errorf("QueryUser(%d) failed! empty user", userid)
		return ""
	}
	b, err := base64.StdEncoding.DecodeString(diary)
	if err == nil {
		rs, err := utils.AesDeCrypt(b, []byte(utils.MD5(user.Code)))
		if err == nil {
			return string(rs)
		} else {
			utils.Log.Errorf("AesDeCrypt(%s) [%d] failed! %v", diary, userid, err)
		}
	} else {
		utils.Log.Errorf("DecodeStringB64(%s) failed! %v", diary, err)
	}
	return ""
}

func GetUserDayDiary(id, userid, year, month, day int) (*model.UserDiary, error) {
	diary := &model.UserDiary{}
	var result *gorm.DB
	if id > 0 {
		result = database.Table(TB_USER_DIARY).Where("id = ? and userid = ? and status in (0,1)", id, userid).First(diary)
	} else {
		result = database.Table(TB_USER_DIARY).Where("userid = ? and year = ? and month = ? and day = ? and status in (0,1)", userid, year, month, day).First(diary)
	}
	if diary.Diary != "" {
		diary.Diary = DecDiary(userid, diary.Diary)
	}
	return diary, result.Error
}

func GetUserMonthDiary(userid, year, month string) []*model.UserDiary {
	res := []*model.UserDiary{}
	result := database.Table(TB_USER_DIARY).Where("userid = ? and year = ? and month = ? and status in (0,1)", userid, year, month).Find(&res)
	if result.Error != nil {
		utils.Log.Errorf("GetUserMonthDiary(%v,%v,%v) failed! %v", userid, year, month, result.Error)
	}
	id, _ := strconv.Atoi(userid)
	for _, item := range res {
		item.Diary = DecDiary(id, item.Diary)
	}
	return res
}
