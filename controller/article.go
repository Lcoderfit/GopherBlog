package controller

import (
	"GopherBlog/constant"
	"GopherBlog/model"
	"GopherBlog/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 获取文章信息
func GetArticleInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	data, code := model.GetArticleInfo(id)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	successWithData(c, data)
}

// 获取文章列表
func GetArticleList(c *gin.Context) {
	pageSize, err := strconv.Atoi(c.Param("pageSize"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	pageNum, err := strconv.Atoi(c.Param("pageNum"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	if pageSize < 10 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	if pageNum <= 0 {
		pageNum = 1
	}

	articles, code := model.GetArticleList(pageSize, pageNum)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	successWithData(c, articles)
}
