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

错误：
一.数据库连接失败
1.config.ini文件中的字段需要与定义的结构体字段名字相同(大小写也必须一致)
2.dsn格式：username:passowrd@tcp(host:port)/dbName?charset=utf8&parseTime=true&loc=Local

二、os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)报错：The system cannot find the path specified
如果filePath参数传入的是带目录的路径，例如:log/output.log，则需要先创建log目录，OpenFile只会创建最后的output.log,
不会创建父目录，所以如果不手动创建则会报错



