package model

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
	"gorm.io/gorm"
)

// 文章结构体
// TODO：foreignKey标签后面的字段名是否大小写敏感？？
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

// 获取单个文章信息
func GetArticleInfo(id int) (article Article, code int) {
	// TODO:Preload中的字段名是否大小写敏感
	// Find不会返回ErrRecordNotFound， 但是First和Take会
	err := db.Preload("Category").Where("id = ?", id).Take(&article).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ArticleNotExistError), err)
		return article, constant.ArticleNotExistError
	}
	// 使用sql表达式更新文章阅读量, 字段名用数据库中的字段名或模型中的字段名均可
	// TODO：阅读量需要设置同一个客户端在一定时间段内只增加一次，不然会存在刷阅读量的行为
	err = db.Model(&article).Where("id = ?", id).UpdateColumn(
		"read_count", gorm.Expr("read_count + ?", 1),
	).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.UpdateReadCountError), err)
		return article, constant.UpdateReadCountError
	}
	return article, constant.SuccessCode
}

// 获取文章列表
func GetArticleList(title string, pageSize, pageNum int) (articles []Article, total int64, code int) {
	// TODO:1.需要增加文章按时间排序、按阅读量排序的功能(或集成redis)
	// TODO:2.需要集成专门的搜索引擎用于搜索，而不是使用MySQL的模糊搜索
	// 默认按照时间降序排序
	if title == "" {
		err := db.Limit(pageSize).Offset(pageSize * (pageNum - 1)).Order(
			"create_at desc",
		).Preload("Category").Find(&articles).Error
		if err != nil {
			utils.Logger.Error(constant.ConvertForLog(constant.GetArticleListInfoError), err)
			return articles, 0, constant.GetArticleListInfoError
		}
		// TODO:其他model中如果有要计算total的需要修改逻辑
		// 计算所有文章的总数，而不是一页的数量
		err = db.Model(&articles).Count(&total).Error
		if err != nil {
			utils.Logger.Error(constant.ConvertForLog(constant.CountArticleListError), err)
			return articles, 0, constant.CountArticleListError
		}
	} else {
		err := db.Limit(pageSize).Offset(pageSize*(pageNum-1)).Order(
			"create_at desc",
		).Preload("Category").Where("title like ?", "%"+title+"%").Find(&articles).Error
		if err != nil {
			utils.Logger.Error(constant.ConvertForLog(constant.GetArticleListInfoError), err)
			return articles, 0, constant.GetArticleListInfoError
		}
		err = db.Model(&articles).Where("title like ?", "%"+title+"%").Error
		if err != nil {
			utils.Logger.Error(constant.ConvertForLog(constant.CountArticleListError), err)
			return articles, 0, constant.CountArticleListError
		}
	}
	return articles, total, constant.SuccessCode
}

// 通过分类Id获取文章列表
func GetArticleListByCategoryId(id, pageSize, pageNum int) (articles []Article, total int64, code int) {
	// TODO:Preload先后顺序是否对结果有影响
	err := db.Preload("Category").Where("cid = ?", id).Limit(pageSize).Offset(
		pageSize * (pageNum - 1),
	).Order("create_at desc").Find(&articles).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.GetArticleListByCategoryIdError), err)
		return articles, 0, constant.GetArticleListByCategoryIdError
	}
	err = db.Model(&articles).Where("cid = ?", id).Count(&total).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CountArticleListError), err)
		return articles, 0, constant.CountArticleListError
	}
	return articles, total, constant.SuccessCode
}
