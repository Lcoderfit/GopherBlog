package controller

import (
	"GopherBlog/constant"
	"GopherBlog/middleware"
	"GopherBlog/model"
	"GopherBlog/utils"
	"GopherBlog/utils/validator"
	"github.com/gin-gonic/gin"
)

// Login 登录接口
func Login(c *gin.Context) {
	var data model.User
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		return
	}

	msg, err := validator.Validate(data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DataVerificationError), err)
		failWithData(c, constant.DataVerificationError, msg)
		return
	}

	// 检查用户名跟密码是否正确
	var code int
	data, code = model.CheckAccount(data.Username, data.Password)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}

	// 设置用户token，使用户无需重复登录
	token, code := middleware.SetToken(data.Username)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}

	// 返回token到客户端
	successWithData(c, gin.H{
		"username": data.Username,
		"id":       data.ID,
		"token":    token,
	})
}

// UpToken 用于存储更新后的token
type UpToken struct {
	Token string `json:"token"`
}

// CheckToken 验证token
func CheckToken(c *gin.Context) {
	var data UpToken
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
	}
	msg, err := validator.Validate(data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DataVerificationError), err)
		failWithData(c, constant.DataVerificationError, msg)
	}

	_, code, err := middleware.CheckToken(data.Token)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(code), err)
		fail(c, code)
	}
	success(c)
}
