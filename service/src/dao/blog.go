package dao

import (
	"sms/service/src/dao/model"
	"time"

	"github.com/jinzhu/gorm"
)

const TB_BLOG = "tb_blog_ctx"

func TestBlog() error {
	blog := &model.BlogCtx{}
	result := database.Debug().Table(TB_BLOG).First(blog)
	if result.Error != nil && result.Error.Error() != "record not found" {
		return result.Error
	}
	if blog.Code == "" {
		blog.Code = "Test1234"
		blog.CreateTime = time.Now()
		blog.UpdateTime = time.Now()
		blog.AuthorId = 1
		blog.Status = 1
		blog.Title = "Welcome !"
		blog.Content = "#  使用介绍\n"
		result = database.Debug().Table(TB_BLOG).Create(blog)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func SaveBlog(ctx *model.BlogCtx) error {
	result := database.Table(TB_BLOG).Save(ctx)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		result = database.Table(TB_BLOG).Create(ctx)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func QueryBlog(id int64, code string) *model.BlogCtx {
	if id == 0 && code == "" {
		return nil
	}
	blog := &model.BlogCtx{}
	var result *gorm.DB
	if id != 0 {
		result = database.Table(TB_BLOG).Where("id = ?", id).First(blog)
	} else {
		result = database.Table(TB_BLOG).Where("code = ?", code).First(blog)
	}
	if result.Error != nil {
		return nil
	}
	return blog
}

func ListBlogs(page, size int, authorIds []int, tags []string) []*model.BlogCtx {
	blogs := []*model.BlogCtx{}
	return blogs
}
