package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

// cors详见：http://www.ruanyifeng.com/blog/2016/04/cors.html
func Cors() gin.HandlerFunc {
	// 1.CORS分为简单请求和非简单请求两种，非简单请求第一次发送一个预检请求，预检通过后会根据MaxAge的设定作为该预检请求的有效期，在有效期内不用再
	// 发送预检请求，之后的请求就跟简单请求一样
	// 2.预检请求的作用：用来询问网页所在的域名是否被服务器所认可，并确认服务器支持的HTTP动词和头部字段
	return cors.New(
		cors.Config{
			AllowAllOrigins: true, // true表示允许任意域名进行跨域请求，该项被设置则无需设置AllowOrigins
			// AllowOrigins:           nil,	//[]string类型，表示被允许进行跨域请求的的域名
			// 是一个用于验证origin是否有效的自定义函数,如果有效则返回true,否则返回false,如果实现了该函数会忽略AllowOrigins的配置
			// AllowOriginFunc: nil,
			// 允许客户端用于跨域请求的方法列表(或者说服务端支持的跨域请求方法)，默认是简单方法（POST和GET）
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			// 表示服务端允许客户端用于非简单请求的头部字段（服务端支持的非简单请求头部字段）
			AllowHeaders: []string{"*"},
			// 指示客户端CORS请求头中是否可以包含用户凭证（如cookie，HTTP身份验证或客户端SSL证书）
			// 注意：如果要发送cookie，则Access-Control-Allow-Origin就不能设置为*,否则无法上传，必须指定与请求网页一致的域名
			// cookie也遵循同源策略，只有服务端域名设置的cookie才能上传，其他域的不能
			// TODO:本项目用的是Authorization传递token，所以不受影响？？
			AllowCredentials: true,
			// CORS请求时，XMLHttpRequest对象的getResponseHeader()方法只能拿到6个基本字段，如果需要获取其他字段则需要在该设置中配置
			// 用于设置前端API在接受到CORS响应时可以获取到的额外字段
			ExposeHeaders: nil,
			// 指示一个“预检”请求可以被缓存多长时间，在此期间不用再发出另一条预检请求
			MaxAge: 12 * time.Hour,
			// true表示支持用通配符来设置允许跨域请求的域名(例如设置为： AllowOrigins: []string{"http://some-domain/*", "https://api.*"})
			// AllowWildcard: false,
			// 是否允许使用浏览器扩展
			// AllowBrowserExtensions: false,
			// 是否允许使用WebSockets协议
			// AllowWebSockets: false,
			// 是否允许使用file://schema (危险！）只有在你百分之百确定需要的时候才使用它
			// AllowFiles: false,
		},
	)
}
