package api

import (
	"context"
	"sms/service/src/api/model"
	"sms/service/src/config"
	"sms/service/src/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/storage"
)

func ServiceUpload(key string) {
	base := "/Users/chenchunjiang/go/src/sms/webapp/src/assets/img/"
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": key,
		},
	}
	upToken := config.GetSimpleUpToken()
	utils.Log.Info("Qiniu GetSimpleUpToken:", upToken)
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, base+key, &putExtra)
	if err != nil {
		utils.Log.Error("Qiniu PutFile:", err)
		return
	}
	utils.Log.Info(ret.Key, "|", ret.Hash)
}

func QiniuUpToken(c *gin.Context) {
	data, ok := utils.GetCache("QiniuUpToken")
	if ok {
		res := data.(model.Response)
		c.JSONP(200, res)
		return
	}
	tk := make(map[string]string)
	tk["token"] = config.GetSimpleUpToken()
	tk["create_time"] = utils.GetStdTime()
	tk["domain"] = config.Qiniu.Domain
	tk["bucket"] = config.Qiniu.Bucket
	tk["prefix"] = time.Now().Format("20060102_")
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
		Data:    tk,
	}
	utils.SetCache("QiniuUpToken", res, 2*time.Minute)
	c.JSONP(200, res)
}
