package model

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
)

type Category struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}

// IsCategoryExist 判断文章分类是否存在
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

// CreateCategory 创建新分类
func CreateCategory(data *Category) int {
	err := db.Create(data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CreateCategoryError), err)
		return constant.CreateCategoryError
	}
	return constant.SuccessCode
}

// GetCategoryInfo 获取文章分类信息
func GetCategoryInfo(id int) (Category, int) {
	var data Category
	err := db.Where("id = ?", id).Take(&data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.GetCategoryInfoError), err)
		return data, constant.GetCategoryInfoError
	}
	return data, constant.SuccessCode
}

// GetCategoryList 获取文章分类列表
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

// EditCategoryInfo 编辑分类信息
func EditCategoryInfo(id int, data *Category) int {
	maps := make(map[string]interface{})
	maps["name"] = data.Name
	// update category set name=xxx where id=yyy
	err := db.Model(&Category{}).Select("id = ?", id).Updates(maps).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.EditCategoryInfoError), err)
		return constant.EditCategoryInfoError
	}
	return constant.SuccessCode
}

// DeleteCategory 删除分类信息
func DeleteCategory(id int) int {
	// TODO:Update用Model，Delete一般用Where？？
	err := db.Where("id = ?", id).Delete(&Category{}).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DeleteCategoryError), err)
		return constant.DeleteCategoryError
	}
	return constant.SuccessCode
}
