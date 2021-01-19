package errmsg

// 状态码
const (
	Success                = 200
	Error                  = 500
	UserAlreadyExistsError = 1001
)

// 状态码信息
var codeMsg = map[int]string{
	Success: "success",
	Error: "error",

	UserAlreadyExistsError: "用户已存在",
}

func GetCodeMsg(code int) string {
	return codeMsg[code]
}
