package middleware

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 自定义jwt的payload部分（所携带的信息）
type MyClaim struct {
	username string
	jwt.StandardClaims
}

// SignedString函数需要传入字节切片
var JwtKey = []byte(utils.JwtKey)

// 设置token
func SetToken(username string) (string, int) {
	// 设置token的过期时间，超出该时间token失效，需要重新登录
	expiredTime := time.Now().Add(7 * 24 * time.Hour)
	// 实例化claim
	myClaim := MyClaim{
		username: username,
		StandardClaims: jwt.StandardClaims{
			// 传入unix时间戳
			ExpiresAt: expiredTime.Unix(),
			// 发行人，可以不设置
			Issuer: "GopherBlog",
		},
	}

	// 第一个参数用于设置header，第二个参数用于设置payload
	newJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaim)
	// 根据选定的加密算法对header和payload进行加密，返回signature（即加密后所得到的一串字符串）
	token, err := newJwt.SignedString(JwtKey)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.SetTokenError), err)
		return "", constant.SetTokenError
	}
	return token, constant.SuccessCode
}

// 校验token
func CheckToken(token string) (*MyClaim, int, error) {
	var myClaim MyClaim
	// 解析出jwt
	parsedJwt, err := jwt.ParseWithClaims(token, &myClaim, func(token *jwt.Token) (i interface{}, e error) {
		return JwtKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, constant.TokenMalformedError, err
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, constant.TokenInvalidError, err
			} else {
				// TODO:token格式错误??
				return nil, constant.CheckTokenError, err
			}
		}
	}

	if parsedJwt == nil {
		return nil, constant.TokenNotExistError, err
	}
	if key, ok := parsedJwt.Claims.(*MyClaim); ok && parsedJwt.Valid {
		return key, constant.SuccessCode, nil
	}
	return nil, constant.CheckTokenError, errors.New("")
}

// jwt中间件
// TODO：将c.JSON替换为响应函数
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization: 授权
		tokenHeader := c.Request.Header.Get("Authorization")
		code := constant.TokenNotExistError
		// token不存在
		if tokenHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":  code,
				"error": constant.CodeMsg[code],
			})
			// 直接返回，该中间件中c.Abort()后面的部分将不再执行，当在该中间件之后的其他中间件还是会执行
			c.Abort()
			// TODO:为什么还要return
			return
		}
		// 正常的规则为： bearer xxxxxxx
		tokenInfo := strings.Split(tokenHeader, " ")
		if len(tokenInfo) != 2 || tokenInfo[0] != "bearer" {
			code = constant.TokenMalformedError
			c.JSON(http.StatusOK, gin.H{
				"code":  code,
				"error": constant.CodeMsg[code],
			})
			c.Abort()
			return
		}

		myClaim, code, err := CheckToken(tokenInfo[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  code,
				"error": constant.CodeMsg[code],
			})
			c.Abort()
			return
		}
		// 设置上下文变量
		c.Set("username", myClaim)
		return
	}
}
