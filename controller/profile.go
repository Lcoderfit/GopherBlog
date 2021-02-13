package controller

import (
	"GopherBlog/constant"
	"GopherBlog/model"
	"GopherBlog/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetProfileInfo 获取个人简介
func GetProfileInfo(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		return
	}
	data, code := model.GetProfileInfo(id)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	successWithData(c, data)
}

// UpdateProfileInfo JWT鉴权接口:更新个人信息
func UpdateProfileInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	var data model.Profile
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	code := model.UpdateProfileInfo(id, &data)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	success(c)
}
