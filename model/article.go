package model

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
	"gorm.io/gorm"
)

// Article 文章结构体
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

// GetArticleInfo 获取单个文章信息
func GetArticleInfo(id int) (article Article, code int) {
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

// GetArticleList 获取文章列表
func GetArticleList(title string, pageSize, pageNum int) (articles []Article, total int64, code int) {
	// TODO:1.需要增加文章按时间排序、按阅读量排序的功能(或集成redis)
	// TODO:2.需要集成专门的搜索引擎用于搜索，而不是使用MySQL的模糊搜索
	// 默认按照时间降序排序
	if title == "" {
		err := db.Limit(pageSize).Offset(pageSize * (pageNum - 1)).Order(
			"created_at desc",
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
			"created_at desc",
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

// GetArticleListByCategoryId 通过分类Id获取文章列表
func GetArticleListByCategoryId(id, pageSize, pageNum int) (articles []Article, total int64, code int) {
	// TODO:Preload先后顺序是否对结果有影响
	err := db.Preload("Category").Where("cid = ?", id).Limit(pageSize).Offset(
		pageSize * (pageNum - 1),
	).Order("created_at desc").Find(&articles).Error
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

// CreateArticle 创建文章
func CreateArticle(data *Article) int {
	err := db.Create(&data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CreateArticleError), err)
		return constant.CreateArticleError
	}
	return constant.SuccessCode
}

// EditArticleInfo 编辑文章
func EditArticleInfo(id int, data *Article) int {
	// updates操作最好还是使用map作为参数传入
	maps := map[string]interface{}{
		"Title":   data.Title,
		"desc":    data.Desc,
		"cid":     data.Cid,
		"content": data.Content,
		"img":     data.Img,
	}
	err := db.Model(&Article{}).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.EditArticleInfoError), err)
		return constant.EditArticleInfoError
	}
	return constant.SuccessCode
}

// DeleteArticle 删除文章
func DeleteArticle(id int) int {
	err := db.Where("id = ?", id).Delete(&Article{}).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DeleteArticleError), err)
		return constant.DeleteArticleError
	}
	return constant.SuccessCode
}
