package utils

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

const mdata = `# Title
## hello word!
![image](http://localhost/a.jpg)

1. 改造数据库表
    1. real_personal_user 个人用户实名表
    2. real_check_record 运单实名表
2. 分库分表策略
    1. mysql 8主8从
    2. 用户实名信息表根据用户userid做hash分库分表
    3. 快件实名表根据bill_code做hash分库分表
3. 应用改造升级流程
    1. 步骤1: 申请MQ，异步将oralce中的历史数据导入到mysql分库额分表中
    2. 步骤2: 读写oracle，异步写入mysql数据
    3. 步骤3: oracle作为主库，实现oracle，mysql双写，mysql 作为查询库，mysql查不到数据时降级到从oracle读取，此过程会持续三个月
    4. 步骤4: 关闭双写逻辑, mysql 作为查询库，查不到时也不从oracle读取
	5. 步骤5: 去除oracle部分代码，下线oracle数据库

id | code
---|---
自增主键 | 用户code
serial | varchar(8) 

`

func TestCode(t *testing.T) {
	t.Error(Gen8RCode())
}

func TestMarkdown(t *testing.T) {
	t.Error(TemplateMarkDown(mdata))
}

func TestRand(t *testing.T) {
	t.Error(GetRandNum(6))
	t.Error(GetRandNum(8))
	t.Error(GetRandNum(9))
}

func TestPic(t *testing.T) {
	text := "# 正则表达式简明教程\n@sum:介绍regexp的使用方法和作用\n![image](http://r5uiv7l5f.hd-bkt.clouddn.com/testregexp.webp)\n\n## 1.Regexp是什么？"
	t.Error(GetPic(text, "null"))
	t.Error(GetSum(text))
}

func TestDec(t *testing.T) {

	data := "YhCbGDqYf5tCSW2UcnXAMJfCV+DOlEdepEzqOEDa9DWUKz+1F+DI33Shnw+To9SOxmvdMGh/abY+LQojth1WSpnzC3Cab8HSS7NmXijMthsl1RCfR2UPdVjqFKZN9/l63QJX9cpugFMQuf6LiB1lp87XThCyKLUh/Vb/h6RUpeheehkOdcH4enS6RrPpFVvaR2bwtMHRf2aV5mBObrIWQLzqoIxgE5pFhASzI7Mcy9rmnJ1wQxJ+pYj0UPClTvPzDFjWLMaN2moHJgh6DAqrVc8VPoxyq1ZZBmReA8qretGcok1tmHQTkR310WNarSJ46pru8glaUVjX0gGKjTTmqeFFD7RZrOkBqJm6wWizrhi9aScfUxfutw2BG+2O/ZWSzQv31Ze/a9p49dViktn9JnW+lRtpkusMPpTSrZ1t8uEPTEyBlHJlAcxtrIWt+e4PXHxpk7ctkvt0xcEv6SjKUw=="
	iv := "xR6vDQfaX6QgREJ1tJ6wUQ=="
	key := "hoO3uGXJDrl/oN/SzmWJ9g=="

	b1, _ := base64.StdEncoding.DecodeString(data)
	k, _ := base64.StdEncoding.DecodeString(key)
	i, _ := base64.StdEncoding.DecodeString(iv)
	t.Error(i)
	res, err := AesCBCDncrypt(b1, k, i)
	t.Error(string(res))
	t.Error(err)
}

var data = `{"userid":5,"data":{"nickName":"å¾®ä¿¡ç<94>¨æ<88>·","gender":0,"language":"","city":"","province":"","country":"","avatarUrl":"https://thirdwx.qlogo.cn/mmopen/vi_32/POgEwh4mIHO4nibH0KlMECNjjGxQUq24ZEaGT4poC6icRiccVGKSyXwibcPq4BWmiaIGuG1icwxaQX6grC9VemZoJ8rg/132"}}`

func TestData(t *testing.T) {
	params := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &params)
	t.Error(err)
	userid := params["userid"].(float64)
	data := params["data"].(map[string]interface{})
	t.Error(userid)
	t.Error(data["avatarUrl"].(string), data["nickName"].(string))
}
