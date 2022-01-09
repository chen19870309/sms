package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var DB database
var Site site
var App app

type Conf struct {
	Db   database `json:"db"`
	Site site     `json:"site"`
	App  app      `json:"app"`
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
		break
	default:
	}
	return conn
}

func InitConf(conf string) error {
	b, e := ioutil.ReadFile(conf)
	if e != nil {
		return e
	}
	config := &Conf{}
	e = json.Unmarshal(b, config)
	if e != nil {
		return e
	}
	DB = config.Db
	Site = config.Site
	App = config.App
	return nil
}
