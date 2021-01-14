package model

import (
	"gorm.io/gorm"
)

// 用户模型
type User struct{
	// gorm.Model是一个包含ID、CreateAt、UpdateAt、DeleteAt字段的结构体，可以内嵌到自定义模型中
	// 这样自己的模型就带有ID（主键）、CreateAt、UpdateAt、DeleteAt等的字段
	gorm.Model
	// gorm标签后面的每个子标签名与设置的值格式为 -》 gorm:"tag1:v1;tag2:v2"
	// validate标签格式： validate:"tag1,tag2=v2,tag3=v3"
	Username string `gorm:"type:varchar(20);not null" json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(20);not null" json:"password" validate:"required,min=6,max=20" label:"密码"`
	Role int `gorm:"type:int;default:2" json:"role" validate:"required,gte=2" label:"角色码"`
}


