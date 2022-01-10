package api

import (
	"bytes"
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
	service.Serv.StaticFile("/", "/Users/chenchunjiang/go/src/sms/webapp/dist/index.html")
}

func (web *WebS) Close() {
	dao.CloseDB()
}

func NewBlog() (w *WebS, err error) {
	dao.InitDB()
	err = dao.TestUser()
	if err != nil {
		utils.Log.Errorf("TestUser Failed![%v]", err)
	}
	err = dao.TestBlog()
	if err != nil {
		utils.Log.Errorf("TestBlog Failed![%v]", err)
	}
	blog := service.Serv.Group("blog")
	blog.GET("/ping", Pong)
	blog.POST("/login", LoginBlog)
	blog.POST("/user/edit", EditUser)
	blog.POST("/save/:code", SaveBlog)
	blog.POST("/posts/:code", GetPosts)
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

func SaveBlog(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	c.JSONP(200, res)
}

func LoginBlog(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
		Data:    dao.QueryUser(1, ""),
	}
	c.JSONP(200, res)
}

func GetPosts(c *gin.Context) {
}

func EditUser(c *gin.Context) {

}
