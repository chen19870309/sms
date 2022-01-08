#  个人Blog建表SQL
##  1.blog信息表

id | code | create_time | update_time | author_id | title | tags | sum | content | status 
---|---|---|---|---|---|---|---|---|---
自增主键 | 文件访问码 |创建时间|更新时间|作者id|标题|标签|摘要|md内容|状态
serial | varchar(8) | timestamp | timestamp | integer|varchar(128)|varchar(128)|varchar(128)|text|smallint

```
create table tb_blog_ctx(
    id serial primary key,
    code varchar(8),
    create_time timestamp,
    update_time timestamp,
    author_id integer,
    title varchar(128),
    tags varchar(128),
    sum varchar(128),
    content text,
    status smallint
);
```

## 2.用户信息表

id | code | create_time | update_time | username | nickname | secret | icon | remark | level | status 
---|---|---|---|---|---|---|---|---|---|--
自增主键 | 用户code |创建时间|更新时间|用户名|昵称|密钥|头像|备注状态|用户级别|状态
serial | varchar(8) | timestamp | timestamp | varchar(64)|varchar(128)|varchar(128)|varchar(128)|text|smallint|smallint

```
create table tb_sms_user(
    id serial primary key,
    code varchar(8),
    create_time timestamp,
    update_time timestamp,
    username varchar(64),
    nickname varchar(128),
    secret varchar(128),
    icon varchar(128),
    remark text,
    level smallint,
    status smallint
);
```