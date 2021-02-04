package middleware

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"time"
)

// Log 日志中间件
func Log() gin.HandlerFunc {
	// 创建日志文件，必须先创建log目录(OpenFile会创建最后的文件，但是不会创建父目录)
	// 否则会报找不到路径的错误(The system cannot find the path specified)
	filePath := "log/output.log"
	src, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		// TODO:创建失败了怎么办？？
		utils.Logger.Error(constant.ConvertForLog(constant.CreateLogFileError), err)
	}

	// 创建logrus实例
	logger := logrus.New()
	// 设置日志级别，会打印该日志级别及以上级别的日志
	logger.SetLevel(logrus.DebugLevel)
	// 将要打印的日志内容重定向到src
	logger.Out = src

	// 创建轮询日志（file-rotatelogs具有日志轮询机制，就是定期清理日志，防止日志文件过大）
	logWriter, err := rotatelogs.New(
		filePath+"%Y%m%d.log",                     // 需要轮询处理的日志文件的格式
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 设置日志被从系统中清除之前的最大保存时间），一周清理一次
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割的时间间隔，一天切割一次
	)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.CreateRotateLogsError), err)
	}

	// 将logrus日志级别映射到io.Writer,多个级别可以共享一个io.Writer, 一个级别不能使用多个io.Writer
	writerMap := lfshook.WriterMap{
		logrus.DebugLevel: logWriter,
		logrus.InfoLevel:  logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.FatalLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	// 创建hook函数,TODO:
	hook := lfshook.NewHook(writerMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// 添加hook到logrus
	logger.AddHook(hook)

	return func(c *gin.Context) {
		startTime := time.Now()
		// 将程序挂起，执行后面的流程，等返回响应后并执行完其他Next后，会继续执行下面的路程
		c.Next()
		// 计算后续中间件及请求处理函数用时
		entTime := time.Since(startTime)
		// math.Ceil: 向上取整
		spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(entTime.Nanoseconds())/1000000.0)))
		// 获取本地系统地址
		hostname, err := os.Hostname()
		if err != nil {
			utils.Logger.Error(constant.ConvertForLog(constant.GetLocalHostnameError), err)
			hostname = "Unknown"
		}
		// 获取状态码、客户端IP、用户代理、数据大小、请求方法和URI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		dataSize := c.Writer.Size()
		// TODO: dataSize会小于0吗
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method
		uri := c.Request.RequestURI
		// logrus字段机制，每次打印日志都会额外带上配置的字段
		log := logrus.WithFields(logrus.Fields{
			"Hostname":  hostname,
			"Status":    statusCode,
			"SpentTime": spendTime,
			"IP":        clientIP,
			"Agent":     userAgent,
			"DataSize":  dataSize,
			"Method":    method,
			"URI":       uri,
		})
		if len(c.Errors) > 0 {
			// ErrorTypePrivate表示私人错误, TODO
			log.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode >= 500 {
			log.Error()
		} else if statusCode >= 400 {
			log.Warn()
		} else {
			log.Info()
		}
	}
}
