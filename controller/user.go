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
	}

	msg, err := validator.Validate(data)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.DataVerificationError), err)
		failWithData(c, constant.DataVerificationError, msg)
	}

	if ok := model.IsUserExist(data.Username); ok {
		// 能否重写Info函数
		utils.Logger.Info(constant.ConvertForLog(constant.UserAlreadyExistsError))
		fail(c, constant.UserAlreadyExistsError)
		return
	}

	// TODO:创建model需要传入指针，为什么？？
	if err := model.CreateUser(&data); err != nil {
		fail(c, constant.CreateUserError)
		return
	}
	// TODO：是否需要对请求成功的情况进行状态码分级
	success(c)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
	}
	var data model.User
	// 这里报错不直接返回，因为需要后面调用验证器返回哪个参数有问题的信息
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
		return
	}

	user, err := model.GetUserInfoById(id)
	if err != nil {
		fail(c, constant.GetUserInfoError)
		return
	}
	successWithData(c, user)
}

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	pageNum, err := strconv.Atoi(c.Query("page_num"))
	username := c.Query("username")

	// 对每一页的数量上下限进行限制
	// TODO: 逻辑是否正常
	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	users, total, err := model.GetUserList(pageSize, pageNum, username)
	if err != nil {
		fail(c, constant.GetUserListError)
		return
	}
	data := map[string]interface{}{
		"users": users,
		"total": total,
	}
	successWithData(c, data)
}

// ChangeUserPassword JWT鉴权接口:修改用户密码
// TODO:需要鉴权的接口都需要传入data参数？?
func ChangeUserPassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
		return
	}

	var data model.User
	// TODO：为什么这里不需要使用validator进行验证
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
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
	//id, err := strconv.Atoi(c.Param("id"))
	//if err != nil {
	//	utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
	//	fail(c, constant.ParamError)
	//}
	var data model.User
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	code := model.EditUserInfo(&data)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	success(c)
}

// DeleteUser JWT鉴权接口:删除用户
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}
	var data model.User
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
		fail(c, constant.ParamError)
	}

	if ok := model.IsUserExist(data.Username); !ok {
		utils.Logger.Info(constant.ConvertForLog(constant.UserAlreadyExistsError))
		fail(c, constant.UserAlreadyExistsError)
	}
	code := model.DeleteUser(id)
	if code != constant.SuccessCode {
		fail(c, code)
	}
	success(c)
}
