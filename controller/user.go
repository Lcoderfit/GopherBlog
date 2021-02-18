package controller

import (
	"GopherBlog/constant"
	"GopherBlog/model"
	"GopherBlog/utils"
	"GopherBlog/utils/validator"
	"github.com/gin-gonic/gin"
	"strconv"
)

// AddUser 添加用户
func AddUser(c *gin.Context) {
	// 1.声明一个模型
	// 2.参数校验
	// 3.数据校验
	// 4.调用模型中声明的与数据库交互的函数进行数据交互
	var data model.User
	// 这里报错不直接返回，因为需要后面调用验证器返回哪个参数有问题的信息
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		return
	}

	msg, err := validator.Validate(data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DataVerificationError), err)
		failWithData(c, constant.DataVerificationError, msg)
		return
	}

	if ok := model.IsUserExistsByName(data.Username); ok {
		// 能否重写Info函数
		utils.Logger.Error(constant.ConvertForLog(constant.UserAlreadyExistsError))
		fail(c, constant.UserAlreadyExistsError)
		return
	}

	// TODO:创建model需要传入指针，为什么？？
	if code := model.CreateUser(&data); code != constant.SuccessCode {
		fail(c, code)
		return
	}
	// TODO：是否需要对请求成功的情况进行状态码分级
	success(c)
}

// GetUserInfo 获取单个用户信息
func GetUserInfo(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		return
	}

	user, code := model.GetUserInfoById(id)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	successWithData(c, user)
}

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.Query("page_number"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	username := c.Query("username")

	// 对每一页的数量上下限进行限制
	// TODO: 逻辑是否正常
	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	if pageNum <= 0 {
		pageNum = 1
	}

	users, total, code := model.GetUserList(pageSize, pageNum, username)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	data := map[string]interface{}{
		"user_list": users,
		"total":     total,
	}
	successWithData(c, data)
}

// ChangeUserPassword JWT鉴权接口:修改用户密码
// TODO:需要鉴权的接口都需要传入data参数？?
func ChangeUserPassword(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		fail(c, constant.ParamError)
		return
	}

	var data model.User
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
		return
	}
	msg, err := validator.Validate(data)
	if err != nil {
		failWithData(c, constant.DataVerificationError, msg)
		return
	}

	code := model.ChangeUserPassword(id, &data)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	success(c)
}

// EditUserInfo JWT鉴权接口:编辑用户信息
func EditUserInfo(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
		return
	}
	// 需要编辑的信息通过json格式参数传入
	var data model.UserEdition
	//var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
		return
	}
	msg, err := validator.Validate(data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DataVerificationError), err)
		failWithData(c, constant.DataVerificationError, msg)
		return
	}

	code := model.EditUserInfo(id, &data)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	success(c)
}

// DeleteUser JWT鉴权接口:删除用户
func DeleteUser(c *gin.Context) {
	id, err := MustInt(c.Param, "id")
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
		return
	}

	if ok := model.IsUserExistsById(id); !ok {
		utils.Logger.Info(constant.ConvertForLog(constant.UserNotExistsError))
		fail(c, constant.UserNotExistsError)
		return
	}
	code := model.DeleteUser(id)
	if code != constant.SuccessCode {
		fail(c, code)
		return
	}
	success(c)
}

// 测试接口
func Test(c *gin.Context) {
	var data model.User
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
		return
	}
	successWithData(c, data)
}
