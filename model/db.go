package model

import (
	"GopherBlog/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var db *gorm.DB
var err error

func InitDB() {
	// data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
		utils.DatabaseCfg.Password,
		utils.DatabaseCfg.Username,
		utils.DatabaseCfg.Host,
		utils.DatabaseCfg.Port,
		utils.DatabaseCfg.Name,
	)

	// gorm.Open接受两个参数，一个是“数据库方言器”（dialector），同一种意思，不同的方言有不同的发音，数据库也一样，mysql有mysql的语法，
	// oracle有oracle的语法，而对用用户来说，不管底层用的是MySQL还是Oracle，上层的orm接口都一样，用户只需要告诉orm使用哪一种"方言"即可
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用默认事务（为保证数据一致性，gorm默认会创建事务来执行单个的创建、删除、更新操作，由此可以推断默认使用的是innodb）
		// 由于创建一个事务都需要消耗一定的资源，所以如果没有特别要求可以禁用该默认选项, 能获得30%+性能提升
		SkipDefaultTransaction:                   true,
		// 命名策略，就是控制表明或者列表按照单数或者复数形式来命名，例如用户表，如果设置表名为复数形式则是users，单数形式为user
		NamingStrategy:                           schema.NamingStrategy{
			SingularTable: true,
		},
		//FullSaveAssociations:                     false,

		// 最低级别，应该相当于trace级别
		Logger:                                   logger.Default.LogMode(logger.Silent),
		//NowFunc:                                  nil,
		//DryRun:                                   false,
		//PrepareStmt:                              false,
		//DisableAutomaticPing:                     false,

		// 在迁移时禁用外键约束, ?????
		DisableForeignKeyConstraintWhenMigrating: true,
		//DisableNestedTransaction:                 false,
		//AllowGlobalUpdate:                        false,
		//QueryFields:                              false,
		//CreateBatchSize:                          0,
		//ClauseBuilders:                           nil,
		//ConnPool:                                 nil,
		//Dialector:                                nil,
		//Plugins:                                  nil,
	})

	if err != nil {
		utils.Logger.Error("数据库连接失败: ", err)
	}
	// 迁移数据库
	err = db.AutoMigrate()
	if err != nil {
		utils.Logger.Error("数据库迁移失败： ", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		utils.Logger.Error("创建数据库实例失败： ", err)
	}
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大连接数, 当设置的最大空闲连接数超过最大连接数，则会将最大空闲连接数重置为与最大连接数相等
	sqlDB.SetMaxOpenConns(100)
	// 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(10 * time.Second)
}
