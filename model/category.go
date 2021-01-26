package model

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
)

type Category struct {
	ID   int    `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}

// 判断文章分类是否存在
func IsCategoryExist(name string) (int, bool) {
	var data Category
	err := db.Where("name = ?", name).Take(&data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DatabaseAccessError), err)
		return constant.DatabaseAccessError, true
	}
	if data.ID > 0 {
		utils.Logger.Error(constant.ConvertForLog(constant.CategoryExistError), err)
		return constant.CategoryExistError, true
	}
	return constant.SuccessCode, false
}

// 创建新分类
func CreateCategory(data *Category) int {
	err := db.Create(data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CreateCategoryError), err)
		return constant.CreateCategoryError
	}
	return constant.SuccessCode
}

// 获取文章分类信息
func GetCategoryInfo(id int) (Category, int) {
	var data Category
	err := db.Where("id = ?", id).Take(&data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.GetCategoryInfoError), err)
		return data, constant.GetCategoryInfoError
	}
	return data, constant.SuccessCode
}

// 获取文章分类列表
func GetCategoryList(pageSize, pageNum int) (data []Category, code int) {
	err := db.Find(&data).Limit(pageSize).Offset(pageSize * (pageNum - 1)).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DatabaseAccessError), err)
		return data, constant.DatabaseAccessError
	}
	if len(data) == 0 {
		utils.Logger.Error(constant.ConvertForLog(constant.CategoryNotExist), err)
		return data, constant.CategoryExistError
	}
	return data, constant.SuccessCode
}
