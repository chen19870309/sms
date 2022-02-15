package api

import (
	"sms/service/src/api/model"
	"sms/service/src/dao"

	"github.com/gin-gonic/gin"
)

func DealGetDiary(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	res.Data = dao.GetUserMonthDiary(c.Query("userid"), c.Query("year"), c.Query("month"))
	c.JSONP(200, res)
}

func DealPutDiary(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	data := &model.DiaryData{}
	err := ParseData(c, data)
	if err == nil {
		err = dao.NewUserDiary(data.Id, data.Userid, data.Year, data.Month, data.Day, data.Remark)
		if err != nil {
			res.Success = false
			res.Message = err.Error()
		}
	} else {
		res.Success = false
		res.Message = "无效的参数"
	}
	c.JSONP(200, res)
}
