package middleware

// TODO:拦截请求，将参数校验和err处理放在中间件中进行, 使得请求处理函数可以直接使用数据
// TODO:当在中间件中参数校验失败时，直接调用c.JSON返回（可以封装成fail函数）
