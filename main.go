// GoBlog
package main

import (
	"GopherBlog/model"
	"GopherBlog/routes"
)

func main() {
	// 初始化数据库
	model.InitDB()
	// 初始化路由
	routes.InitRouter()
}
