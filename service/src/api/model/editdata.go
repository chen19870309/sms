package model

type UserData struct {
	Id       int64  `json:"id"`
	OP       string `json:"OP"`
	NP       string `json:"NP"`
	Icon     string `json:"icon"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Remark   string `json:"remark"`
}
