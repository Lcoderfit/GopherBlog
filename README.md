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





