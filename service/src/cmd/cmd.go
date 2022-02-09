package main

import (
	"fmt"
	"sms/service/src/config"
	"sms/service/src/utils"
	"strconv"
	"strings"
	"sync"

	"sms/service/src/api"

	"github.com/judwhite/go-svc"
	"github.com/urfave/cli/v2"
)

type program struct {
	once sync.Once
	blog *api.WebS
	ctx  *cli.Context
}

func (p *program) Init(env svc.Environment) error {
	config.InitConf(p.ctx.String("conf"))
	if config.App.ListenPort == 0 {
		addr := p.ctx.String("address")
		addrs := strings.Split(addr, ":")
		config.App.ListenPort, _ = strconv.ParseInt(addrs[1], 10, 64)
	}
	blog, err := api.NewBlogService()
	if err != nil {
		return err
	}
	if config.WX.OrignId != "" {
		api.InitWeixinService(blog)
	}
	p.blog = blog
	return nil
}

func (p *program) Start() error {
	addr := fmt.Sprintf(":%d", config.App.ListenPort)
	utils.Log.Info("Address:", addr)
	go p.blog.Serv.Run(addr)
	return nil
}

func (p *program) Stop() error {
	p.once.Do(func() {
		// exit
		utils.Log.Info("call stop!")
		api.CloseEngine()
	})
	return nil
}
