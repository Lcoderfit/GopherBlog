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

	msg, err := validator.Validate(data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.Error,
			"message": msg,
		})
		// 当上面的响应报错时，调用c.Abort()可以确保其不会影响下面的请求处理程序
		c.Abort()
	}

	var code int
	if ok := model.IsUserExists(data.Username); ok {
		code = errmsg.UserAlreadyExistsError
	} else {
		code = errmsg.Success
		model.CreateUser(&data)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetCodeMsg(code),
	})
}
