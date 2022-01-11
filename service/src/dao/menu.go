package dao

import (
	"sms/service/src/dao/model"
	"time"

	"github.com/jinzhu/gorm"
)

const TB_MENU = "tb_book_menu"

func TestMenu() error {
	menu := &model.BookMenu{}
	var result *gorm.DB
	result = database.Table(TB_MENU).First(menu)
	if result.Error != nil && result.Error.Error() == "record not found" {
		menu.Pid = 0
		menu.Name = "main menu"
		menu.Status = 1
		menu.CreateTime = time.Now()
		menu.UpdateTime = time.Now()
		menu.Remark = "blogs的主菜单,默认按月创建chapter,blog挂在chapter下"
		return SaveMenu(menu)
	}
	return result.Error
}

func CreateMonthMenu() (int64, error) {
	menu := &model.BookMenu{}
	chapter := time.Now().Format("2006,01")
	var result *gorm.DB
	result = database.Table(TB_MENU).Where("name = ?", chapter).First(menu)
	if result.Error != nil && result.Error.Error() == "record not found" {
		menu.Pid = 1
		menu.Name = chapter
		menu.Status = 1
		menu.CreateTime = time.Now()
		menu.UpdateTime = time.Now()
		menu.Remark = "按月创建chapter,blog挂在chapter下"
		err := SaveMenu(menu)
		if err != nil {
			return 0, err
		}
	}
	return menu.Id, result.Error
}

func CreateBookMenu(pid int64, blog *model.BlogCtx) error {
	menu := QueryMenu(pid, blog.Code)
	menu.Pid = pid
	menu.Name = blog.Title
	menu.Sum = blog.Sum
	menu.Code = blog.Code
	menu.Remark = blog.Tags
	menu.Status = 1
	menu.CreateTime = blog.CreateTime
	menu.UpdateTime = time.Now()
	return SaveMenu(menu)
}

func SaveMenu(ctx *model.BookMenu) error {
	if ctx.Id != 0 {
		result := database.Table(TB_MENU).Save(ctx)
		if result.Error != nil {
			return result.Error
		}
	} else {
		result := database.Table(TB_MENU).Create(ctx)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func QueryMenu(pid int64, code string) *model.BookMenu {
	menu := &model.BookMenu{}
	var result *gorm.DB
	result = database.Table(TB_MENU).Where("pid = ? and code = ?", pid, code).First(menu)
	if result.Error != nil {
		return nil
	}
	return menu
}

func QueryMenus(pid int64) []*model.BookMenu {
	menus := []*model.BookMenu{}
	var result *gorm.DB
	result = database.Table(TB_MENU).Where("pid = ?", pid).Find(&menus)
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
