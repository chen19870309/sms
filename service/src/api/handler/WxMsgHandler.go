package handler

import (
	"encoding/json"
	"sms/service/src/utils"

	"github.com/arstd/weixin"
)

func DefaultHandler(msg *weixin.Message) weixin.ReplyMsg {
	utils.Log.Infof("defaultHandler:%+v", msg)

	event := weixin.NewRecvEvent(msg)
	js, _ := json.Marshal(event)

	// echo message
	ret := &weixin.ReplyText{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   msg.CreateTime,
		Content:      string(js),
	}

	utils.Log.Infof("replay message: %+v", ret)
	return ret
}

func EchoMsgText(m *weixin.RecvText) weixin.ReplyMsg {
	utils.Log.Infof("receive message: %+v", m)

	// echo message
	ret := &weixin.ReplyText{
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   m.CreateTime,
		Content:      m.FromUserName + ", " + m.Content,
	}

	utils.Log.Infof("replay message: %+v", ret)
	return ret
}

func EchoMsgImage(m *weixin.RecvImage) weixin.ReplyMsg {
	utils.Log.Infof("%+v", m)

	// echo message
	ret := &weixin.ReplyImage{
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   m.CreateTime,
		PicUrl:       m.PicUrl,
		MediaId:      m.MediaId,
	}

	utils.Log.Infof("%+v", ret)
	return ret
}

func EchoMsgVoice(m *weixin.RecvVoice) weixin.ReplyMsg {
	utils.Log.Infof("%+v", m)

	// echo message
	ret := &weixin.ReplyVoice{
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   m.CreateTime,
		MediaId:      m.MediaId,
	}

	utils.Log.Infof("%+v", ret)
	return ret
}

func EchoMsgVideo(m *weixin.RecvVideo) weixin.ReplyMsg {
	utils.Log.Infof("%+v", m)

	// MediaId ???
	ret := &weixin.ReplyVideo{
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   m.CreateTime,
		MediaId:      m.MediaId,
		Title:        "video",
		Description:  "thist is a test desc...",
	}

	utils.Log.Infof("%+v", ret)
	return ret
}

func EchoMsgShortVideo(m *weixin.RecvVideo) weixin.ReplyMsg {
	utils.Log.Infof("%+v", m)

	// MediaId ???
	ret := &weixin.ReplyVideo{
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   m.CreateTime,
		MediaId:      m.ThumbMediaId,
		Title:        "shortvideo",
		Description:  "thist is a test desc...",
	}

	utils.Log.Infof("%+v", ret)
	return ret
}

func EchoMsgLocation(m *weixin.RecvLocation) weixin.ReplyMsg {
	utils.Log.Infof("%+v", m)

	// echo message
	ret := &weixin.ReplyText{
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   m.CreateTime,
		Content:      weixin.AccessToken(),
	}

	utils.Log.Infof("replay message: %+v", ret)
	return ret
}

func EchoMsgLink(m *weixin.RecvLink) weixin.ReplyMsg {
	utils.Log.Infof("%+v", m)

	// 回复图文消息

	return nil
}
