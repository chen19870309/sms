package dao

import (
	"sms/service/src/dao/model"
	"sms/service/src/utils"
	"time"

	"github.com/jinzhu/gorm"
)

const TB_BLOG = "tb_blog_ctx"

func TestBlog() error {
	blog := &model.BlogCtx{}
	result := database.Table(TB_BLOG).First(blog)
	if result.Error != nil && result.Error.Error() != "record not found" {
		return result.Error
	}
	if blog.Code == "" {
		blog.Code = utils.Gen8RCode()
		blog.CreateTime = time.Now()
		blog.UpdateTime = time.Now()
		blog.AuthorId = 1
		blog.Status = 1
		blog.Title = "Welcome !"
		blog.Content = "#  使用介绍\n"
		result = database.Table(TB_BLOG).Create(blog)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func AutoSaveBlog(code, theme, data string, authorid uint) error {
	blog := &model.BlogCtx{
		Content:    data,
		UpdateTime: time.Now(),
		Title:      utils.GetMdTitle(data),
		Tags:       utils.GetMdTags(data, theme),
	}
	result := database.Debug().Table(TB_BLOG).Where("code = ? and author_id = ?", code, authorid).Update(blog)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func SaveBlog(ctx *model.BlogCtx) error {
	if ctx.Id != 0 {
		result := database.Table(TB_BLOG).Save(ctx)
		if result.Error != nil {
			return result.Error
		}
	} else {
		result := database.Table(TB_BLOG).Create(ctx)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func NewBlog(authorId uint) (*model.BlogCtx, error) {
	blog := &model.BlogCtx{
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		AuthorId:   authorId,
		Status:     0,
		Content:    "#  New Edit!\n@create time:" + time.Now().Format("2006-01-02 15:04:05") + "\n",
	}
	for {
		code := utils.Gen8RCode()
		result := database.Table(TB_BLOG).Where("code = ?", code).First(blog)
		if result.Error != nil {
			if result.Error.Error() != "record not found" {
				return nil, result.Error
			} else {
				blog.Code = code
				break
			}
		}
	}
	result := database.Table(TB_BLOG).Create(blog)
	if result.Error != nil {
		return nil, result.Error
	}
	return blog, nil
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

// 发布文章
func PutBlog(code, data string) *model.BlogCtx {
	blog := QueryBlog(0, code)
	if blog != nil {
		blog.Content = data
		blog.Title = utils.GetMdTitle(data)
		blog.Tags = utils.GetMdTags(data, "")
		blog.Status = 1
		SaveBlog(blog)
		pid, err := CreateMonthMenu()
		if err != nil {
			utils.Log.Errorf("CreateMonthMenu failed![%v]", err)
			return nil
		} else {
			err = CreateBookMenu(pid, blog)
			if err != nil {
				utils.Log.Errorf("CreateBookMenu failed![%v]", err)
				return nil
			}
		}
	}
	return blog
}

func QueryBlogCaches(auther_id int) []*model.BlogCtx {
	blogs := []*model.BlogCtx{}
	result := database.Table(TB_BLOG).Where("status = 0 and author_id = ?", auther_id).Find(&blogs)
	if result.Error != nil {
		return nil
	}
	return blogs
}

func ListBlogs(page, size int, authorIds []int, tags []string) []*model.BlogCtx {
	blogs := []*model.BlogCtx{}
	result := database.Table(TB_BLOG).Where("status = 0").Find(&blogs)
	if result.Error != nil {
		return nil
	}
	return blogs
}
