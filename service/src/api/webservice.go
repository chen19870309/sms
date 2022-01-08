package api

import (
	"github.com/gin-gonic/gin"
)

type BlogS struct {
	Serv *gin.Engine
}

func NewBlog() (b *BlogS, err error) {
	b = &BlogS{}
	router := gin.Default()
	blog := router.Group("blog")
	blog.GET("/ping", Pong)
	b.Serv = router
	return b, nil
}

func Pong(c *gin.Context) {
	c.JSON(200, "pong!")
}
