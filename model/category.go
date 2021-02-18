package model

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
)

type Category struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(20);not null" validate:"required" json:"name"`
}

// IsCategoryExistsByName 根据名称判断文章分类是否存在
func IsCategoryExistsByName(name string) (int, bool) {
	var data Category
	// 如果不存在才会返回err(gorm.ErrRecordNotFound), 判断是否存在可以忽略err
	db.Where("name = ?", name).Take(&data)
	if data.ID > 0 {
		utils.Logger.Error(constant.ConvertForLog(constant.CategoryExistsError), err)
		return constant.CategoryExistsError, true
	}
	return constant.SuccessCode, false
}

// IsCategoryExistsById 根据ID判断文章分类是否存在
func IsCategoryExistsById(id int) (int, bool) {
	var data Category
	err := db.Where("id = ?", id).Take(&data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CategoryNotExistsError))
		return constant.CategoryNotExistsError, false
	}
	return constant.SuccessCode, true
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
		utils.Logger.Error(constant.ConvertForLog(constant.CategoryNotExistsError), err)
		return data, constant.CategoryNotExistsError
	}
	return data, constant.SuccessCode
}

// GetCategoryList 获取文章分类列表
func GetCategoryList(pageSize, pageNum int) (data []Category, code int) {
	err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DatabaseAccessError), err)
		return data, constant.DatabaseAccessError
	}
	return data, constant.SuccessCode
}

// EditCategoryInfo 编辑分类信息
func EditCategoryInfo(id int, data *Category) int {
	// update category set name=xxx where id=yyy
	err := db.Model(&Category{}).Where("id = ?", id).Updates(data).Error
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
