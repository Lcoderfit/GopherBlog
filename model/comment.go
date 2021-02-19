package model

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
	"gorm.io/gorm"
)

// Comment gorm.Model中设置了ID为主键， 如果没有使用这个，需要显示设置一个字段为"primaryKey"
// TODO:需要设置成belongs to或者has many关系
// TODO：comment has many comment
type Comment struct {
	gorm.Model
	UserId    int    `json:"user_id"`
	ArticleId int    `json:"article_id"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	Content   string `gorm:"type:varchar(500);not null" json:"content"`                       // 评论不能为空
	Status    int8   `gorm:"tinyint;not null;default:2" validate:"oneof=0 1 2" json:"status"` // 状态默认为2，1表示审核通过，2表示待审核，
	// 0表示审核不通过或被撤销
}

// CreateComment 新增评论
// TODO:前端如果传入的用户ID和Username不一致，则数据库中数据就会错乱
// TODO:添加事务，如果两个sql操作，其中一个成功，另一个不成功，则会导致数据不一致
func CreateComment(data *Comment) int {
	err := db.Create(data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CreateCommentError), err)
		return constant.CreateCommentError
	}

	// 更新文章评论数； TODO:评论和文章是否均需要审核？？？
	//err = db.Model(&Article{}).Where("id = ?", data.ArticleId).UpdateColumn(
	//	"comment_count", gorm.Expr("comment_count + ?", 1),
	//).Error
	//if err != nil {
	//	utils.Logger.Error(constant.ConvertForLog(constant.UpdateCommentCountError), err)
	//	return constant.UpdateCommentCountError
	//}
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

// GetCommentListByArticleId 前端获取评论列表
func GetCommentListByArticleId(articleId, pageSize, pageNum int) (comments []Comment, total int64, code int) {
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
		utils.Logger.Error(constant.ConvertForLog(constant.GetCommentListByArticleIdError), err)
		return comments, 0, constant.GetCommentListByArticleIdError
	}
	return comments, total, constant.SuccessCode
}

// GetCommentList 获取评论列表
func GetCommentList(pageNum, pageSize int) (comments []Comment, total int64, code int) {
	err := db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("created_at desc").Find(&comments).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.GetCommentListError), err)
		return nil, 0, constant.GetCommentListError
	}
	err = db.Model(comments).Count(&total).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.GetCommentListError), err)
		// TODO:返回评论列表但不返回评论数？？
		return comments, 0, constant.GetCommentListError
	}
	return comments, total, constant.SuccessCode
}

// UpdateCommentStatus 更新评论状态
func UpdateCommentStatus(id int, data *Comment) int {
	var comment Comment
	// 判断评论是否存在
	err := db.Where("id = ?", id).Take(&comment).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CommentNotExistError), err)
		return constant.CommentNotExistError
	}
	// 防止评论状态重复设置，一来可以避免无效的数据库查询，二来可以防止无效的设置导致的评论数自增
	if data.Status == comment.Status {
		utils.Logger.Error(constant.ConvertForLog(constant.CommentStatusRepeatSet))
		return constant.CommentStatusRepeatSet
	}

	// 判断是增加评论数还是减少评论数,前面已经排除了传入的status和表中现有的status相同的情况
	// 如果是从其他状态变成状态1，则评论数+1，如果是从状态1变成其他状态，则评论数-1，其他情况不变
	var count int
	if data.Status == 1 {
		count = 1
	} else if comment.Status == 1 {
		count = -1
	}
	// 使用map，字段为零值也会更新
	maps := map[string]interface{}{
		"status": data.Status,
	}

	// 更新评论状态
	err = db.Model(&comment).Updates(maps).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.UpdateCommentStatusError), err)
		return constant.UpdateCommentStatusError
	}
	// 更新文章的评论数
	err = db.Model(&Article{}).Where("id = ?", comment.ArticleId).UpdateColumn(
		"comment_count", gorm.Expr("comment_count + ?", count),
	).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.AddCommentCountError), err)
		return constant.AddCommentCountError
	}
	return constant.SuccessCode
}

//// TakeDownComment 撤销审核
//func TakeDownComment(id int, data *Comment) int {
//	maps := map[string]interface{}{
//		"status": data.Status,
//	}
//	var comment Comment
//	// 更新status
//	err := db.Model(&Comment{}).Where("id = ?", id).Updates(maps).Take(&comment).Error
//	if err != nil {
//		utils.Logger.Error(constant.ConvertForLog(constant.TakeDownCommentError), err)
//		return constant.TakeDownCommentError
//	}
//	// 评论数-1
//	err = db.Model(&Article{}).Where("id = ?", comment.ArticleId).UpdateColumn(
//		"comment_count", gorm.Expr("comment_count + ?", 1),
//	).Error
//	if err != nil {
//		utils.Logger.Error(constant.ConvertForLog(constant.ReduceCommentCountError), err)
//		return constant.ReduceCommentCountError
//	}
//	return constant.SuccessCode
//}

// DeleteComment 删除评论
func DeleteComment(id int) int {
	err := db.Where("id = ?", id).Delete(&Comment{}).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DeleteCommentError), err)
		return constant.DeleteCommentError
	}
	return constant.SuccessCode
}
