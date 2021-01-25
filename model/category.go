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
func IsCategoryExists(name string) (bool, error) {
	var data Category
	err := db.Where("name = ?", name).Take(&data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DatabaseAccessError), err)
		return false, err
	}
	if data.ID > 0 {
		utils.Logger.Error(constant.ConvertForLog(constant.CategoryExistsError), err)
		return false, err
	}
	return true, nil
}
