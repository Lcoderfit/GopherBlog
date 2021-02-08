// GoBlog
package main

import (
	"fmt"
	"strconv"
)

func main() {
	a, err := strconv.Atoi("1a0")
	fmt.Println(a, err)

	//// 初始化数据库
	//model.InitDB()
	//// 初始化路由
	//routes.InitRouter()
}
