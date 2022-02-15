package model

import "time"

type UserCardRes struct {
	Id         int64     `json:"id" `
	Userid     int64     `json:"userid"`
	ResId      int64     `json:"res_id"`
	Scope      string    `json:"scope"`
	Gp         string    `json:"group"`
	Word       string    `json:"word"`        //word
	Remark     string    `json:"remark"`      //text
	CreateTime time.Time `json:"create_time"` //create_time timestamp default CURRENT_TIMESTAMP ,
	UpdateTime time.Time `json:"update_time"` //expire_time timestamp default CURRENT_TIMESTAMP
	Status     int       `json:"status"`      //0生字 1 已学会
}

type UserScopeCount struct {
	Id    int
	Scope string
	Gp    string
	Color string
	Icon  string
	Cnt   int
	Ucnt  int
}
