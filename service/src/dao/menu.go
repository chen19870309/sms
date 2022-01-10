package dao

import (
	"sms/service/src/dao/model"

	"github.com/jinzhu/gorm"
)

const TB_MENU = "tb_book_menu"

func SaveMenu(ctx *model.BookMenu) error {
	result := database.Table(TB_MENU).Save(ctx)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		result = database.Table(TB_MENU).Create(ctx)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func QueryMenu(code string) *model.BookMenu {
	menu := &model.BookMenu{}
	var result *gorm.DB
	result = database.Table(TB_MENU).Where("code = ?", code).First(menu)
	if result.Error != nil {
		return nil
	}
	return menu
}

func QueryMenus(pid int64) []*model.BookMenu {
	menus := []*model.BookMenu{}
	var result *gorm.DB
	result = database.Table(TB_MENU).Where("pid = ?", pid).Find(menus)
	if result.Error != nil {
		return nil
	}
	return menus
}

func QueryBook(name string) *model.BookMenu {
	menu := &model.BookMenu{}
	var result *gorm.DB
	result = database.Table(TB_MENU).Where("name = ?", name).First(menu)
	if result.Error != nil {
		return nil
	}
	return menu
}
