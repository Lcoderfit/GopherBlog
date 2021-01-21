package controller

import (
	"GopherBlog/constant"
	"GopherBlog/model"
	"GopherBlog/utils"
	"GopherBlog/utils/validator"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 添加用户
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
		failWithData(c, constant.DataVerificationError, msg)
		// 当上面的响应报错时，调用c.Abort()可以确保其不会影响下面的请求处理程序
		//c.Abort()
	}

	if ok := model.IsUserExists(data.Username); ok {
		// 能否重写Info函数
		utils.Logger.Info(constant.ConvertForLog(constant.UserAlreadyExistsError))
		fail(c, constant.UserAlreadyExistsError)
	}

	// TODO:创建model需要传入指针，为什么？？
	if err := model.CreateUser(&data); err != nil {
		fail(c, constant.CreateUserError)
	}
	// TODO：是否需要对请求成功的情况进行状态码分级
	success(c)
}

// 获取用户信息
func GetUserInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), err)
	}

	user, err := model.GetUserInfoById(id)
	if err != nil {
		fail(c, constant.GetUserInfoError)
	}
	successWithData(c, user)
}

// 获取用户列表
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
	}
	data := map[string]interface{}{
		"users": users,
		"total": total,
	}
	successWithData(c, data)
}

//
