package utils

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// 可以定义任意数量的实例
var Logger = logrus.New()

var (
	ServerCfg   *ServerConfig
	DatabaseCfg *DatabaseConfig
)

// 服务器配置
type ServerConfig struct {
	AppMode  string
	HttpPort string
}

// 数据库配置
type DatabaseConfig struct {
	Type     string
	Name     string
	Host     string
	Port     string
	Username string
	Password string
}

func init() {
	// 设置日志格式
	Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// 设置日志级别， logrus默认级别是什么？？
	Logger.SetLevel(logrus.DebugLevel)
	// 设置打印信息中显示调用打印语句的函数所在的文件路径及函数名和在文件中的哪一行
	Logger.SetReportCaller(true)

	cfg, err := ini.Load("../config/config.ini")
	if err != nil {
		Logger.Error("config.ini读取错误: ", err)
	}
	err = cfg.Section("server").MapTo(ServerCfg)
	if err != nil {
		Logger.Error("读取server配置错误: ", err)
	}
	err = cfg.Section("database").MapTo(DatabaseCfg)
	if err != nil {
		Logger.Error("读取database配置错误： ", err)
	}
}
