package api

import (
	"log"
	"sms/service/src/dao"

	"github.com/gin-gonic/gin"
)

type WebS struct {
	Serv *gin.Engine
}

var service *WebS

func init() {
	service = &WebS{}
	service.Serv = gin.Default()
}

func (web *WebS) Close() {
	dao.CloseDB()
}

func NewBlog() (w *WebS, err error) {
	dao.InitDB()
	err = dao.TestUser()
	if err != nil {
		log.Println(err)
	}
	err = dao.TestBlog()
	if err != nil {
		log.Println(err)
	}
	blog := service.Serv.Group("blog")
	blog.GET("/ping", Pong)
	blog.POST("/login", LoginBlog)
	blog.POST("/save", SaveBlog)
	blog.POST("/posts/:code", GetPosts)
	return service, nil
}

func Pong(c *gin.Context) {
	c.JSON(200, "pong!")
}

func SaveBlog(c *gin.Context) {
}

func LoginBlog(c *gin.Context) {
}

func GetPosts(c *gin.Context) {
}
