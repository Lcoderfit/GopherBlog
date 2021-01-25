package middleware

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
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
		return nil, constant.TokenIsNilError, err
	}
	if key, ok := parsedJwt.Claims.(*MyClaim); ok && parsedJwt.Valid {
		return key, constant.SuccessCode, nil
	}
	return nil, constant.CheckTokenError, errors.New("")
}
