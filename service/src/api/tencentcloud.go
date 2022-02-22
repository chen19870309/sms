package api

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sms/service/src/config"
	"sms/service/src/dao"
	"sms/service/src/utils"
	"strings"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	tts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tts/v20190823"
)

var ttscClient *tts.Client

const format = `{
	"ModelType":1,
	"Volume":5,
	"Speed":-1.2,
	"VoiceType":101015
}`

func InitTCloudAPI() {
	credential := common.NewCredential(config.TCloud.AppId, config.TCloud.Secret)
	ttscClient, _ = tts.NewClient(credential, regions.Shanghai, profile.NewClientProfile())
}

func TextToVoice(key, text string) error {
	request := tts.NewTextToVoiceRequest()
	request.FromJsonString(format)
	sessionid := utils.Gen8RCode() + fmt.Sprintf("-%v", time.Now().UnixMicro())
	request.SessionId = &sessionid
	request.Text = &text
	res, err := ttscClient.TextToVoice(request)
	if err != nil {
		utils.Log.Error("TextToVoice failed!", err)
		return err
	} else {
		utils.Log.Info("TextToVoice:", res)
		data, err := base64.StdEncoding.DecodeString(*res.Response.Audio)
		if err != nil {
			utils.Log.Error("TextToVoice debase64 failed!", res.ToJsonString())
		} else {
			err = ioutil.WriteFile(config.App.StaticPath+"sound/"+key+".wav", data, 0755)
			if err != nil {
				utils.Log.Error("WriteFile failed!", res.ToJsonString())
			}
		}
		return err
	}
}

//扫描目录生成字库
func AutoGenScope() {
	//GenScopeDataByFile(config.App.BasePath + "常用字.csv")
	//TextToVoice("静夜思", "静夜思 唐，李白 床前明月光，疑是地上霜。举头望明月，低头思故乡")
	//TextToVoice("咏鹅", "咏鹅 唐，骆宾王 鹅鹅鹅，曲项向天歌。白毛浮绿水，红掌拨清波")
	//TextToVoice("A", "A,apple")
	//TextToVoice("B", "B,banana")
}

func GenScopeDataByFile(file string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		utils.Log.Errorf("GenScopeDataByFile(%s) failed! %v", file, err)
		return err
	}
	fname := strings.Split(filepath.Base(file), ".")
	body := "# 小西的字库|" + fname[0] + "\n@create time:" + time.Now().Format("2006-01-02 15:04:05") + "\n@private\n@tag:scope,cardres\n@scope=" + fname[0] + "\n"
	bdata := []byte(body)
	txt := string(b)
	items := strings.Split(txt, "\n")
	for _, item := range items { //ai,t|小,小西的小|image|sound
		utils.Log.Infof("start deal[%s]", item)
		args := strings.Split(item, ",")
		if len(args) != 4 {
			continue
		}
		tk := strings.Split(args[1], "，")
		sound := ""
		if args[0] == "ai" { //调用云平台的ai语言合成
			TextToVoice(tk[0], args[1])
			sound = "https://www.xiaoxibaby.xyz/static/sound/gp1/" + tk[0] + ".m4a"
		} else {
			sound = args[3]
		}
		bstr := []byte("## ")
		bstr = append(bstr, []byte(tk[0]+"\n")...)                                    //标题
		bstr = append(bstr, []byte("> "+args[1]+"\n")...)                             //内容
		bstr = append(bstr, []byte(fmt.Sprintf("[sound](%s)\n", sound))...)           //声音
		bstr = append(bstr, []byte(fmt.Sprintf("![image](%s)\n***\n\n", args[2]))...) //图片

		bdata = append(bdata, bstr...)
	}
	blog, err := dao.NewBlog(1)
	if err != nil {
		utils.Log.Errorf("NewBlog(%s) failed! %v", file, err)
		return err
	}
	dao.PutBlog(blog.Code, string(bdata), 1)
	return nil
}
