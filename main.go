package main

import (
	"GopherBlog/model"
	"GopherBlog/routes"
)

func main() {
	//user := model.User{
	//	Username: "l00",
	//	Password: "fajskdjf",
	//	Role:     4,
	//}
	//msg, err := validator.Validate(&user)
	//if err != nil {
	//	utils.Logger.Error(err)
	//	//utils.Logger.Error(msg)
	//	fmt.Println("hk")
	//}
	//utils.Logger.Info(msg)
	// 初始化数据库
	model.InitDB()
	// 初始化路由
	routes.InitRouter()
}
