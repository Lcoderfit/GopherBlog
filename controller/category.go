package controller

import (
	"GopherBlog/constant"
	"GopherBlog/model"
	"GopherBlog/utils"
	"GopherBlog/utils/validator"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// AddCategory 添加文章分类
func AddCategory(c *gin.Context) {
	var data model.Category
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		return
	}
	// TODO:优化逻辑
	// 防止输入标签为空格
	data.Name = strings.TrimSpace(data.Name)

	msg, err := validator.Validate(data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DataVerificationError), err)
		failWithData(c, constant.DataVerificationError, msg)
		return
	}

	// 检查分类名称是否已经存在
	code, ok := model.IsCategoryExistsByName(data.Name)
	if ok {
		fail(c, code)
		return
	}

	code = model.CreateCategory(&data)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	// TODO:为什么需要返回data
	successWithData(c, data)
}

// GetCategoryInfo 获取分类信息
func GetCategoryInfo(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		fail(c, constant.ParamError)
		return
	}
	data, code := model.GetCategoryInfo(id)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	successWithData(c, data)
}

// GetCategoryList 获取分类列表
func GetCategoryList(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.Query("page_number"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	// 设置每页的数据量上下限
	if pageSize < 10 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}
	if pageNum <= 0 {
		pageNum = 1
	}

	// 获取分类列表
	categories, code := model.GetCategoryList(pageSize, pageNum)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	successWithData(c, categories)
}

// EditCategoryInfo JWT鉴权接口:编辑分类
func EditCategoryInfo(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		fail(c, constant.ParamError)
		return
	}
	var data model.Category
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
		return
	}
	// TODO:优化
	data.Name = strings.TrimSpace(data.Name)
	msg, err := validator.Validate(data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DataVerificationError), msg)
		failWithData(c, constant.DataVerificationError, msg)
		return
	}

	code, ok := model.IsCategoryExistsByName(data.Name)
	if ok {
		fail(c, code)
		return
	}
	code = model.EditCategoryInfo(id, &data)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	success(c)
}

// DeleteCategory JWT鉴权接口:删除分类
func DeleteCategory(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		fail(c, constant.ParamError)
		return
	}

	// 如果分类不存在，则无需删除
	code, ok := model.IsCategoryExistsById(id)
	if !ok {
		fail(c, code)
		return
	}
	code = model.DeleteCategory(id)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	success(c)
}
