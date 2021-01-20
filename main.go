package main

import (
	"GopherBlog/model"
	"GopherBlog/utils"
	"GopherBlog/utils/validator"
	"fmt"
)

func main() {
	user := model.User{
		Username: "l00",
		Password: "fajskdjf",
		Role:     4,
	}
	msg, err := validator.Validate(&user)
	if err != nil {
		utils.Logger.Error(err)
		//utils.Logger.Error(msg)
		fmt.Println("hk")
	}
	utils.Logger.Info(msg)

}
