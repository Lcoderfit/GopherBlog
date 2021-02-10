/*
controller 定义控制器层逻辑
*/
package controller

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// success 请求成功
func success(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    constant.SuccessCode,
		"message": constant.CodeMsg[constant.SuccessCode],
		"data":    nil,
	})
}

func successWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    constant.SuccessCode,
		"message": constant.CodeMsg[constant.SuccessCode],
		"data":    data,
	})
}

func successWithStatusCode(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"code":    constant.SuccessCode,
		"message": constant.CodeMsg[constant.SuccessCode],
		"data":    data,
	})
}

// fail 请求失败, 如果请求失败了，则没有必要返回data
func fail(c *gin.Context, code int) {
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"error": constant.CodeMsg[code],
	})
}

// failWithData 请求失败，并返回data
func failWithData(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"error": constant.CodeMsg[code],
		"data":  data,
	})
}

func MustInt(f func(string) string, param string) (result int, err error) {
	result, err = strconv.Atoi(f(param))
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.ParamError), param+", ", err)
		return result, err
	}
	return result, nil
}

func MustIntArray(f func(string) string, params ...string) (results []int, err error) {
	for _, v := range params {
		param, err := strconv.Atoi(f(v))
		if err != nil {
			utils.Logger.Error(constant.ConvertForLog(constant.ParamError), v+", ", err)
			return nil, err
		}
		results = append(results, param)
	}
	return results, nil
}
