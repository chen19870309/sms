package dao

import (
	"fmt"
	"sms/service/src/dao/model"
	"sms/service/src/utils"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

const TB_CARD_RES = "tb_card_res"

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
			if len(items) == 3 {
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

func QueryScopedCard(scope, ResType, gp string, pageSize int) ([]*model.CardRes, error) {
	cards := []*model.CardRes{}
	key := fmt.Sprintf("WORD_%v_%v_%v_%d", scope, ResType, gp, pageSize)
	data, ok := utils.GetCache(key)
	if ok {
		cards = data.([]*model.CardRes)
		return cards, nil
	}
	var result *gorm.DB
	if gp != "" {
		result = database.Debug().Table(TB_CARD_RES).Limit(pageSize).Where("scope = ? and res_type = ? and gp = ?", scope, ResType, gp).Order("id desc").Find(&cards)
	} else { //查询公开的blog
		result = database.Debug().Table(TB_CARD_RES).Limit(pageSize).Where("scope = ? and res_type = ? and gp = 'PB'", scope, ResType).Order("id desc").Find(&cards)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	// 缓存新字列表30分钟，较少数据库查询
	utils.SetCache(key, cards, 30*time.Minute)
	return cards, nil
}
