/*
constant 用于存放请求处理过程中可能出现的错误状态码及状态码所代表的含义
*/
package constant

import (
	"fmt"
)

// 状态码, 由5位构成，服务级别错误码：（eg：由1位表示，1表示系统级错误，2表示普通错误）
// 模块级错误码：由第二位和第三位两位构成（eg：01表示用户模块，02表示文章模块）
// 具体错误码：由于第四位和第五位表示，表示具体的错误（eg：01表示用户名错误，02表示密码错误等）
const (
	SuccessCode = 0

	// 系统级别错误
	ServerError           = 10001
	StartProjectError     = 10002
	ParamError            = 10003
	SetValidatorError     = 10004
	DataVerificationError = 10005
	// TODO:是否需要继续拆分code等级
	ReadConfigFileError     = 10006
	ReadServerConfigError   = 10007
	ReadDatabaseConfigError = 10008

	DatabaseAccessError         = 10009
	DatabaseConnectError        = 10010
	DatabaseMigrateError        = 10011
	DatabaseInstanceCreateError = 10012

	// 普通级别错误-01用户模块
	UserAlreadyExistsError  = 20101
	CreateUserError         = 20102
	UserNotExistError       = 20103
	GetUserListError        = 20104
	UsernameNotExistsError  = 20105
	UserPasswordError       = 20106
	EncryptPasswordError    = 20107
	SavePasswordError       = 20108
	UpdatePasswordError     = 20109
	UserRoleError           = 20110
	ChangeUserPasswordError = 20111
	EditUserInfoError       = 20112
	DeleteUserError         = 20113

	// 普通级别错误-02中间件模块
	SetTokenError         = 20201
	TokenMalformedError   = 20202
	TokenInvalidError     = 20203
	CheckTokenError       = 20204
	TokenNotExistError    = 20205
	CreateLogFileError    = 20206
	CreateRotateLogsError = 20207
	GetLocalHostnameError = 20208

	// 普通级别错误-03个人信息模块
	GetProfileInfoError    = 20301
	UpdateProfileInfoError = 20302

	// 普通级别错误-04文章分类模块
	CategoryExistError    = 20401
	CreateCategoryError   = 20402
	GetCategoryInfoError  = 20403
	CategoryNotExist      = 20404
	EditCategoryInfoError = 20405
	DeleteCategoryError   = 20406

	// 普通级别错误-05文章模块
	ArticleNotExistError            = 20501
	UpdateReadCountError            = 20502
	GetArticleListInfoError         = 20503
	CountArticleListError           = 20504
	GetArticleListByCategoryIdError = 20505
	CreateArticleError              = 20506
	EditArticleInfoError            = 20507
	DeleteArticleError              = 20508

	// 普通级别错误-06评论模块
	CreateCommentError      = 20601
	GetCommentInfoError     = 20602
	GetCommentCountError    = 20603
	GetCommentListError     = 20604
	CountCommentError       = 20605
	ApproveCommentError     = 20606
	AddCommentCountError    = 20607
	TakeDownCommentError    = 20608
	ReduceCommentCountError = 20609
	DeleteCommentError      = 20610
)

// CodeMsg 状态码信息字典,key和value分别对应状态码(int)和状态码信息(string)
var CodeMsg = map[int]string{
	SuccessCode: "ok",

	// 系统级别错误
	ServerError:           "服务异常",
	StartProjectError:     "项目启动失败",
	ParamError:            "参数错误",
	SetValidatorError:     "翻译器设置失败",
	DataVerificationError: "数据校验错误",

	ReadConfigFileError:     "配置文件读取失败",
	ReadServerConfigError:   "读取server配置错误",
	ReadDatabaseConfigError: "读取database配置错误",

	DatabaseAccessError:         "数据库访问异常",
	DatabaseConnectError:        "数据库连接失败",
	DatabaseMigrateError:        "数据库迁移失败",
	DatabaseInstanceCreateError: "创建数据库实例失败",

	// 模块级别错误-01用户模块
	UserAlreadyExistsError:  "用户已存在",
	CreateUserError:         "用户创建失败",
	UserNotExistError:       "用户不存在",
	GetUserListError:        "获取用户列表失败",
	UsernameNotExistsError:  "用户名不存在",
	UserPasswordError:       "密码错误",
	EncryptPasswordError:    "密码加密错误",
	SavePasswordError:       "密码保存失败",
	UpdatePasswordError:     "密码更新失败",
	UserRoleError:           "用户角色码错误",
	ChangeUserPasswordError: "密码修改失败",
	EditUserInfoError:       "用户信息编辑失败",
	DeleteUserError:         "用户删除失败",

	// 普通级别错误-02中间件模块
	SetTokenError:         "token设置失败",
	TokenMalformedError:   "token格式错误",
	TokenInvalidError:     "token已失效",
	CheckTokenError:       "token不正确",
	TokenNotExistError:    "token不存在",
	CreateLogFileError:    "日志文件创建失败",
	CreateRotateLogsError: "轮询日志创建失败",
	GetLocalHostnameError: "本地系统主IP获取失败",

	// 普通级别错误-03个人信息模块
	GetProfileInfoError:    "个人信息获取失败",
	UpdateProfileInfoError: "个人信息更新错误",

	// 普通级别错误-04文章分类模块
	CategoryExistError:    "分类已存在",
	CreateCategoryError:   "新增分类失败",
	GetCategoryInfoError:  "分类获取失败",
	CategoryNotExist:      "分类不存在",
	EditCategoryInfoError: "文章分类编辑错误",
	CreateArticleError:    "文章创建失败",
	DeleteCategoryError:   "文章分类删除失败",

	// 普通级别错误-05文章模块
	ArticleNotExistError:            "文章信息获取失败",
	UpdateReadCountError:            "阅读量更新失败",
	GetArticleListInfoError:         "文章列表获取失败",
	CountArticleListError:           "文章列表总数获取失败",
	GetArticleListByCategoryIdError: "该分类文章列表获取失败",
	EditArticleInfoError:            "文章编辑失败",
	DeleteArticleError:              "文章删除失败",

	// 普通级别错误-06评论模块
	CreateCommentError:      "评论创建失败",
	GetCommentInfoError:     "评论获取失败",
	GetCommentCountError:    "评论数获取失败",
	GetCommentListError:     "评论列表获取失败",
	CountCommentError:       "评论数获取失败",
	ApproveCommentError:     "评论通过失败",
	AddCommentCountError:    "评论数增加失败",
	TakeDownCommentError:    "评论撤销失败",
	ReduceCommentCountError: "评论数减少失败",
	DeleteCommentError:      "评论删除错误",
}

// 将code转换为末尾带有“：”的message, 用于打印log信息
func ConvertForLog(code int) string {
	return fmt.Sprintf("code: %d, msg: %s, err: ", code, CodeMsg[code])
}
