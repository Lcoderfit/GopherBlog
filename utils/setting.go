package utils

import (
	"GopherBlog/constant"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// 可以定义任意数量的实例
var Logger = logrus.New()

var (
	JwtKey      string
	ServerCfg   = new(ServerConfig)
	DatabaseCfg = new(DatabaseConfig)
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
	// 当使用Fatal及以上级别打印日志消息时，调用Fatal或Panic函数，在调用位置之后的程序将不会执行
	Logger.SetLevel(logrus.DebugLevel)
	// 设置打印信息中显示调用打印语句的函数所在的文件路径及函数名和在文件中的哪一行
	Logger.SetReportCaller(true)

	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		Logger.Error(constant.ConvertForLog(constant.ReadConfigFileError), err)
	}
	err = cfg.Section("server").MapTo(ServerCfg)
	if err != nil {
		Logger.Error(constant.ConvertForLog(constant.ReadServerConfigError), err)
	}
	err = cfg.Section("database").MapTo(DatabaseCfg)
	if err != nil {
		Logger.Error(constant.ConvertForLog(constant.ReadDatabaseConfigError), err)
	}
	// 获取jwt加密密钥
	JwtKey = cfg.Section("server").Key("JwtKey").String()
}
