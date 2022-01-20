package model

type Response struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Count   int         `json:"count"`
}

type SearchReq struct {
	Text string `json:"text"`
}

type BlogAutoSave struct {
	Data     string `json:"data"`
	Theme    string `json:"theme"`
	AuthorId uint   `json:"author_id"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type BookItem struct {
	Id    int64  `json:"id"`
	Code  string `json:"code"`
	Url   string `json:"url"`
	Title string `json:"title"`
	Day   string `json:"day"`
	Sum   string `json:"sum"`
	Pic   string `json:"pic"`
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
