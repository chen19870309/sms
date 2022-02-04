package model

import "time"

type CardRes struct {
	Id         int64     `json:"id" `         //id serial primary key,
	ResType    string    `json:"res_type"`    //res_type varchar(64) not null default 'word',
	Pic        string    `json:"pic"`         //pic varchar(128) not null default '',
	Sound      string    `json:"sound"`       //sound varchar(128) not null default '',
	Pinyin     string    `json:"pinyin"`      //pinyin varchar(64) not null default '',
	Scope      string    `json:"scope"`       //scope varchar(64) not null default '',
	Gp         string    `json:"gp"`          //gp varchar(32) not null default '',
	Word       string    `json:"word"`        //word
	Data       string    `json:"data"`        //text
	CreateTime time.Time `json:"create_time"` //create_time timestamp default CURRENT_TIMESTAMP ,
	ExpireTime time.Time `json:"expire_time"` //expire_time timestamp default CURRENT_TIMESTAMP
}
