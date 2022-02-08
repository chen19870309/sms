#  个人Blog建表SQL
##  1.blog信息表

id | code | create_time | update_time | author_id | title | tags | sum | content | status 
---|---|---|---|---|---|---|---|---|---
自增主键 | 文件访问码 |创建时间|更新时间|作者id|标题|标签|摘要|md内容|状态
serial | varchar(8) | timestamp | timestamp | integer|varchar(128)|varchar(128)|varchar(128)|text|smallint

```
--postgres
drop table tb_blog_ctx;
create table tb_blog_ctx(
    id serial primary key,
    code varchar(8) default '',
    create_time timestamp default CURRENT_TIMESTAMP,
    update_time timestamp default CURRENT_TIMESTAMP,
    author_id integer default 0,
    title varchar(128) default '',
    tags varchar(128) default 'private',
    sum varchar(128) default '',
    content text,
    status smallint default 0
);
--mysql
create table tb_blog_ctx(
    id bigint(20) not null auto_increment,
    code varchar(8),
    create_time datetime,
    update_time datetime,
    author_id integer,
    title varchar(128),
    tags varchar(128),
    sum varchar(128),
    content text,
    status integer,
    primary key(id)
) default charset=utf8;
```

## 2.用户信息表

id | code | create_time | update_time | username | nickname | secret | icon | remark | level | status 
---|---|---|---|---|---|---|---|---|---|--
自增主键 | 用户code |创建时间|更新时间|用户名|昵称|密钥|头像|备注状态|用户级别|状态
serial | varchar(8) | timestamp | timestamp | varchar(64)|varchar(128)|varchar(128)|varchar(128)|text|smallint|smallint

```
--postgres
drop table tb_sms_user;
create table tb_sms_user(
    id serial primary key,
    code varchar(8) default '',
    create_time timestamp default CURRENT_TIMESTAMP,
    update_time timestamp default CURRENT_TIMESTAMP,
    username varchar(64) not null,
    nickname varchar(128) not null default '<nil>',
    secret varchar(128) not null,
    icon varchar(128) default '',
    remark text ,
    level smallint default 0,
    status smallint default 0,
    email varchar(64) default '',
    login_ip varchar(20) default ''
);
--mysql
create table tb_sms_user(
    id bigint(20) not null auto_increment,
    code varchar(8),
    create_time datetime,
    update_time datetime,
    username varchar(64),
    nickname varchar(128),
    secret varchar(128),
    icon varchar(128),
    remark text,
    level integer,
    status integer,
    email varchar(64),
    login_ip varchar(20) default '',
    primary key(id)
) default charset=utf8;
```

## 3.书页菜单表
id | pid | create_time | update_time | name | sum | code | pic | remark | day | status 
---|---|---|---|---|---|---|---|---|---|--
自增主键(菜单ID) | 父菜单ID |创建时间|更新时间|菜单名称|摘要|访问码|刊页图片|备注状态|发布日期|状态
serial | bigint(20) | timestamp | timestamp | varchar(64)|varchar(128)|varchar(8)|varchar(128)|text|varchar(32)|smallint

```
--mysql
create table tb_book_menu(
    id bigint(20) not null auto_increment,
    pid bigint(20),
    create_time datetime,
    update_time datetime,
    name varchar(64),
    sum varchar(128),
    code varchar(8),
    pic varchar(128),
    remark text,
    day varchar(32),
    status integer,
    author_id integer,
    primary key(id)
) default charset=utf8;
--postgres
create table tb_book_menu(
    id serial primary key,
    pid integer default 0,
    create_time timestamp default CURRENT_TIMESTAMP ,
    update_time timestamp default CURRENT_TIMESTAMP,
    name varchar(64) not null,
    sum varchar(128) default '',
    code varchar(8) default '',
    pic varchar(128) default '',
    remark text,
    day varchar(32) default '',
    status integer default 0,
    author_id integer default 0
);
```

## 4.资源表

```
--mysql
create table tb_resource(
    id bigint(20) not null auto_increment,
    res_type varchar(64) not null default 'pic',
    uri varchar(128) not null default '',
    res_val varchar(64) not null default '',
    create_time datetime default CURRENT_TIMESTAMP ,
    expire_time datetime default CURRENT_TIMESTAMP,
    primary key(id)
) default charset=utf8;
--postgres
create table tb_resource(
  id serial primary key,
  res_type varchar(64) not null default 'pic',
  uri varchar(128) not null default '',
  res_val varchar(64) not null default '',
  create_time timestamp default CURRENT_TIMESTAMP ,
  expire_time timestamp default CURRENT_TIMESTAMP
);
```

## 5.卡片资源表
```
--mysql
drop table tb_card_res;
create table tb_card_res(
    id bigint(20) not null auto_increment,
    res_type varchar(64) not null default 'word',
    data text,
    word varchar(10) not null default '',
    pic varchar(128) not null default '',
    sound varchar(128) not null default '',
    pinyin varchar(64) not null default '',
    scope varchar(64) not null default '',
    gp varchar(32) not null default '',
    create_time datetime default CURRENT_TIMESTAMP ,
    expire_time datetime default CURRENT_TIMESTAMP,
    primary key(id)
) default charset=utf8;
--postgres
drop table tb_card_res;
create table tb_card_res(
    id serial primary key,
    res_type varchar(64) not null default 'word',
    data text,
    word varchar(10) not null default '',
    pic varchar(128) not null default '',
    sound varchar(128) not null default '',
    pinyin varchar(64) not null default '',
    scope varchar(64) not null default '',
    gp varchar(32) not null default '',
    create_time timestamp default CURRENT_TIMESTAMP ,
    expire_time timestamp default CURRENT_TIMESTAMP
);
```

## 6.用户资源表
```
drop table tb_user_card_res;
create table tb_user_card_res(
    id bigint(20) not null auto_increment,
    userid bigint(20) not null ,
    res_id bigint(20) not null ,
    word varchar(10) not null default '',
    scope varchar(64) not null default '',
    gp varchar(32) not null default '',
    create_time datetime default CURRENT_TIMESTAMP ,
    update_time datetime default CURRENT_TIMESTAMP,
    status int not null default 0 comment '0生字,1已学会',
    remark text,
    primary key(id)
) default charset=utf8;

create table tb_user_card_res(
    id serial primary key,
    userid integer not null ,
    res_id integer not null ,
    word varchar(10) not null default '',
    scope varchar(64) not null default '',
    gp varchar(32) not null default '',
    create_time timestamp default CURRENT_TIMESTAMP ,
    update_time timestamp default CURRENT_TIMESTAMP,
    status int not null default 0 comment '0生字,1已学会',
    remark text,
    primary key(id)
)
```