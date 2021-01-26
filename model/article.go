package model

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
	"gorm.io/gorm"
)

// 文章结构体
type Article struct {
	gorm.Model
	Category     Category `gorm:"foreignKey:Cid"`
	Cid          int      `gorm:"type:int;not null" json:"cid"`
	Title        string   `gorm:"type:varchar(100)" json:"title"`
	Content      string   `gorm:"type:longtext" json:"content"`
	Img          string   `gorm:"type:varchar(100)" json:"img"`
	Desc         string   `gorm:"type:varchar(200)" json:"desc"`
	CommentCount int      `gorm:"type:int;not null;default:0" json:"comment_count"`
	ReadCount    int      `gorm:"type:int;not null;default:0" json:"read_count"`
}

// 获取文章信息
func GetArticleInfo(id int) (Article, int) {
	var data Article
	err := db.Where("id = ?", id).Take(&data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.GetArticleInfoError), err)
		return data, constant.GetArticleInfoError
	}
	return data, constant.SuccessCode
}

// 获取文章列表
func GetArticleList(pageSize, pageNum int) (articles []Article, code int) {
	err := db.
}