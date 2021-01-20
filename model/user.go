package model

import (
	"gorm.io/gorm"
)

// 用户模型
type User struct {
	// gorm.Model是一个包含ID、CreateAt、UpdateAt、DeleteAt字段的结构体，可以内嵌到自定义模型中
	// 这样自己的模型就带有ID（主键）、CreateAt、UpdateAt、DeleteAt等的字段
	// 自增ID作为主键，默认从1开始自增
	gorm.Model
	// gorm标签后面的每个子标签名与设置的值格式为 -》 gorm:"tag1:v1;tag2:v2"
	// validate标签格式： validate:"tag1,tag2=v2,tag3=v3"
	Username string `gorm:"type:varchar(20);not null" json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(20);not null" json:"password" validate:"required,min=6,max=20" label:"密码"`
	Role     int    `gorm:"type:int;default:2" json:"role" validate:"required,gte=2" label:"角色码"`
}

// 检查用户是否存在
func IsUserExists(name string) bool {
	var user User
	// 1.字段名用数据库字段名称或者模型中的字段名都可以
	//
	// 2.Take、First、Find这些函数会根据传入参数选择查询对应的表，例如Take(&user)表示选择user表，并且将查询结果存放到user中
	// 应为Take和First只查询单条数据，所以user中会存储查询到的值，也可以用 result := .....Take(&user)来获取查询到的结果数量,
	// 但是Find会查询多条数据，需要传入users（user类型的切片）作为参数来存储多个值，result.RowsAffected表示查询到的数据条数
	//
	// 3.Take相当于:limit 1， First相当于：order by id limit 1, 所以一般Take效率比较高
	db.Select("id").Where("username = ?", name).Take(&user)
	if user.ID > 0 {
		return true
	}
	return false
}

// 创建新用户
//func CreateUser(user *User) error {
//	db.Cre
//}
