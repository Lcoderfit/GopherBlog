package controller

import (
	"GopherBlog/constant"
	"GopherBlog/model"
	"GopherBlog/utils"
	"GopherBlog/utils/validator"
	"github.com/gin-gonic/gin"
	"strconv"
)

// AddComment 新增评论
func AddComment(c *gin.Context) {
	var data model.Comment
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		return
	}

	msg, err := validator.Validate(data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DataVerificationError), err)
		failWithData(c, constant.DataVerificationError, msg)
		return
	}

	code := model.CreateComment(&data)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	successWithData(c, data)
}

// GetCommentInfo 获取评论
func GetCommentInfo(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		fail(c, constant.ParamError)
		return
	}
	comments, code := model.GetCommentInfo(id)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	successWithData(c, comments)
}

// GetCommentCount 获取文章评论数量
func GetCommentCount(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		fail(c, constant.ParamError)
		return
	}
	count, code := model.GetCommentCount(id)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	successWithData(c, count)
}

// GetCommentList 前端获取评论列表
func GetCommentList(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		fail(c, constant.ParamError)
		return
	}
	pageNum, _ := strconv.Atoi(c.Query("page_number"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

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
		return
	}
	successWithData(c, gin.H{
		"data":  comments,
		"total": total,
	})
}

// ApproveComment JWT鉴权接口:审核评论
func ApproveComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	var data model.Comment
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	code := model.ApproveComment(id, &data)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	success(c)
}

// TakeDownComment JWT鉴权接口:撤销评论
func TakeDownComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	var data model.Comment
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	code := model.TakeDownComment(id, &data)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	success(c)
}

// DeleteComment JWT鉴权接口:删除评论
func DeleteComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	var data model.Comment
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	code := model.DeleteComment(id)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	success(c)
}
