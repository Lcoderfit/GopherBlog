package model

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	// gorm.Model是一个包含ID、CreateAt、UpdateAt、DeleteAt字段的结构体，可以内嵌到自定义模型中
	// 这样自己的模型就带有ID（主键）、CreateAt、UpdateAt、DeleteAt等的字段
	// 自增ID作为主键，默认从1开始自增
	gorm.Model
	// 1.gorm标签后面的每个子标签名与设置的值格式为 -》 gorm:"tag1:v1;tag2:v2"
	// 2.validate标签格式： validate:"tag1,tag2=v2,tag3=v3"
	// 3.json标签定义的是结构体转换为json数据时对应的字段名称
	// 4.设置结构体时注意字段对应的gorm标签类型，例如什么时候设主键，什么时候需要设置primaryKey,什么时候设not null
	Username string `gorm:"type:varchar(20);not null" json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(20);not null" json:"password" validate:"required,min=6,max=20" label:"密码"`
	Role     int    `gorm:"type:int;default:2" json:"role" validate:"required,gte=2" label:"角色码"`
}

// IsUserExist 检查用户是否存在
func IsUserExist(name string) bool {
	var user User
	// 1.字段名用数据库字段名称或者模型中的字段名都可以
	//
	// 2.Take、First、Find这些函数会根据传入参数选择查询对应的表，例如Take(&user)表示选择user表，并且将查询结果存放到user中
	// 因为Take和First只查询单条数据，所以user中会存储查询到的值，也可以用 result := .....Take(&user)来获取查询到的结果数量,
	// 但是Find会查询多条数据，需要传入users（user类型的切片）作为参数来存储多个值，result.RowsAffected表示查询到的数据条数
	//
	// 3.Take相当于:limit 1， First相当于：order by id limit 1, 所以一般Take效率比较高
	db.Select("id").Where("username = ?", name).Take(&user)
	if user.ID > 0 {
		return true
	}
	return false
}

// CreateUser 创建新用户, 实现BeforeSave和BeforeUpdate接口，对密码进行创建和更新时都会自动进行加密
func CreateUser(data *User) error {
	err := db.Create(data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CreateUserError), err)
		return err
	}
	return nil
}

// BeforeSave 保存密码到数据库中时，自动对密码进行加密
func (u *User) BeforeSave(_ *gorm.DB) (err error) {
	u.Password, err = EncryptPassword(u.Password)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.SavePasswordError))
		return
	}
	utils.Logger.Info("password: ", u.Password)
	return
}

// BeforeUpdate 更新密码在保存到数据库之前自动对密码进行加密
func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	u.Password, err = EncryptPassword(u.Password)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.UpdatePasswordError), err)
		return
	}
	return
}

// EncryptPassword 对密码进行bcrypt加密
func EncryptPassword(password string) (string, error) {
	const cost = 10
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.EncryptPasswordError), err)
		return "", err
	}
	return string(hashPwd), nil
}

// GetUserInfoById 通过用户ID获取用户数据
func GetUserInfoById(id int) (User, int) {
	var user User
	// 如果没找到会返回record not found错误
	err := db.Where("id = ?", id).Take(&user).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.UserNotExistError), err)
		return user, constant.UserNotExistError
	}
	return user, constant.SuccessCode
}

// GetUserList 获取用户列表
// 可能存储非常多的数据，total需要用长整型
func GetUserList(pageSize, pageNum int, username string) (users []User, total int64, err error) {
	if username != "" {
		err = db.Select("id, username, role").Where(
			"username like ?", "%"+username+"%",
		).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Error
	} else {
		err = db.Select("id, username, role").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Error
	}
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.GetUserListError), err)
		return nil, total, err
	}
	total = int64(len(users))
	return users, total, nil
}

// CheckAccount 检查账户密码是否正确
func CheckAccount(username, password string) (user User, code int) {
	err = db.Where("username = ?", username).Take(&user).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DatabaseAccessError), err)
		return user, constant.DatabaseAccessError
	}
	if user.ID == 0 {
		utils.Logger.Error(constant.ConvertForLog(constant.UsernameNotExistsError))
		return user, constant.UsernameNotExistsError
	}

	// 判断密码是否是已加密密码对应的明文
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.UserPasswordError), err)
		return user, constant.UserPasswordError
	}
	if user.Role != 1 {
		utils.Logger.Error(constant.ConvertForLog(constant.UserRoleError), err)
		return user, constant.UserRoleError
	}
	return user, constant.SuccessCode
}

// ChangeUserPassword 修改用户密码
// TODO:跟一般的逻辑比是否过于简单
func ChangeUserPassword(id int, data *User) int {
	// Select可以选择需要修改的字段， 下面的语句等效于：update user set=xxx where id=yyy
	// xxx表示Where中的id，yyy表示data中的ID值
	err := db.Select("password").Where("id = ?", id).Updates(&data).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ChangeUserPasswordError), err)
		return constant.ChangeUserPasswordError
	}
	return constant.SuccessCode
}

// EditUserInfo 编辑用户信息
func EditUserInfo(data *User) int {
	// TODO：传入data和传入&data的区别？？
	var maps = make(map[string]interface{})
	maps["username"] = data.Username
	maps["role"] = data.Role
	// 假设data中的ID为111， 下面语句相当于update user set username=xxx and role=yyy where id='111'
	err := db.Model(&data).Updates(maps)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.EditUserInfoError), err)
		return constant.EditUserInfoError
	}
	return constant.SuccessCode
}

// DeleteUser 删除用户
func DeleteUser(id int) int {
	// 如果模型中定义了gorm.DeleteAt字段，则会具有“软删除”功能, 即并不是真正的删除，而是将DeleteAt字段更新为删除语句执行的时间
	// 可以通过db.UnScoped().Where()查询软删除的数据，要永久删除可以使用db.UnScoped().Delete()
	err := db.Where("id = ?", id).Delete(&User{}).Error
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DeleteUserError), err)
		return constant.DeleteUserError
	}
	return constant.SuccessCode
}
