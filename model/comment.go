package model

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
	"gorm.io/gorm"
)

// Comment gorm.Model中设置了ID为主键， 如果没有使用这个，需要显示设置一个字段为"primaryKey"
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

// CreateComment 新增评论, 评论内容可以重复，需要不需要像User的model一样判断是否存在
func CreateComment(data *Comment) int {
	err := db.Create(data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CreateCommentError), err)
		return constant.CreateCommentError
	}
	return constant.SuccessCode
}

// GetCommentInfo 获取评论
func GetCommentInfo(id int) (comment Comment, code int) {
	err := db.Where("id = ?", id).Take(&comment).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CommentNotExistError), err)
		return comment, constant.CommentNotExistError
	}
	return comment, constant.SuccessCode
}

// GetCommentCount 获取评论列表
func GetCommentCount(articleId int) (count int64, code int) {
	// TODO:status为1表示是正常用户，为2表示什么呢？
	err := db.Model(&Comment{}).Where("article_id = ? and status = ?", articleId, 1).Count(&count).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.GetCommentCountError), err)
		return 0, constant.GetCommentCountError
	}
	return count, constant.SuccessCode
}

// GetCommentList 前端获取评论列表
func GetCommentList(articleId, pageSize, pageNum int) (comments []Comment, total int64, code int) {
	err := db.Model(Comment{}).Where("article_id = ?", articleId).Count(&total).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CountCommentError), err)
		return comments, 0, constant.CountCommentError
	}
	// 1.Joins也可以通过“？”传参
	// 2.当有多个Joins链式连接时，可以把前一个joins得到的结果作为与后一个Joins进行连接的数据集
	// 可以直接在最后面使用Find，前面就可以省略Model(&Comment{}), 或者前面用db.Model(&Comment{})后面用Scan(&comment)
	// 获取同一文章下的所有用户发表的评论
	err = db.Model(&Comment{}).Where("article_id = ?", articleId).Limit(pageSize).Offset(
		pageSize * (pageNum - 1),
	).Order(
		"created_at desc",
	).Select(
		"comment.id, comment.user_id, user.username, comment.article_id, " +
			"article.title, comment.content, comment.status, comment.created_at, comment.deleted_at",
	).Joins(
		"left join user on comment.user_id=user.id",
	).Joins("left join article on comment.article_id=article.id").Scan(&comments).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.GetCommentListError), err)
		return comments, 0, constant.GetCommentListError
	}
	return comments, total, constant.SuccessCode
}

// ApproveComment 通过评论
func ApproveComment(id int, data *Comment) int {
	maps := map[string]interface{}{
		"status": data.Status,
	}
	var comment Comment
	// 更新评论状态
	err := db.Model(&Comment{}).Where("id = ?", id).Updates(maps).Take(&comment).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ApproveCommentError), err)
		return constant.ApproveCommentError
	}
	// 将文章的评论数增加1
	err = db.Model(&Article{}).Where("id = ?", comment.ArticleId).UpdateColumn(
		"comment_count", gorm.Expr("comment_count + ?", 1),
	).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.AddCommentCountError), err)
		return constant.AddCommentCountError
	}
	return constant.SuccessCode
}

// TakeDownComment 撤销审核
func TakeDownComment(id int, data *Comment) int {
	maps := map[string]interface{}{
		"status": data.Status,
	}
	var comment Comment
	// 更新status
	err := db.Model(&Comment{}).Where("id = ?", id).Updates(maps).Take(&comment).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.TakeDownCommentError), err)
		return constant.TakeDownCommentError
	}
	// 评论数-1
	err = db.Model(&Article{}).Where("id = ?", comment.ArticleId).UpdateColumn(
		"comment_count", gorm.Expr("comment_count + ?", 1),
	).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ReduceCommentCountError), err)
		return constant.ReduceCommentCountError
	}
	return constant.SuccessCode
}

// DeleteComment 删除评论
func DeleteComment(id int) int {
	err := db.Where("id = ?", id).Delete(&Comment{}).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DeleteCommentError), err)
		return constant.DeleteCommentError
	}
	return constant.SuccessCode
}
