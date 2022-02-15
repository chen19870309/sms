package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

var DB database
var Site site
var App app

var Qiniu qiniu
var Mail mail

var WX weixin

var TCloud tencentCloud

type Conf struct {
	Db     database     `json:"db"`
	Site   site         `json:"site"`
	App    app          `json:"app"`
	Qiniu  qiniu        `json:"qiniu"`
	Mail   mail         `json:"mail"`
	Weixin weixin       `json:"weixin"`
	TCloud tencentCloud `json:"tcloud"`
}

type database struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
	DbName   string `json:"dbname"`
	Driver   string `json:"driver"`
}

type site struct {
	Id     int    `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
	BeiAn  string `json:"beian"`
}

type app struct {
	ListenPort int64  `json:"listen-port"`
	LogLevel   string `json:"log-level"`
	BasePath   string `json:"base-path"`
	StaticPath string `json:"static-path"`
	Secret     string `json:"secret"`
}

type qiniu struct {
	Ak     string `json:"ak"`
	Sk     string `json:"sk"`
	Cb     string `json:"cb"`
	Domain string `json:"domain"`
	Bucket string `json:"bucket"`
}

type mail struct {
	FromName  string `json:"from-name"`
	FromEmail string `json:"from-email"`
	FromSec   string `json:"from-sec"`
	Smtp      string `json:"smtp"`
	SmtpPort  int    `json:"smtp-port"`
}

type weixin struct {
	OrignId        string `json:"origin-id"`
	AppId          string `json:"app-id"`
	AppSecret      string `json:"app-secret"`
	Token          string `json:"token"`
	EncodingAeskey string `json:"encoding-aeskey"`
}

type tencentCloud struct {
	AppId  string `json:"app-id"`
	Secret string `json:"secret"`
}

func InitApp(port int64, level, base, static string) {
	if port <= 0 {
		port = 8080
	}
	if level == "" {
		level = "info"
	}
	if base == "" {
		base = "/tmp/blog"
	}
	if static == "" {
		base = "/tmp/blog/static"
	}
	App = app{
		ListenPort: port,
		LogLevel:   level,
		BasePath:   base,
		StaticPath: static,
	}
}

func InitSite(id int, code, name, domain, beian string) {
	if id <= 0 {
		id = 1
	}
	if code == "" {
		code = "0001"
	}
	if name == "" {
		name = "Sms Blog"
	}
	if domain == "" {
		domain = "localhost"
	}
	Site = site{
		Id:     id,
		Code:   code,
		Name:   name,
		Domain: domain,
		BeiAn:  beian,
	}
}

func InitDB(port int, host, dbname, username, password, driver string) {
	DB = database{
		Host:     host,
		Port:     port,
		DbName:   dbname,
		UserName: username,
		PassWord: password,
		Driver:   driver,
	}
}

func GetConnArgs() string {
	var conn = ""
	switch DB.Driver {
	case "mysql":
		conn = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", DB.UserName, DB.PassWord, DB.Host, DB.Port, DB.DbName)
		break
	case "postgres":
		conn = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", DB.Host, DB.Port, DB.DbName, DB.UserName, DB.PassWord)
		break
	default:
	}
	return conn
}

func InitConf(conf string) error {
	b, e := ioutil.ReadFile(conf)
	if e != nil {
		log.Panic(e)
		return e
	}
	config := &Conf{}
	e = json.Unmarshal(b, config)
	if e != nil {
		log.Panic(e)
		return e
	}
	DB = config.Db
	Site = config.Site
	App = config.App
	Qiniu = config.Qiniu
	Mail = config.Mail
	WX = config.Weixin
	TCloud = config.TCloud
	return nil
}

func GetSimpleUpToken() string {
	putPolicy := storage.PutPolicy{
		Scope: Qiniu.Bucket,
		//CallbackURL:      qiniu.cb,
		CallbackBodyType: "application/json",
		ReturnBody:       `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	if Qiniu.Cb != "" {
		putPolicy.CallbackURL = Qiniu.Cb
	}
	mac := qbox.NewMac(Qiniu.Ak, Qiniu.Sk)
	return putPolicy.UploadToken(mac)
}
