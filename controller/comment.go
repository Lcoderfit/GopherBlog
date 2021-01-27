package controller

import (
	"GopherBlog/constant"
	"GopherBlog/model"
	"GopherBlog/utils"
	"GopherBlog/utils/validator"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 新增评论
func AddComment(c *gin.Context) {
	var data model.Comment
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
	}

	msg, err := validator.Validate(data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DataVerificationError), err)
		failWithData(c, constant.DataVerificationError, msg)
	}

	code := model.CreateComment(&data)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	successWithData(c, data)
}

// 获取评论
func GetCommentInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	comments, code := model.GetCommentInfo(id)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	successWithData(c, comments)
}

// 获取文章评论数量
func GetCommentCount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	count, code := model.GetCommentCount(id)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	successWithData(c, count)
}

// 前端获取评论列表
func GetCommentList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	pageSize, err := strconv.Atoi(c.Param("page_size"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	pageNum, err := strconv.Atoi(c.Param("page_num"))
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

	comments, total, code := model.GetCommentList(id, pageSize, pageNum)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	successWithData(c, gin.H{
		"data":  comments,
		"total": total,
	})
}
