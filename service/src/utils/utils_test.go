package utils

import (
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
