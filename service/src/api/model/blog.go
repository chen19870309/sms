package model

type Response struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Jwt     string      `json:"jwt_token"`
	Data    interface{} `json:"data"`
	Count   int         `json:"count"`
}

type SearchReq struct {
	Text string `json:"text"`
}

type BlogAutoSave struct {
	Data     string `json:"data" validate:"required"`
	Theme    string `json:"theme"`
	AuthorId uint   `json:"author_id" validate:"required"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type BookItem struct {
	Id         int64  `json:"id"`
	Code       string `json:"code"`
	Url        string `json:"url"`
	Title      string `json:"title"`
	Day        string `json:"day"`
	Sum        string `json:"sum"`
	Pic        string `json:"pic"`
	Author     string `json:"author"`
	Status     uint   `json:"status"`
	UpdateTime string `json:"updatetime"`
}

type BookChapter struct {
	Id    int64      `json:"id"`
	Name  string     `json:"name"`
	Sum   string     `json:"sum"`
	Pic   string     `json:"pic"`
	Books []BookItem `json:"books"`
}

type BookMenu struct {
	Id       int64         `json:"id"`
	Name     string        `json:"name"`
	Chepters []BookChapter `json:"chepters"`
}
