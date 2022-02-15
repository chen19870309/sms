package model

import "time"

type UserDiary struct {
	Id         int64     `json:"id" `          //id serial primary key,
	Userid     int64     `json:"userid" `      //userid integer not null ,
	Year       int64     `json:"year" `        //year integer not null ,
	Month      int64     `json:"month" `       //month integer not null ,
	Day        int64     `json:"day" `         //day integer not null ,
	CreateTime time.Time `json:"create_time" ` //create_time timestamp default CURRENT_TIMESTAMP ,
	UpdateTime time.Time `json:"update_time" ` //update_time timestamp default CURRENT_TIMESTAMP,
	Status     int64     `json:"status" `      //status int not null default 0,
	Color      string    `json:"color" `       //color varchar(32) not null default 'orange',
	Diary      string    `json:"diary" `       //diary text
}
