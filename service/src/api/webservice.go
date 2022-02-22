package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sms/service/src/api/model"
	"sms/service/src/config"
	"sms/service/src/dao"
	"sms/service/src/utils"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
	gomail "github.com/go-mail/gomail"
)

type WebS struct {
	Serv     *gin.Engine
	SecCache map[string]*model.UserData
}

const _USERDATA = "_USERDATA_"
const def = "/bj1.jpeg"

var service *WebS

// prometheus监控指标
var httpHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace:   "http_server",
	Subsystem:   "",
	Name:        "requests_seconds",
	Help:        "Histogram of response latency (seconds) of http handlers.",
	ConstLabels: nil,
	Buckets:     nil,
}, []string{"method", "code", "uri", "ip"})

func init() {
	service = &WebS{
		SecCache: make(map[string]*model.UserData),
	}
	prometheus.MustRegister(httpHistogram)
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
		t0 := time.Now()
		c.Next()
		uris := strings.Split(c.Request.RequestURI, "?")
		httpHistogram.WithLabelValues(
			c.Request.Method,
			fmt.Sprintf("%v", c.Writer.Status()),
			uris[0],
			c.ClientIP(),
		).Observe(time.Since(t0).Seconds())
	}
}

func BlogAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := c.GetRawData()
		if err != nil {
			utils.Log.Error("GetRawData Failed!:", err)
		} else {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // 关键点
		}
		ok := true
		var user *model.UserData
		jwtStr := c.Request.Header["Authorization"]
		//utils.Log.Info("Auth JWT:", jwtStr)
		if len(jwtStr) != 0 && jwtStr[0] != "" { // 使用jwt校验用户登陆状态
			userid, err := utils.CheckJwt(jwtStr[0], config.App.Secret)
			if err != nil {
				utils.Log.Errorf("Jwt[%v] Auth Failed![%v]", jwtStr, err)
			} else {
				//utils.Log.Errorf("Jwt[%v] Auth Success![%v]", jwtStr, text)
				user = &model.UserData{
					Id: userid,
				}
			}
		} else { // 使用Cookie 授权登陆
			auth, _ := c.Cookie("auth_token")
			user = service.SecCache[auth]
			if user == nil || user.Id == 0 {
				u := dao.CheckAuthCode(auth)
				if u != nil {
					service.SecCache[auth] = &model.UserData{
						Id:       u.Id,
						Username: u.Username,
						Nickname: u.Nickname,
						Remark:   u.Remark,
					}
				}
			}
		}
		if user != nil {
			utils.Log.Info("Auth User:", user)
			c.Set(_USERDATA, user)
		}
		url := c.Request.URL.String()
		utils.Log.Infof("[%v]uri=[%v] IP=[%v] req=[%v]", c.Request.Method, url, c.ClientIP(), string(data))
		if strings.Contains(url, "newblog") || strings.Contains(url, "save") || strings.Contains(url, "blogcaches") || strings.Contains(url, "edit") {
			if user == nil {
				ok = false
				res := model.Response{
					Code:    401,
					Success: false,
					Message: "auth failed please login!",
				}
				c.JSONP(401, res)
				c.AbortWithStatus(401)
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

func SendMail(email, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(config.Mail.FromEmail, config.Mail.FromName)) //这种方式可以添加别名，即“XX官方”
	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])
	mailTo := []string{email}
	m.SetHeader("To", mailTo...) //发送给多个用户
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	r := gomail.NewPlainDialer(config.Mail.Smtp, config.Mail.SmtpPort, config.Mail.FromEmail, config.Mail.FromSec)
	return r.DialAndSend(m)
}

func SendMailWithCode(email, code string) error {
	body := fmt.Sprintf("欢迎注册:%s<br /> 请使用<strong><font color=\"#A52A2A\">%s</font></strong>验证邮箱，验证码10分钟内有效<br />如非本人操作请忽略此邮件^_^<br />", config.Mail.FromName, code)
	return SendMail(email, "账号邮箱激活码", body)
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
			res.Jwt = utils.GenJwt(user.Id, user.Username, config.App.Secret)
			res.Data = user
			service.SecCache[user.Secret] = &model.UserData{
				Id:       user.Id,
				Username: user.Username,
				Nickname: user.Nickname,
				Remark:   user.Remark,
			}
			err = dao.NewAuthCode(user.Id, user.Secret)
			if err != nil {
				utils.Log.Errorf("NewAuthCode(%d,%s) => %v", user.Id, user.Secret, err)
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
	userid := 0
	user, ok := c.Get(_USERDATA)
	if ok {
		userid = int(user.(*model.UserData).Id)
	}
	blog := dao.QueryBlog(0, code)
	if blog == nil {
		res.Code = -1
		res.Success = false
		res.Message = "无效的链接!"
	} else {
		if strings.Contains(blog.Tags, "private") && userid != int(blog.AuthorId) {
			res.Code = -1
			res.Success = false
			res.Message = "私有文章，授权失败!"
		}
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
		if strings.Contains(item.Remark, "private") {
			continue
		}
		user := dao.QueryUser(int64(item.AuthorId), "")
		//utils.Log.Info(item.AuthorId, user)
		book := model.BookItem{
			Id:    item.Id,
			Title: item.Name,
			Url:   "/page/" + item.Code,
			Day:   item.CreateTime.Format("2006-01-02"),
		}
		if user != nil {
			book.Author = user.Nickname
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
					Pic:   utils.GetPic(item.Content, config.Qiniu.Domain+def),
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
			Id:         item.Id,
			Code:       item.Code,
			Title:      item.Title,
			Sum:        utils.GetSum(item.Content),
			Pic:        utils.GetPic(item.Content, config.Qiniu.Domain+def),
			Url:        utils.GetBookUrl(item.Code),
			Day:        utils.GetMdTags(item.Content),
			UpdateTime: item.UpdateTime.Format("2006-01-02"),
			Status:     item.Status,
		})
	}
	return rs
}

func IndexBlog(c *gin.Context) {
	more := c.Query("more")
	size := 8
	if more != "" {
		size, _ = strconv.Atoi(more)
	}
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
		Data:    WaperBlogs(0, 1, size),
		Count:   dao.CountUserBlogs(0),
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
		user := dao.QueryUserEmail(b.Email)
		if user != nil {
			res.Code = -1
			res.Success = false
			res.Message = "此邮箱已经注册过账号,不能重复使用!"
		} else {
			code, err := dao.GenEmailCode(b.Email)
			if err != nil {
				res.Code = -1
				res.Success = false
				res.Message = err.Error()
			} else {
				//发送邮件 b.Email
				utils.Log.Info("Email:", code)
				err := SendMailWithCode(b.Email, code)
				if err != nil {
					res.Code = -1
					res.Success = false
					res.Message = err.Error()
				}
			}
		}
	}
	c.JSONP(200, res)
}

func MoreBlog(c *gin.Context) {
}
