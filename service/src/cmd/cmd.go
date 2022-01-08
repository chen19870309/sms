package main

import (
	"log"
	"sync"

	"sms/service/src/api"

	"github.com/judwhite/go-svc"
	"github.com/urfave/cli/v2"
)

type program struct {
	once sync.Once
	blog *api.BlogS
	ctx  *cli.Context
}

func (p *program) Init(env svc.Environment) error {
	blog, err := api.NewBlog()
	if err != nil {
		return err
	}
	p.blog = blog
	return nil
}

func (p *program) Start() error {
	go p.blog.Serv.Run(p.ctx.String("port"))
	return nil
}

func (p *program) Stop() error {
	p.once.Do(func() {
		// exit
		log.Println("call stop!")
	})
	return nil
}
