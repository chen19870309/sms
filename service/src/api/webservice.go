package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sms/service/src/api/model"
	"sms/service/src/dao"
	"sms/service/src/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type WebS struct {
	Serv     *gin.Engine
	SecCache map[string]*model.UserData
}

const _USERDATA = "_USERDATA_"
const def = "http://r5uiv7l5f.hd-bkt.clouddn.com/bj1.jpeg"

var service *WebS

func init() {
	service = &WebS{
		SecCache: make(map[string]*model.UserData),
	}
	service.Serv = gin.Default()
	service.Serv.Use(Cors())
	service.Serv.Use(BlogAuth())
	service.Serv.StaticFS("/static", http.Dir("/Users/chenchunjiang/go/src/sms/webapp/dist/static"))
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
	InitEngine()
	blog := service.Serv.Group("blog")
	blog.GET("/ping", Pong)
	blog.GET("/menu", GetMenu)
	blog.POST("/login", LoginBlog)
	blog.POST("/checkemail", CheckEmail)
	blog.POST("/regist", RegistBlog)
	blog.POST("/edituser", EditUser)
	blog.POST("/editpwd", EditPwd)
	blog.GET("/uptoken", QiniuUpToken)
	blog.GET("/newblog", NewBlogPage)
	blog.GET("/blogcaches", BlogCaches)
	blog.POST("/save/:code", SaveBlog)
	blog.GET("/posts/:code", GetPosts)
	blog.PUT("/posts/:code", PutPosts)
	blog.GET("/mainindex", IndexBlog)
	blog.GET("/usereditindex", IndexUserEdit)
	blog.GET("/more", MoreBlog)
	blog.POST("/search", SearchBlog)
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
		ok := true
		auth, err := c.Cookie("auth_token")
		ip := c.ClientIP()
		url := c.Request.URL.String()
		if strings.Contains(url, "newblog") || strings.Contains(url, "save") || strings.Contains(url, "blogcaches") || strings.Contains(url, "edit") {
			user := service.SecCache[auth]
			if user != nil {
				utils.Log.Infof("auth=%v ip=%v url=[%v] userid=[%d]", auth, ip, c.Request.URL, user.Id)
				utils.Log.Infof("Params:%v", c.Params)
			}
			if user == nil || user.Id == 0 {
				ok = false
				res := model.Response{
					Code:    401,
					Success: false,
					Message: "auth failed!",
				}
				c.JSONP(200, res)
				c.AbortWithStatus(401)
			} else {
				c.Set(_USERDATA, user)
				d, e := c.Get(_USERDATA)
				utils.Log.Info("USER=", d, e)
			}
		}
		// m := make(map[string]interface{})
		// c.BindJSON(&m)
		// log.Println(m)
		if ok {
			c.Next()
		}
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
		user, ok := c.Get(_USERDATA)
		if ok {
			b.AuthorId = uint(user.(*model.UserData).Id)
			err = dao.AutoSaveBlog(code, b.Theme, b.Data, b.AuthorId)
		} else {
			err = errors.New("用户授权信息已过期")
		}
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
		secure, _ := utils.PwdCode(req.Password)
		user, err := dao.AuthUser(req.Username, secure, c.ClientIP())
		if err != nil {
			res.Code = -1
			res.Success = false
			res.Message = "用户名或密码错误!"
		} else {
			res.Data = user
			service.SecCache[user.Secret] = &model.UserData{
				Id:       user.Id,
				Username: user.Username,
				Nickname: user.Nickname,
				Remark:   user.Remark,
			}
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
	user, ok := c.Get(_USERDATA)
	var err error
	if ok {
		res.Data, err = dao.NewBlog(uint(user.(*model.UserData).Id))
	} else {
		err = errors.New("用户授权已过期")
	}
	if err != nil {
		res.Code = -1
		res.Success = false
		res.Message = err.Error()
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
		blog := dao.PutBlog(code, b.Data, b.AuthorId)
		if blog == nil {
			res.Code = -1
			res.Success = false
			res.Message = "wrong posts code!(文章发布失败)"
		} else {
			res.Data = blog
			go AddIndex(blog) //刷新索引
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
	userid := 0
	user, ok := c.Get(_USERDATA)
	if ok {
		userid = int(user.(*model.UserData).Id)
	}
	ms := dao.QueryMenus(0, userid)
	if ms == nil {
		res.Code = -1
		res.Success = false
		res.Message = "get empty menu!"
	} else {
		menu.Id = ms[0].Id
		menu.Name = ms[0].Name
		menu.Chepters = getChapters(menu.Id, userid)
		res.Data = menu
	}
	c.JSONP(200, res)
}

func getChapters(pid int64, userid int) []model.BookChapter {
	chapters := []model.BookChapter{}
	ms := dao.QueryMenus(pid, userid)
	for _, item := range ms {
		chapter := model.BookChapter{
			Id:    item.Id,
			Name:  item.Name,
			Books: getBooks(item.Id, userid),
		}
		chapters = append(chapters, chapter)
	}
	return chapters
}

func getBooks(pid int64, userid int) []model.BookItem {
	books := []model.BookItem{}
	ms := dao.QueryMenus(pid, userid)
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
	userid := 0
	user, ok := c.Get(_USERDATA)
	if ok {
		userid = int(user.(*model.UserData).Id)
	}
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
		Data:    dao.QueryBlogCaches(userid),
	}
	c.JSONP(200, res)
}

func EditPwd(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	b := &model.UserData{}
	userid := 0
	user, ok := c.Get(_USERDATA)
	if ok {
		userid = int(user.(*model.UserData).Id)
	}
	err := ParseData(c, b)
	if err != nil {
		res.Code = -1
		res.Success = false
		res.Message = err.Error()
	} else {
		op, _ := utils.PwdCode(b.OP)
		np, _ := utils.PwdCode(b.NP)
		utils.Log.Info("ExchangeUserPwd by user:", userid)
		_, err = dao.ExchangeUserPwd(op, np, userid)
		if err != nil {
			res.Code = -2
			res.Success = false
			res.Message = err.Error()
		}
	}
	c.JSONP(200, res)
}

func EditUser(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	b := &model.UserData{}
	userid := 0
	user, ok := c.Get(_USERDATA)
	if ok {
		userid = int(user.(*model.UserData).Id)
	}
	err := ParseData(c, b)
	if err != nil {
		res.Code = -1
		res.Success = false
		res.Message = err.Error()
	} else {
		_, err = dao.UpdateUserInfo(b.Icon, b.Nickname, b.Remark, userid)
		if err != nil {
			res.Code = -2
			res.Success = false
			res.Message = err.Error()
		}
	}
	c.JSONP(200, res)
}

func SearchBlog(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	b := &model.SearchReq{}
	ParseData(c, b)
	ls, err := Search(b.Text)
	if err != nil {
		res.Code = -1
		res.Success = false
		res.Message = err.Error()
	} else {
		rs := []*model.BookItem{}
		for _, id := range ls {
			item := dao.QueryBlog(int64(id), "")
			if item != nil {
				data := model.BookItem{
					Id:    item.Id,
					Code:  item.Code,
					Title: item.Title,
					Sum:   utils.GetSum(item.Content),
					Pic:   utils.GetPic(item.Content, def),
					Url:   utils.GetBookUrl(item.Code),
				}
				rs = append(rs, &data)
			}
		}
		res.Data = rs
	}
	c.JSONP(200, res)
}

func WaperBlogs(userid, page, pageSize int) []*model.BookItem {
	utils.Log.Infof("WaperBlogs(%d,%d,%d)", userid, page, pageSize)
	rs := []*model.BookItem{}
	for _, item := range dao.GetUserBlogs(userid, page, pageSize) {
		rs = append(rs, &model.BookItem{
			Id:    item.Id,
			Code:  item.Code,
			Title: item.Title,
			Sum:   utils.GetSum(item.Content),
			Pic:   utils.GetPic(item.Content, def),
			Url:   utils.GetBookUrl(item.Code),
			Day:   "",
		})
	}
	return rs
}

func IndexBlog(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
		Data:    WaperBlogs(0, 1, 8),
	}
	c.JSONP(200, res)
}

func IndexUserEdit(c *gin.Context) {
	index := 1
	page_index, err := c.Cookie("user_page_index")
	utils.Log.Info("IndexUserEdit:", page_index, err)
	if page_index != "" {
		index, _ = strconv.Atoi(page_index)
	}
	userid := 0
	user, ok := c.Get(_USERDATA)
	if ok {
		userid = int(user.(*model.UserData).Id)
	}
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
		Data:    WaperBlogs(userid, index, 10),
		Count:   dao.CountUserBlogs(userid),
	}
	c.JSONP(200, res)
}

func RegistBlog(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	b := &model.RegistData{}
	err := ParseData(c, b)
	if err != nil {
		res.Code = -1
		res.Success = false
		res.Message = err.Error()
	} else {
		secure, _ := utils.PwdCode(b.Password)
		user, err := dao.RegistUser(b.Username, secure, b.Email, b.Code, c.ClientIP())
		if err != nil {
			res.Code = -1
			res.Success = false
			res.Message = err.Error()
		} else {
			res.Data = user
			service.SecCache[user.Secret] = &model.UserData{
				Id:       user.Id,
				Username: user.Username,
				Nickname: user.Nickname,
				Remark:   user.Remark,
			}
		}
	}
	c.JSONP(200, res)
}

func CheckEmail(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	b := &model.RegistData{}
	err := ParseData(c, b)
	if err != nil {
		res.Code = -1
		res.Success = false
		res.Message = err.Error()
	} else {
		code, err := dao.GenEmailCode(b.Email)
		if err != nil {
			res.Code = -1
			res.Success = false
			res.Message = err.Error()
		} else {
			//发送邮件 b.Email
			utils.Log.Info("Email:", code)
			//utils.SendMailWithCode(b.Email, code)
		}
	}
	c.JSONP(200, res)
}

func MoreBlog(c *gin.Context) {
}
