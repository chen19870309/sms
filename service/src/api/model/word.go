package model

type Word struct {
	Id     int    `json:"Id"`
	Word   string `json:"Word"`
	PinYin string `json:"PinYin"`
	Pic    string `json:"Pic"`
	Sound  string `json:"Sound"`
	Scope  string `json:"scope"`
	Group  string `json:"group"`
}

type CheckWord struct {
	Word
	OpenId string `json:"openid"`
	Userid int    `json:"userid"`
	Status int    `json:"status"`
}
