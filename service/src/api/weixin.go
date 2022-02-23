package api

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"sms/service/src/api/handler"
	"sms/service/src/api/model"
	"sms/service/src/config"
	"sms/service/src/dao"
	"sms/service/src/utils"

	"github.com/arstd/weixin"
	"github.com/gin-gonic/gin"
)

var aeskey []byte

func InitWeixinService(web *WebS) {
	weixin.Initialize(config.WX.OrignId, config.WX.AppId, config.WX.AppSecret, config.WX.Token, config.WX.EncodingAeskey)

	weixin.RecvTextHandler = handler.EchoMsgText             // 注册文本消息处理器
	weixin.RecvImageHandler = handler.EchoMsgImage           // 注册图片消息处理器
	weixin.RecvVoiceHandler = handler.EchoMsgVoice           // 注册语音消息处理器
	weixin.RecvVideoHandler = handler.EchoMsgVideo           // 注册视频消息处理器
	weixin.RecvShortVideoHandler = handler.EchoMsgShortVideo // 注册小视频消息处理器
	weixin.RecvLocationHandler = handler.EchoMsgLocation     // 注册位置消息处理器
	weixin.RecvLinkHandler = handler.EchoMsgLink             // 注册链接消息处理器
	weixin.RecvDefaultHandler = handler.DefaultHandler       // 注册默认处理器

	weixin.EventSubscribeHandler = handler.EventSubscribeHandler     // 注册关注事件处理器
	weixin.EventUnsubscribeHandler = handler.EventUnsubscribeHandler // 注册取消关注事件处理器
	weixin.EventLocationHandler = handler.EventLocationHandler       // 注册上报地理位置事件处理器
	weixin.EventClickHandler = handler.EventClickHandler             // 注册点击自定义菜单事件处理器
	weixin.EventViewHandler = handler.EventViewHandler               // 注册点击菜单跳转链接时的事件处理器
	// 模版消息发送结果通知事件
	weixin.EventTemplateSendJobFinishHandler = handler.EventTemplateSendJobFinishHandler
	weixin.EventDefaultHandler = handler.EventDefaultHandler // 注册默认处理器
	wx := service.Serv.Group("weixin")
	wx.GET("", CheckWeixin)
	wx.POST("", HandleWxMesage)
	wx.GET("/scopes", DealScopes)               //获取字库分类
	wx.GET("/words", DealGetWords)              //获取字库
	wx.POST("/words", DealCheckWords)           //更新个人字库
	wx.GET("/diary", DealGetDiary)              //获取个人日记
	wx.POST("/diary", DealPutDiary)             //插入日记
	wx.POST("/:appid/userinfo", DealWXUserInfo) //小程序用户上传用户信息
	wx.POST("/:appid/login", DealLoginWX)       //小程序用户登陆
}

func DealWXUserInfo(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	appid := c.Param("appid")
	params := make(map[string]interface{})
	body, err := c.GetRawData()
	if err == nil {
		err = json.Unmarshal(body, &params)
		utils.Log.Infof("DealWXUserInfo:%v|%v", appid, params)
		//if reflect.TypeOf(params["userid"]) == "string"
		userid := params["userid"].(float64)
		data := params["data"].(map[string]interface{})
		u, err := dao.UpdateUserInfo(data["avatarUrl"].(string), data["nickName"].(string), string(body), int(userid))
		if err != nil {
			utils.Log.Error("UpdateUserInfo failed!", err)
		} else {
			res.Jwt = utils.GenJwt(u.Id, u.Username, config.App.Secret)
			res.Data = u
		}
	}
	c.JSONP(200, res)
}

func DealGetWords(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	data := []*model.Word{}
	ls, err := dao.QueryScopedCard(c.Query("userid"), c.Query("scope"), "words", c.Query("group"), c.Query("word"), c.Query("mode"), 10)
	if err != nil {
		res.Success = false
		res.Message = err.Error()
	} else {
		for _, a := range ls {
			item := &model.Word{
				Id:     int(a.Id),
				Word:   a.Word,
				PinYin: a.Pinyin,
				Pic:    a.Pic,
				Sound:  a.Sound,
				Scope:  a.Scope,
				Group:  a.Gp,
			}
			data = append(data, item)
		}
		res.Data = data
		res.Count = len(data)
	}
	utils.Log.Info("DealGetWords result = >", res)
	c.JSONP(200, res)
}

func DealCheckWords(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	check := &model.CheckWord{}
	err := ParseData(c, check)
	if err == nil {
		if check.Id > 0 && check.Userid > 0 {
			//check.Status = 0 生字本 1 熟悉字本
			err = dao.EditUserCardRes(check.Userid, check.Id, check.Status)
		} else {
			err = errors.New("check params failed!")
		}
	}
	if err != nil {
		res.Success = false
		res.Message = err.Error()
	}
	c.JSONP(200, res)
}

func DealScopes(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	res.Data = dao.GetResScopes(c.Query("userid"), "words")
	c.JSONP(200, res)
}

func DealLoginWX(c *gin.Context) {
	res := model.Response{
		Code:    0,
		Success: true,
		Message: "ok",
	}
	appid := c.Param("appid")
	params := make(map[string]string)
	body, err := c.GetRawData()
	utils.Log.Infof("GetRawData:%v", err)
	err = json.Unmarshal(body, &params)
	utils.Log.Infof("DealLoginWX:%v|%v", appid, params)
	secret := "f9e9734c2310119423c9cf726c5bc209"
	openid, err := utils.GetUserOpenidByCode(appid, secret, params["code"])
	utils.Log.Infof("GetUserOpenidByCode:%v|%v", openid, err)
	if err == nil {
		user, err := utils.GetUserInfoByOpenId(appid, secret, openid)
		utils.Log.Infof("GetUserInfo:%v|%v", user, err)
		if err != nil {
			u, err := dao.SaveWxUser(openid, "", "", "")
			res.Jwt = utils.GenJwt(u.Id, u.Username, config.App.Secret)
			res.Data = u
			utils.Log.Infof("SaveWxUser:%v|%v", u, err)
		} else {
			b, _ := json.Marshal(user)
			u, err := dao.SaveWxUser(openid, user["nickname"].(string), user["headimgurl"].(string), string(b))
			res.Jwt = utils.GenJwt(u.Id, u.Username, config.App.Secret)
			res.Data = u
			utils.Log.Infof("SaveWxUser:%v|%v", u, err)
		}
	}
	c.JSONP(200, res)
}

func check(c *gin.Context) bool {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")

	// 每次都验证 URL，以判断来源是否合法
	return weixin.ValidateURL(config.WX.Token, timestamp, nonce, signature)
}

func CheckWeixin(c *gin.Context) {
	// 每次都验证 URL，以判断来源是否合法
	if !check(c) {
		utils.Log.Error("CheckWeixin failed!", c.ClientIP())
		c.AbortWithStatus(http.StatusNonAuthoritativeInfo)
	} else {
		utils.Log.Info("CheckWeixin OK!", c.ClientIP(), "|", c.Query("echostr"))
		c.Writer.Write([]byte(c.Query("echostr")))
	}
}

func HandleWxMesage(c *gin.Context) {
	// 每次都验证 URL，以判断来源是否合法
	if !check(c) {
		utils.Log.Error("CheckWeixin failed!", c.ClientIP())
		c.AbortWithStatus(http.StatusNonAuthoritativeInfo)
	} else {
		msg, err := parseBody(c)
		if err != nil {
			utils.Log.Error("parseBody failed!", err)
			c.AbortWithError(http.StatusBadRequest, err)
		} else {
			utils.Log.Info("HandleMessage !", msg)
			reply := weixin.HandleMessage(msg)
			// 如果返回为 nil，则默认返回""
			ret := []byte("")
			if reply != nil {
				// 如果返回不为 nil，表示需要回复
				ret, err = packReply(reply, c.Query("encrypt_type"), c.Query("timestamp"), c.Query("nonce"))
				if err != nil {
					c.AbortWithError(http.StatusBadRequest, err)
					return
				}
			}
			utils.Log.Infof("to weixin: %s", string(ret))
			c.Header("Content-Type", "text/xml; charset=utf-8")
			c.Writer.Write(ret)
		}
	}
}

func parseBody(c *gin.Context) (msg *weixin.Message, err error) {
	body, err := c.GetRawData()
	msg = &weixin.Message{}
	// 如果报文被加密了，先要验签解密
	if c.Query("encrypt_type") == "aes" {
		encMsg := &weixin.EncMessage{}
		// 解析加密的 xml
		err = xml.Unmarshal(body, encMsg)
		if err != nil {
			return nil, err
		}
		msg.ToUserName = encMsg.ToUserName
		msg.Encrypt = encMsg.Encrypt

		if !weixin.CheckSignature(config.WX.Token, c.Query("timestamp"), c.Query("nonce"), encMsg.Encrypt, c.Query("msg_signature")) {
			return nil, errors.New("check signature error")
		}
		body, err = weixin.DecryptMsg(encMsg.Encrypt, getAesKey(), config.WX.AppId)
		if err != nil {
			return nil, err
		}
		utils.Log.Infof("receive: %s", body)
	}

	// 解析 xml
	err = xml.Unmarshal(body, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func getAesKey() []byte {
	if aeskey == nil || len(aeskey) != 32 {
		aeskey, _ = base64.StdEncoding.DecodeString(config.WX.EncodingAeskey + "=")
	}
	return aeskey
}

func packReply(reply weixin.ReplyMsg, encryptType, timestamp, nonce string) (ret []byte, err error) {
	switch reply.(type) {
	case *weixin.ReplyText:
		reply.SetMsgType(weixin.MsgTypeText)
	case *weixin.ReplyImage:
		reply.SetMsgType(weixin.MsgTypeImage)
	case *weixin.ReplyVoice:
		reply.SetMsgType(weixin.MsgTypeVoice)
	case *weixin.ReplyVideo:
		reply.SetMsgType(weixin.MsgTypeVideo)
	case *weixin.ReplyMusic:
		reply.SetMsgType(weixin.MsgTypeMusic)
	case *weixin.ReplyNews:
		reply.SetMsgType(weixin.MsgTypeNews)
	default:
		utils.Log.Error("unexpected custom message type")
	}

	ret, err = xml.MarshalIndent(reply, "", "  ")
	if err != nil {
		return nil, err
	}
	utils.Log.Infof("replay: %s", ret)

	// 如果接收的消息加密了，那么回复的消息也需要签名加密
	if encryptType == "aes" {
		b64Enc, err := weixin.EncryptMsg(ret, getAesKey(), config.WX.AppId)
		if err != nil {
			return nil, err
		}
		encMsg := weixin.EncMessage{
			Encrypt:      b64Enc,
			MsgSignature: weixin.Signature(config.WX.Token, timestamp, nonce, b64Enc),
			TimeStamp:    timestamp,
			Nonce:        nonce, // 随机数
		}
		ret, err = xml.MarshalIndent(encMsg, "", "    ")
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}
