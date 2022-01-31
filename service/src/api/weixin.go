package api

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"net/http"
	"sms/service/src/api/handler"
	"sms/service/src/config"
	"sms/service/src/utils"

	"github.com/arstd/weixin"
	"github.com/gin-gonic/gin"
)

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
		keys, err := base64.StdEncoding.DecodeString(config.WX.EncodingAeskey + "=")
		if err != nil {
			return nil, err
		}
		body, err = weixin.DecryptMsg(encMsg.Encrypt, keys, config.WX.AppId)
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
		keys, err := base64.StdEncoding.DecodeString(config.WX.EncodingAeskey + "=")
		if err != nil {
			return nil, err
		}
		b64Enc, err := weixin.EncryptMsg(ret, keys, config.WX.AppId)
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
