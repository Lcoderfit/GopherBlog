# GopherBlog
基于Gin+Vue+MySQL开发的博客系统

### 普通接口

### 鉴权接口

### 功能介绍

### 项目部署

### 技术栈

### 项目简介

1.定义接口标准
2.返回Error的场景
3.打印log的场景和等级
4.状态码分级
5.常量模块化
6.常量如何用于日志？？
7.log中的中文内容是否需要在constant中设置特定的code和msg
8.针对不同场景的响应方法：success，fail
9.当返回err不合适时，可以用返回err code，然后调用方判断返回的是否为successCode来返回响应

Q:
1.log存在重复
2.log和err系统不统一
3.统一修改成以code返回的形式？？/
4.记录调试时会遇到的错误
5.利用pporf分析性能
6.压测工具
7.什么时候返回error，什么时候返回panic
8.json tag用下划线命名法，gorm标签用小驼峰，其他用大小驼峰，router用什么命名法？？
9.GraphQL替换RESTful， grpc等等
10.controller的参数获取需要修改成全部按照ShouldJSONBind接收
11.需要有两个系统,一个用于用户使用的系统,另一个用来给管理员使用
12.复用代码，简化流程
13.拦截请求，将参数校验和err处理放在中间件中进行, 使得请求处理函数可以直接使用数据，
当在中间件中参数校验失败时，直接调用c.JSON返回（可以封装成fail函数）
// 14.如果c.Next()和c.Abort()有多层嵌套，则执行顺序是怎样的？？
15.设置了JSONFormatter输出还是text格式
16.对项目进行热部署，修改代码可以实时热更新
17.c.ShouldBindJSON(&data), 如果data有多个字段，但是请求中只包含其中部分字段，则会设置那一部分字段的值，其他字段会取默认值；
请求的参数与data中字段名的对应关系对大小写不敏感，例如请求中有个字段“id”， 但是data中的字段为ID，也是可以匹配的；
如果请求中同时存在"id"和“ID”，id在ID上面，则会取更下面的（即“ID”）字段值传给data中的字段
但是有一个问题，内嵌了gorm.Model的结构体，在输出结构体实例的时候ID，CreateAt，UpdateAt，DeleteAt都是大写的

//18.通过goland的configuration的environment选项可以绕过管理员权限修改环境变量
//19.post/put方法一般会传递json数据，而get和delete方法传递的一般是路由参数（Param）或者url参数（Query）;
路由参数和url参数不可避免的需要手动对数据进行校验
20.后端针对get请求的url或者路由参数先采取直接返回策略，等前后端联调时再进行架构考虑
21.rest设计规范，注意：由HTTP动词+URI名词组成，URI中不能包含动词；
URI中的名词均用复数表示；
使用连字符‘-’提高url可读性，而不是下划线;
末尾不以‘/’结尾
get /v1/users
get /v1/users/id
get /v1/users
post /v1/users

22.id一般作为路由参数，page_num,page_size，limit这种一般作为详情参数


23.First/Take/Last方法如果没有找到都会返回gorm.ErrRecordNotFound错误
Find方法不会返回ErrRecordNotFound错误，但是使用limit和offeset时候需要
先使用limit和offeset再使用find，否则得到的结果集将一直是全集
正确用法：db.Limit().Offeset().Find(&data)
错误：db.Find(&data).Limit().Offeset()

24.获取分类列表接口不传入分页参数时，是否返回默认排序的第一页数据

25.Preload中的字段大小写不敏感，例如结构体为Category，则db.Preload("category")...也是有效的
不过一般都使用db.Preload("Category")....

model中定义模型时，gorm标签的gorm:"foreignKey:Cid"也是大小写不敏感的

26.page_number和page_size参数需要对err进行忽略，因为如果传入参数错误需要保证返回默认的数据列表页

27.Joins查询几种示例
正确：
db.Model(&Comment{}).Where......Joins("left join user on .....").Scan(&comment)
db.Model(&Comment{}).Where......Joins("left join user on .....").Find(&comment)
db.Where......Joins("left join user on .....").Find(&comment)

错误：
db.Where......Joins("left join user on .....").Find(&comment)

28.total返回的是所有的数据量，而不是单单一个页面的数据量

29./admin/check_token修改成Rest接口 ----> /admin/token-check

30.gorm.Model如果定义了DeleteAt字段，则具有软删除功能；即并不是真正的删除，而是将DeleteAt字段更新为删除语句执行的时间
可以通过db.UnScoped().Where()查询软删除的数据，要永久删除可以使用db.UnScoped().Delete()

31.加入pv，uv等redis统计功能；Nginx负载均衡

32.自建issue

33.用户更新接口更新用户信息时没有对用户名进行判重
34.如果忘记密码怎么办，现有的修改密码接口需要进行JWT校验，万一JWT失效了呢

错误：
一.数据库连接失败
1.config.ini文件中的字段需要与定义的结构体字段名字相同(大小写也必须一致)
2.dsn格式：username:passowrd@tcp(host:port)/dbName?charset=utf8&parseTime=true&loc=Local

二、os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)报错：The system cannot find the path specified
如果filePath参数传入的是带目录的路径，例如:log/output.log，则需要先创建log目录，OpenFile只会创建最后的output.log,
不会创建父目录，所以如果不手动创建则会报错

三、内嵌了gorm.Model的结构体，通过响应返回实例时候ID，CreateAt，UpdateAt，DeleteAt都是大写的
需要修改gorm中model.go的源码，添加json tag

四、gorm.Model中的CreateAt，UpdateAt，DeleteAt时区错误
gorm.Model中的CreateAt和UpdateAt为time.Time类型，而DeleteAt为sql.NullTime类型（结构体）
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}
bak.sql的测试数据中id=1的用户创建时间、更新时间、删除时间均为null，所以从数据库中取出来之后，删除时间为null，
而创建时间和更新时间为time.Time的零值，即：0001-01-01 00:00:00 +0000 UTC（apifox取出来为：0001-01-01T00:00:00Z）
Z就是世界协调时间，跟UTC是一样的

五、windows终端下logrus打印日志没有颜色
import "github.com/shiena/ansicolor"
Logger.SetFormatter(&logrus.TextFormatter{
    ForceColors:     true,  // 这个要设置为true
})
// fix:解决logrus在windows终端下输出无颜色区别的问题
Logger.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))

六、Only one usage of each socket address (protocol/network address/port)
每个套接字地址只有一种用法，也就是说端口被其他程序占用了

七、密码加密导致密码长度超过数据库限制
将数据库字段长度调大，sql：

八、data返回格式不统一：单个数据返回字典，多个数据返回列表？？

九、db.Find(&data).Limit(pageSize).Offset((pageNum - 1) * pageSize).Error
无论怎么查data都会保存数据库里的数据

改成：db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&data).Error就没有问题

十、err: unsupported data type: 0xc0003622d0: Table not set, please set it like: db.Model(&user) or db.Table("users")
在查询时没有指明需要查的是哪个表
错误：
db.Where("article_id = ?", articleId).Count(&total).Error
正确：
db.Model(Comment{}).Where("article_id = ?", articleId).Count(&total).Error

十一、err: sql: expected 11 destination arguments in Scan, not 1
错误：用Take来获取int64类型的值，而Take只能接受Comment类型结构体指针
db.Model(&Comment{}).Where("article_id = ? and status = ?", articleId, 1).Take(&count).Error
正确:
注意：count需要是int64类型
db.Model(&Comment{}).Where("article_id = ? and status = ?", articleId, 1).Count(&count).Error

十二、获取用户列表，只需要获取用户的username, role和id三个字段，但是查出来的每一个用户都包含所有字段.
可以通过设置users为[]map[string]interface{}类型，然后：
db.Model(&User{}).Select("id, role, username")......Find(&users)

十三、 err: invalid character '-' in numeric literal
通常是客户端发起的请求参数格式与服务端不一致导致的，例如客户端发送表单类型而服务端接受json类型

十四、更新用户时，不加validator验证输入非法内容会导致问题，加上validator验证又必须所有字段都得传入

除非必要字段，需要更新的字段都需要传入?????
