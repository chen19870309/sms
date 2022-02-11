package dao

import (
	"fmt"
	"sms/service/src/dao/model"
	"sms/service/src/utils"
	"strconv"
	"strings"
	"time"

	"github.com/Chain-Zhang/pinyin"
	"github.com/jinzhu/gorm"
)

const TB_CARD_RES = "tb_card_res"

func GetWordById(id int) *model.CardRes {
	card := &model.CardRes{}
	database.Debug().Table(TB_CARD_RES).Where("res_type = 'words' and id = ?", id).First(card)
	if card.Id > 0 {
		return card
	}
	return nil
}

func SaveCardRes(ctx *model.CardRes) error {
	utils.Log.Infof("SaveCardRes[%v]", ctx)
	if ctx.Id != 0 {
		result := database.Table(TB_CARD_RES).Save(ctx)
		if result.Error != nil {
			return result.Error
		}
	} else {
		result := database.Table(TB_CARD_RES).Create(ctx)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func NewCardRes(ctx *model.CardRes) (*model.CardRes, error) {
	d, _ := time.ParseDuration("24h")
	card := &model.CardRes{}
	var result *gorm.DB
	result = database.Debug().Table(TB_CARD_RES).Where("res_type = ? and gp = ? and word = ?", ctx.ResType, ctx.Gp, ctx.Word).First(card)
	if result.Error != nil && result.Error.Error() != "record not found" {
		return nil, result.Error
	}
	card.CreateTime = time.Now()
	card.ExpireTime = time.Now().Add(d)
	card.Gp = ctx.Gp           //资源分组
	card.Scope = ctx.Scope     //资源范围
	card.ResType = ctx.ResType //资源类型
	card.Word = ctx.Word       //字卡文字
	card.Data = ctx.Data       //卡片内容
	card.Pinyin = ctx.Pinyin   //字卡拼音
	card.Pic = ctx.Pic         //卡片图片
	card.Sound = ctx.Sound     //卡片声音
	err := SaveCardRes(card)
	return card, err
}

func SaveCardInfo(data string) {
	ls := strings.Split(data, "\n")
	beginCard := false
	card := &model.CardRes{
		Gp:      "PB",
		ResType: "words",
		Scope:   "",
	}
	for _, item := range ls {
		str := strings.TrimSpace(item)
		if str == "" {
			continue
		}
		if strings.HasPrefix(str, "@scope=") {
			card.Scope = str[7:]
		}
		if strings.HasPrefix(str, "@group=") {
			card.Gp = str[7:]
		}
		if strings.HasPrefix(str, "## ") {
			beginCard = true
			items := strings.Split(str, " ")
			if len(items) == 2 {
				card.Word = items[1]
				py, err := pinyin.New(card.Word).Split(" ").Mode(pinyin.Tone).Convert()
				if err != nil {
					card.Pinyin = ""
					utils.Log.Errorf("pinyin Convert(%v) => error [%v]", card.Word, err)
				} else {
					card.Pinyin = py
				}
			} else if len(items) == 3 {
				card.Word = items[1]
				card.Pinyin = items[2]
			} else {
				beginCard = false
			}
			continue
		}
		if beginCard {
			if strings.HasPrefix(str, "![image](") {
				card.Pic = str[9 : len(str)-1]
			} else if strings.HasPrefix(str, "[sound](") {
				card.Sound = str[8 : len(str)-1]
			} else if strings.HasPrefix(str, "> ") {
				card.Data = str[2:]
			} else if strings.HasPrefix(str, "***") {
				beginCard = false
				NewCardRes(card)
			}
		}
	}
}

func GetResScopes(userid, resType string) map[string]interface{} {
	res := make(map[string]interface{})
	rs := make([]model.UserScopeCount, 0)
	sp1 := model.UserScopeCount{
		Scope: "生字本",
		Color: "orange",
		Icon:  "favor",
		Gp:    "private",
	}
	sp2 := model.UserScopeCount{
		Scope: "已学会",
		Color: "green",
		Icon:  "favorfill",
		Gp:    "private",
	}
	result := database.Table(TB_CARD_RES).Select("scope, gp, count(id) as cnt").Group("scope,gp").Where("res_type = 'words'").Find(&rs)
	if result.Error != nil {
		utils.Log.Error("GetResScopes:", result.Error)
	} else {
		for _, item := range rs {
			sp1.Cnt += CountUserScopeWords(userid, item.Scope, item.Gp, 0) //获取用户导入生字
			c1 := CountUserScopeWords(userid, item.Scope, item.Gp, 1)      //获取用户已学会字数
			item.Ucnt = c1
			sp2.Cnt += c1
			item.Color = "cyan"
			item.Icon = "favor"
			if c1 == item.Cnt {
				item.Icon = "favorfill"
			}
			res[item.Scope] = item
		}
	}
	//添加生字本
	res[sp1.Scope] = sp1
	//添加已学会
	res[sp2.Scope] = sp2
	return res
}

func CountUserScopeWords(userid, scope, gp string, status int) int {
	ucnt := 0
	result := database.Table(TB_USER_CARD_RES).Where("userid = ? and status = ? and scope = ? and gp = ?", userid, status, scope, gp).Count(&ucnt)
	if result.Error != nil {
		utils.Log.Error("CountUserScopeWords:", result.Error)
	} else {
		return ucnt
	}
	return 0
}

func QueryScopedCard(userid, scope, ResType, gp string, pageSize int) ([]*model.CardRes, error) {
	if userid != "" {
		uid, err := strconv.Atoi(userid)
		if err == nil {
			res, err := QueryUserStdWordCards(scope, gp, pageSize, uid)
			if err == nil {
				return res, err
			}
		}
		utils.Log.Error("QueryUserStdWordCards:", err)
	}
	cards := []*model.CardRes{}
	key := fmt.Sprintf("WORD_%v_%v_%v_%d", scope, ResType, gp, pageSize)
	data, ok := utils.GetCache(key)
	if ok {
		cards = data.([]*model.CardRes)
		return cards, nil
	}
	var result *gorm.DB
	if gp != "" {
		result = database.Table(TB_CARD_RES).Limit(pageSize).Where("scope = ? and res_type = ? and gp = ?", scope, ResType, gp).Order("id desc").Find(&cards)
	} else { //查询公开的blog
		result = database.Table(TB_CARD_RES).Limit(pageSize).Where("scope = ? and res_type = ? and gp = 'PB'", scope, ResType).Order("id desc").Find(&cards)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	// 缓存新字列表30分钟，较少数据库查询
	utils.SetCache(key, cards, 30*time.Minute)
	return cards, nil
}

func QueryUserStdWordCards(scope, group string, pageSize, userid int) ([]*model.CardRes, error) {
	cards := []*model.CardRes{}
	var result *gorm.DB
	sub := database.Table(TB_USER_CARD_RES).Select("res_id").Where("userid = ? and status = 1", userid).SubQuery() //已学会
	if scope == "生字本" {
		sub = database.Table(TB_USER_CARD_RES).Select("res_id").Where("userid = ? and status = 0", userid).SubQuery() //生字本
	}
	result = database.Table(TB_CARD_RES).Limit(pageSize).Where("res_type = 'words' and id in ?", sub).Order("id desc").Find(&cards)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(cards) > 0 && (scope == "生字本" || scope == "已学会") {
		return cards, nil
	}
	ls, err := GetUserCardsByScope(scope, group, userid, 0, pageSize) //获取没学会的
	if err != nil {
		return nil, err
	}
	if len(ls) < pageSize { //没有生字或者生字不满足本次请求page数量
		//已学会
		sub := database.Debug().Table(TB_USER_CARD_RES).Select("res_id").Where("userid = ? and status = 1", userid).SubQuery()
		//查询未学的字
		result = database.Debug().Table(TB_CARD_RES).Limit(pageSize).Where("scope = ? and gp = ? and res_type = 'words' and id NOT in ?", scope, group, sub).Order("id desc").Find(&cards)
		if result.Error != nil {
			return nil, result.Error
		}
		if len(cards) > 0 {
			for _, word := range cards {
				NewUserCardRes(&model.UserCardRes{
					Userid: int64(userid),
					ResId:  word.Id,
					Scope:  word.Scope,
					Gp:     word.Gp,
					Word:   word.Word,
					Status: 0, //添加到生字本
				})
			}
		}
	}
	//重新拉取生字
	sub = database.Debug().Table(TB_USER_CARD_RES).Select("res_id").Where("userid = ? and status = 0", userid).SubQuery()
	result = database.Debug().Table(TB_CARD_RES).Limit(pageSize).Where("scope = ? and gp = ? and res_type = 'words' and id in ? ", scope, group, sub).Order("id desc").Find(&cards)
	if result.Error != nil {
		return nil, result.Error
	}
	return cards, nil
}
