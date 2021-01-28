package controller

import (
	"GopherBlog/constant"
	"GopherBlog/model"
	"GopherBlog/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 获取单个文章信息
func GetArticleInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	data, code := model.GetArticleInfo(id)
	if code != constant.SuccessCode {
		// TODO:如果阅读量更新失败，但是文章正常获取到了还是可以返回
		failWithData(c, code, data)
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
	// 根据标题进行搜索
	title := c.Query("title")

	if pageSize < 10 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	if pageNum <= 0 {
		pageNum = 1
	}

	articles, total, code := model.GetArticleList(title, pageSize, pageNum)
	if code != constant.SuccessCode {
		failWithData(c, code, gin.H{
			"data":  articles,
			"total": total,
		})
	}
	// TODO:修改响应模式，支持自动total
	successWithData(c, gin.H{
		"data":  articles,
		"total": total,
	})
}

// 获取同一分类下的文章
func GetArticleListByCategoryId(c *gin.Context) {
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
	// 获取分类id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	if pageSize < 10 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}
	if pageNum < 0 {
		pageNum = 1
	}

	articles, total, code := model.GetArticleListByCategoryId(id, pageSize, pageNum)
	if code != constant.SuccessCode {
		failWithData(c, code, gin.H{
			"data":  articles,
			"total": total,
		})
	}
	successWithData(c, gin.H{
		"data":  articles,
		"total": total,
	})
}

// JWT:新增文章
func AddArticle(c *gin.Context) {
	var data model.Article
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	// TODO:不用判断文章是否已经存在吗？
	code := model.CreateArticle(&data)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	successWithData(c, data)
}

// JWT:编辑文章
func EditArticleInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	var data model.Article
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	code := model.EditArticleInfo(id, &data)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	success(c)
}

// JWT:删除文章
func DeleteArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	var data model.Article
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	code := model.DeleteArticle(id)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	success(c)
}
