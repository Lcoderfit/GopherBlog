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

// GetCommentListByArticleId 获取同一文章下的所有评论
func GetCommentListByArticleId(c *gin.Context) {
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

	comments, total, code := model.GetCommentListByArticleId(id, pageSize, pageNum)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	successWithData(c, gin.H{
		"comment_list": comments,
		"total":        total,
	})
}

// GetCommentList JWT-获取评论列表
func GetCommentList(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.Query("page_number"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize < 10 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	comments, total, code := model.GetCommentList(pageNum, pageSize)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	successWithData(c, gin.H{
		"comment_list": comments,
		"total":        total,
	})
}

// UpdateCommentStatus JWT鉴权接口:审核通过评论
func UpdateCommentStatus(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		fail(c, constant.ParamError)
		return
	}
	var data model.Comment
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
		return
	}
	msg, err := validator.Validate(data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DataVerificationError), err)
		failWithData(c, constant.DataVerificationError, msg)
		return
	}
	code := model.UpdateCommentStatus(id, &data)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	success(c)
}

//// TakeDownComment JWT鉴权接口:撤销评论
//func TakeDownComment(c *gin.Context) {
//	id, err := MustInt(c.Param, "id")
//	if err != nil {
//		fail(c, constant.ParamError)
//		return
//	}
//	var data model.Comment
//	if err := c.ShouldBindJSON(&data); err != nil {
//		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
//		fail(c, constant.ParamError)
//		return
//	}
//
//	code := model.TakeDownComment(id, &data)
//	if code != constant.SuccessCode {
//		fail(c, code)
//		return
//	}
//	success(c)
//}

// DeleteComment JWT鉴权接口:删除评论
func DeleteComment(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		fail(c, constant.ParamError)
		return
	}

	code := model.DeleteComment(id)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	success(c)
}
