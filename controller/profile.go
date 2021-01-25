package controller

import (
	"GopherBlog/constant"
	"GopherBlog/model"
	"GopherBlog/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 获取个人简介
func GetProfile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
	}
	data, code := model.GetProfileInfo(id)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	successWithData(c, data)
}

//
