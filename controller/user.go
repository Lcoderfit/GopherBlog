package controller

import (
	"GopherBlog/model"
	"GopherBlog/utils"
	"GopherBlog/utils/errmsg"
	"GopherBlog/utils/validator"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 添加用户
func AddUser(c *gin.Context) {
	// 1.声明一个模型
	// 2.参数校验
	// 3.数据校验
	// 4.调用模型中声明的与数据库交互的函数进行数据交互
	var data model.User
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error("参数输入错误： ", err)
	}

	msg, code := validator.Validate(data)
	if code != errmsg.ERROR {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": msg,
		})
		c.Abort()
	}

}
