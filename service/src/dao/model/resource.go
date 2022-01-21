package model

import (
	"time"
)

type Resource struct {
	Id         int64     `json:"id"`          // id integer not null auto_increment,
	ResType    string    `json:"res_type"`    // res_type varchar(64) not null default 'pic',
	Uri        string    `json:"uri"`         // uri varchar(128) not null default '',
	ResVal     string    `json:"res_val"`     // res_val varchar(64) not null default '',
	CreateTime time.Time `json:"create_time"` // create_time datetime default CURRENT_TIMESTAMP ,
	ExpireTime time.Time `json:"expire_time"` // expire_time datetime default CURRENT_TIMESTAMP,
}
