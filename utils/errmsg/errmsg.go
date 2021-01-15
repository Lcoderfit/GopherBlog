package errmsg

// 状态码
const (
	SUCCESS = 200
	ERROR   = 500
)

// 状态码信息
var codeMsg = map[int]string{}

func GetCodeMsg(code int) string {
	return codeMsg[code]
}
