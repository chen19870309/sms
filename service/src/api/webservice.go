package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sms/service/src/api/model"
	"sms/service/src/dao"
	"sms/service/src/utils"

	"github.com/gin-gonic/gin"
)

type WebS struct {
	Serv *gin.Engine
}

var service *WebS

func init() {
	service = &WebS{}
	service.Serv = gin.Default()
	service.Serv.Use(Cors())
	service.Serv.Use(BlogAuth())
	service.Serv.StaticFS("/static", http.Dir("/Users/chenchunjiang/go/src/sms/webapp/dist/static"))
	// service.Serv.StaticFile("/", "/Users/chenchunjiang/go/src/sms/webapp/dist/index.html")
	service.Serv.StaticFile("/", "/Users/chenchunjiang/go/src/sms/webapp/dist/index.html")
}

func (web *WebS) Close() {
	dao.CloseDB()
}

func NewBlogService() (w *WebS, err error) {
	dao.InitDB()
	err = dao.TestUser()
	if err != nil {
		utils.Log.Errorf("TestUser Failed![%v]", err)
	}
	err = dao.TestBlog()
	if err != nil {
		utils.Log.Errorf("TestBlog Failed![%v]", err)
	}
	err = dao.TestMenu()
	if err != nil {
		utils.Log.Errorf("TestMenu Failed![%v]", err)
	}
	blog := service.Serv.Group("blog")
	blog.GET("/ping", Pong)
	blog.GET("/menu", GetMenu)
	blog.POST("/login", LoginBlog)
	blog.POST("/user/edit", EditUser)
	blog.GET("/newblog", NewBlogPage)
	blog.GET("/blogcaches", BlogCaches)
	blog.POST("/save/:code", SaveBlog)
	blog.GET("/posts/:code", GetPosts)
	blog.PUT("/posts/:code", PutPosts)
	return service, nil
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func BlogAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := c.GetRawData()
		if err != nil {
			utils.Log.Error("GetRawData Failed!:", err)
		} else {
			utils.Log.Infof("req=[%v]", string(data))
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // 关键点
		}

		auth := c.GetHeader("Auth-Token")
		ip := c.ClientIP()
		utils.Log.Infof("auth=%v ip=%v url=[%v]", auth, ip, c.Request.URL)
		utils.Log.Infof("Params:%v", c.Params)
		// m := make(map[string]interface{})
		// c.BindJSON(&m)
		// log.Println(m)
		c.Next()
	}
}

func Pong(c *gin.Context) {
	utils.Log.Info("got pong!")
	c.JSON(200, "pong!")
}

func ParseData(c *gin.Context, obj interface{}) error {
	data, err := c.GetRawData()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, obj)
	return err
}

func SaveBlog(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	code := c.Param("code")
	b := &model.BlogAutoSave{}
	err := ParseData(c, b)
	if err != nil {
		res.Code = -1
		res.Success = false
		res.Message = err.Error()
	} else {
		err = dao.AutoSaveBlog(code, b.Theme, b.Data, b.AuthorId)
		if err != nil {
			res.Code = -2
			res.Success = false
			res.Message = err.Error()
		}
	}
	c.JSONP(200, res)
}

func LoginBlog(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	req := &model.UserLogin{}
	err := ParseData(c, req)
	if err != nil {
		res.Code = -1
		res.Success = false
		res.Message = err.Error()
	} else {
		secure, _ := utils.EnPwdCode([]byte(req.Password))
		user, err := dao.AuthUser(req.Username, secure)
		if err != nil {
			res.Code = -1
			res.Success = false
			res.Message = "用户名或密码错误!"
		} else {
			res.Data = user
		}
	}
	c.JSONP(200, res)
}

func NewBlogPage(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	blog, err := dao.NewBlog(1)
	if err != nil {
		res.Code = -1
		res.Success = false
		res.Message = err.Error()
	} else {
		res.Data = blog
	}
	c.JSONP(200, res)
}

func GetPosts(c *gin.Context) {
	code := c.Param("code")
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	blog := dao.QueryBlog(0, code)
	if blog == nil {
		res.Code = -1
		res.Success = false
		res.Message = "wrong posts code!"
	} else {
		res.Data = blog
	}
	c.JSONP(200, res)
}

func PutPosts(c *gin.Context) {
	code := c.Param("code")
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	b := &model.BlogAutoSave{}
	err := ParseData(c, b)
	if err != nil {
		res.Code = -1
		res.Success = false
		res.Message = err.Error()
	} else {
		blog := dao.PutBlog(code, b.Data)
		if blog == nil {
			res.Code = -1
			res.Success = false
			res.Message = "wrong posts code!"
		} else {
			res.Data = blog
		}
	}
	c.JSONP(200, res)

}

func GetMenu(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	menu := &model.BookMenu{}
	ms := dao.QueryMenus(0)
	if ms == nil {
		res.Code = -1
		res.Success = false
		res.Message = "get empty menu!"
	} else {
		menu.Id = ms[0].Id
		menu.Name = ms[0].Name
		menu.Chepters = getChapters(menu.Id)
		res.Data = menu
	}
	c.JSONP(200, res)
}

func getChapters(pid int64) []model.BookChapter {
	chapters := []model.BookChapter{}
	ms := dao.QueryMenus(pid)
	for _, item := range ms {
		chapter := model.BookChapter{
			Id:    item.Id,
			Name:  item.Name,
			Books: getBooks(item.Id),
		}
		chapters = append(chapters, chapter)
	}
	return chapters
}

func getBooks(pid int64) []model.BookItem {
	books := []model.BookItem{}
	ms := dao.QueryMenus(pid)
	for _, item := range ms {
		book := model.BookItem{
			Id:    item.Id,
			Title: item.Name,
			Url:   "/page/" + item.Code,
			Day:   item.CreateTime.Format("01,02,2006"),
		}
		books = append(books, book)
	}
	return books
}

func BlogCaches(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
		Data:    dao.QueryBlogCaches(1),
	}
	c.JSONP(200, res)
}

func EditUser(c *gin.Context) {

}
