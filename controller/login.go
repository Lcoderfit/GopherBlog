package controller

import (
	"GopherBlog/constant"
	"GopherBlog/model"
	"GopherBlog/utils"
	"GopherBlog/utils/validator"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var data model.User
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
	}

	msg, err := validator.Validate(data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DataVerificationError), err)
		failWithData(c, constant.DataVerificationError, msg)
	}

	// 检查用户名跟密码是否正确
	var code int
	data, code, err = model.CheckAccount(data.Username, data.Password)
	if err != nil {
		// TODO:遇到c.JSON是否不会返回？？
		fail(c, code)
		return
	}

	// 设置用户token，使用户无需重复登录
	//token, code = middleware.SetToken()
	successWithData(c, gin.H{
		"username": data.Username,
		"id":       data.ID,
		//"token": token,
	})
}
