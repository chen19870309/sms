package utils

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	tts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tts/v20190823"
)

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

var client = &http.Client{Transport: tr}

type WxResp struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpireIn    int    `json:"expire_in"`
}

func Json2Map(b []byte) map[string]interface{} {
	res := make(map[string]interface{})
	err := json.Unmarshal(b, &res)
	if err != nil {
		res["errcode"] = -1
		res["errmsg"] = err.Error()
	}
	return res
}

func GetAppToken(appid, secret string) (string, error) {
	key := fmt.Sprintf("%s:%s", appid, secret)
	data, ok := GetCache(key)
	if ok {
		return data.(string), nil
	}
	seedUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appid, secret)
	resp, err := client.Get(seedUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	rs := Json2Map(body)
	if rs["errmsg"] != nil && rs["errmsg"].(string) != "" {
		Log.Error("GetAppToken failed!", string(body), rs)
		return "", errors.New(rs["errmsg"].(string))
	}
	token := rs["access_token"].(string)
	SetCache(key, token, 100*time.Minute)
	return token, nil
}

func GetUserOpenidByCode(appid, secret, code string) (string, error) {
	seedUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appid, secret, code)
	resp, err := client.Get(seedUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	rs := Json2Map(body)
	if rs["errmsg"] != nil && rs["errmsg"].(string) != "" {
		Log.Error("GetUserOpenidByCode failed!", string(body), rs)
		return "", errors.New(rs["errmsg"].(string))
	}
	Log.Info("GetUserOpenidByCode =>", string(body))
	openid := rs["openid"].(string)
	SetCache("SESSION_KEY_"+openid, rs["session_key"], 100*time.Minute)
	return openid, nil
}

func GetUserInfoByOpenId(appid, secret, openid string) (map[string]interface{}, error) {
	token, err := GetAppToken(appid, secret)
	if err != nil {
		return nil, err
	}
	seedUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN", token, openid)
	resp, err := client.Get(seedUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	rs := Json2Map(body)
	if rs["errmsg"] != nil && rs["errmsg"].(string) != "" {
		Log.Error("GetUserOpenidByCode failed!", string(body), rs)
		return nil, errors.New(rs["errmsg"].(string))
	}
	return rs, nil
}

var ttscClient *tts.Client

const format = `{
	"ModelType":1,
	"Volume":5,
	"Speed":-0.8,
	"VoiceType":1001
}`

func init() {
	credential := common.NewCredential("AKIDaz80b3FAoDMEp9DzVjhOoxpRmqd9OM3U", "B5i4AzXnoY1D9oCJ9i01lEQLzNUBaiaB")
	ttscClient, _ = tts.NewClient(credential, regions.Shanghai, profile.NewClientProfile())
}

func TextToVoice(text string) error {
	request := tts.NewTextToVoiceRequest()
	request.FromJsonString(format)
	sessionid := Gen8RCode() + fmt.Sprintf("-%v", time.Now().UnixMicro())
	request.SessionId = &sessionid
	request.Text = &text
	res, err := ttscClient.TextToVoice(request)
	if err != nil {
		Log.Error("TextToVoice failed!", err)
		return err
	} else {
		Log.Info("TextToVoice:", res)
		data, err := base64.StdEncoding.DecodeString(*res.Response.Audio)
		if err != nil {
			Log.Error("TextToVoice debase64 failed!", res.ToJsonString())
		} else {
			err = ioutil.WriteFile("./test.wav", data, 0755)
		}
		return err
	}
}
