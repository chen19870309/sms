package dao

import (
	"errors"
	"fmt"
	"sms/service/src/dao/model"
	"sms/service/src/utils"
	"time"

	"github.com/jinzhu/gorm"
)

const TB_USER_CARD_RES = "tb_user_card_res"

func SaveUserCardRes(ctx *model.UserCardRes) error {
	utils.Log.Infof("SaveUserCardRes[%v]", ctx)
	if ctx.Id != 0 {
		result := database.Table(TB_USER_CARD_RES).Save(ctx)
		if result.Error != nil {
			return result.Error
		}
	} else {
		result := database.Table(TB_USER_CARD_RES).Create(ctx)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func EditUserCardRes(userid, cardid, status int) error {
	//1. 查找用户信息
	user := QueryUser(int64(userid), "")
	if user == nil {
		return errors.New("无效的用户ID")
	}
	//2. 查找字库信息
	word := GetWordById(cardid)
	if word == nil {
		return errors.New("无效的字库ID")
	}
	//3. 更新用户字库
	res := &model.UserCardRes{
		Userid: user.Id,
		ResId:  word.Id,
		Word:   word.Word,
		Gp:     word.Gp,
		Scope:  word.Scope,
		Status: status,
	}
	_, err := NewUserCardRes(res)
	return err
}

func NewUserCardRes(ctx *model.UserCardRes) (*model.UserCardRes, error) {
	card := &model.UserCardRes{}
	var result *gorm.DB
	result = database.Debug().Table(TB_USER_CARD_RES).Where("res_id = ? and gp = ? and word = ?", ctx.ResId, ctx.Gp, ctx.Word).First(card)
	if result.Error != nil && result.Error.Error() != "record not found" {
		return nil, result.Error
	}
	if card.Id == 0 {
		card.CreateTime = time.Now()
		card.Userid = ctx.Userid
		card.ResId = ctx.ResId //资源类型
		card.Word = ctx.Word   //字卡文字
		card.Gp = ctx.Gp       //资源分组
		card.Scope = ctx.Scope //资源范围
	}
	card.UpdateTime = time.Now()
	card.Remark = ctx.Remark //备注信息
	card.Status = ctx.Status //0生字1已学会
	key := fmt.Sprintf("USER_%v_%v_%d_%d", ctx.Scope, ctx.Gp, ctx.Status, ctx.Userid)
	utils.DelCache(key) //删除缓存
	err := SaveUserCardRes(card)
	return card, err
}

// 获取用户的生字本或者熟字本
func GetUserCardsByScope(scope, group string, userid, status, pageSize int) ([]*model.UserCardRes, error) {
	if scope == "" || group == "" {
		return nil, errors.New("scope should be inited!")
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	cards := []*model.UserCardRes{}
	key := fmt.Sprintf("USER_%v_%v_%d_%d", scope, group, status, userid)
	data, ok := utils.GetCache(key)
	if ok {
		cards = data.([]*model.UserCardRes)
		return cards, nil
	}
	var result *gorm.DB
	result = database.Debug().Table(TB_USER_CARD_RES).Limit(pageSize).Where("scope = ? and gp = ? and status = ? and userid = ?", scope, group, status, userid).Order("create_time desc").Find(&cards)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(cards) == pageSize {
		// 缓存新字列表30分钟，较少数据库查询
		utils.SetCache(key, cards, 30*time.Minute)
	}
	return cards, nil
}
