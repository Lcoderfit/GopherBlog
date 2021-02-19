package controller

import (
	"GopherBlog/constant"
	"GopherBlog/model"
	"GopherBlog/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetArticleInfo 获取单个文章信息
func GetArticleInfo(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		fail(c, constant.ParamError)
		return
	}

	data, code := model.GetArticleInfo(id)
	if code != constant.SuccessCode {
		// TODO:如果阅读量更新失败，但是文章正常获取到了还是可以返回
		fail(c, code)
		return
	}
	successWithData(c, data)
}

// GetArticleList 获取文章列表
func GetArticleList(c *gin.Context) {
	// 如果传入的参数错误，那会页码和每页的数据个数取默认值，后面仍然能正常返回数据
	// 如果这里对错误进行处理并返回，则文章列表显示为空
	pageNum, _ := strconv.Atoi(c.Query("page_number"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
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
			"article_list": articles,
			"total":        total,
		})
		return
	}
	// TODO:修改响应模式，支持自动total
	successWithData(c, gin.H{
		"article_list": articles,
		"total":        total,
	})
}

// GetArticleListByCategoryId 获取同一分类下的文章
func GetArticleListByCategoryId(c *gin.Context) {
	// 获取分类id
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
	if pageNum < 0 {
		pageNum = 1
	}

	articles, total, code := model.GetArticleListByCategoryId(id, pageSize, pageNum)
	if code != constant.SuccessCode {
		failWithData(c, code, gin.H{
			"category_list": articles,
			"total":         total,
		})
		return
	}
	successWithData(c, gin.H{
		"category_list": articles,
		"total":         total,
	})
}

// AddArticle JWT鉴权接口：新增文章
func AddArticle(c *gin.Context) {
	var data model.Article
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
		return
	}
	// TODO:参数校验，定义新模型：ArticleAddition
	// TODO:一开始创建如果不添加分类，后期如何为文章添加分类，通过更新文章接口？？
	code := model.CreateArticle(&data)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	successWithData(c, data)
}

// EditArticleInfo JWT鉴权接口:编辑文章
func EditArticleInfo(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		fail(c, constant.ParamError)
		return
	}
	var data model.Article
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
		return
	}
	// TODO:是否要加参数校验
	code := model.EditArticleInfo(id, &data)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	success(c)
}

// DeleteArticle JWT鉴权接口:删除文章
func DeleteArticle(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		fail(c, constant.ParamError)
		return
	}

	code := model.DeleteArticle(id)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	success(c)
}
