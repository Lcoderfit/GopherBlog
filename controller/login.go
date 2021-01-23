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
	err = model.CheckAccount(data.Username, data.Password)
}
