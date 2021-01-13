package errmsg

// 状态码
const (

)

// 状态码信息
var codeMsg = map[int]string{

}

func GetErrMsg(code int) string {
	return  codeMsg[code]
}