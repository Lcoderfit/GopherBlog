package model

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
)

// 个人信息
type Profile struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"type:varchar(20)" json:"name"`
	Desc   string `gorm:"type:varchar(200)" json:"desc"`
	QqChat string `gorm:"type:varchar(200)" json:"qq_chat"`
	WeChat string `gorm:"type:varchar(200)" json:"wechat"`
	Weibo  string `gorm:"type:varchar(200)" json:"weibo"`
	Bili   string `gorm:"type:varchar(200)" json:"bili"`
	Email  string `gorm:"type:varchar(200)" json:"email"`
	Img    string `gorm:"type:varchar(200)" json:"img"`
	Avatar string `gorm:"type:varchar(200)" json:"avatar"`
}

// 获取个人信息
func GetProfileInfo(id int) (Profile, int) {
	var data Profile
	// 或db.Select("id").Take(&data, id)
	err := db.Where("id = ?", id).Take(&data).Error
	if err != nil {
		// 1.内函数打log，则不返回err，只返回err code，外函数直接通过判断是否与constant.SuccessCode相等从而返回对应的类型的响应
		// 2.内函数不打log，则需要返回err,可以不返回err code, 外函数打log
		utils.Logger.Error(constant.ConvertForLog(constant.GetProfileInfoError), err)
		return data, constant.GetProfileInfoError
	}
	return data, constant.SuccessCode
}
