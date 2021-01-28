package controller

import (
	"GopherBlog/constant"
	"GopherBlog/model"
	"GopherBlog/utils"
	"GopherBlog/utils/validator"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 添加文章分类
func AddCategory(c *gin.Context) {
	var data model.Category
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
	}

	msg, err := validator.Validate(data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DataVerificationError), err)
		failWithData(c, constant.DataVerificationError, msg)
	}

	// 检查分类名称是否已经存在
	code, ok := model.IsCategoryExist(data.Name)
	if ok {
		fail(c, code)
	}

	code = model.CreateCategory(&data)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	// TODO:为什么需要返回data
	successWithData(c, data)
}

// 获取分类信息
func GetCategoryInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	data, code := model.GetCategoryInfo(id)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	successWithData(c, data)
}

// 获取分类列表
func GetCategoryList(c *gin.Context) {
	pageNum, err := strconv.Atoi(c.Param("pageNum"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	pageSize, err := strconv.Atoi(c.Param("pageSize"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	// 设置每页的数据量上下限
	if pageSize < 10 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}
	if pageNum < 0 {
		pageNum = 1
	}

	// 获取分类列表
	categories, code := model.GetCategoryList(pageSize, pageNum)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	successWithData(c, categories)
}

// JWT:编辑分类
func EditCategoryInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	var data model.Category
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	code, ok := model.IsCategoryExist(data.Name)
	if !ok {
		fail(c, code)
	}
	code = model.EditCategoryInfo(id, &data)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	success(c)
}

// JWT:删除分类
func DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	//var data model.Category
	//if err := c.ShouldBindJSON(&data); err != nil {
	//	utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
	//	fail(c, constant.ParamError)
	//}

	// TODO:为什么删除接口不需要判断分类是否存在
	code := model.DeleteCategory(id)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	success(c)
}
