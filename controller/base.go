/*
controller 定义控制器层逻辑
*/
package controller

import (
	"GopherBlog/constant"
	"github.com/gin-gonic/gin"
	"net/http"
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

func DetailParams(f func()) {

}
