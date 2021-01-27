package model

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
	"gorm.io/gorm"
)

// gorm.Model中设置了ID为主键， 如果没有使用这个，需要显示设置一个字段为"primaryKey"
// TODO:需要设置成belongs to或者has many关系
type Comment struct {
	gorm.Model
	UserId    int    `json:"user_id"`
	ArticleId int    `json:"article_id"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	Content   string `gorm:"type:varchar(500);not null" json:"content"` // 评论不能为空
	Status    int8   `gorm:"tinyint;not null;default:2" json:"status"`  // 状态默认为2
}

// 新增评论, 评论内容可以重复，需要不需要像User的model一样判断是否存在
func CreateComment(data *Comment) int {
	err := db.Create(data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CreateCommentError), err)
		return constant.CreateCommentError
	}
	return constant.SuccessCode
}

// 获取评论
func GetCommentInfo(id int) (comment Comment, code int) {
	err := db.Where("id = ?", id).Take(&comment).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.GetCommentInfoError), err)
		return comment, constant.GetCommentInfoError
	}
	return comment, constant.SuccessCode
}

// 获取评论列表
func GetCommentCount(articleId int) (count, code int) {
	// TODO:status为1表示是正常用户，为2表示什么呢？
	err := db.Where("article_id = ? and status = ?", articleId, 1).Take(&count).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.GetCommentCountError), err)
		return 0, constant.GetCommentCountError
	}
	return count, constant.SuccessCode
}

// 前端获取评论列表
func GetCommentList(articleId, pageSize, pageNum int) (comments []Comment, total int64, code int) {
	err := db.Where("article_id = ?", articleId).Count(&total).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CountCommentError), err)
		return comments, 0, constant.CountCommentError
	}
	// 1.Joins也可以通过“？”传参
	// 2.当有多个Joins链式连接时，可以把前一个joins得到的结果作为与后一个Joins进行连接的数据集
	err = db.Where("article_id = ?", articleId).Limit(pageSize).Offset(pageSize * (pageNum - 1)).Order(
		"create_at desc",
	).Select(
		"comment.id, comment.user_id, user.username, comment.article_id, " +
			"article.title, content, status, create_at, delete_at",
	).Joins(
		"left join user on comment.user_id=user.id",
	).Joins("left join article on comment.article_id=article.id").Find(&comments).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.GetCommentListError), err)
		return comments, 0, constant.GetCommentListError
	}
	return comments, total, constant.SuccessCode
}
