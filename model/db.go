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
		SkipDefaultTransaction: true,
		// 命名策略，就是控制表明或者列表按照单数或者复数形式来命名，例如用户表，如果设置表名为复数形式则是users，单数形式为user
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		//FullSaveAssociations:                     false,

		// 最低级别，应该相当于trace级别
		Logger: logger.Default.LogMode(logger.Silent),
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

	// 设置连接池最大空闲连接数，由于数据库建立连接很消耗资源，而连接池的原理即是事先建立好一定数量的连接放着
	// 如果需要就直接从池子里面拿，并将其标记为忙，用完后放回，标记为空闲。这样就不用每次都建立数据库连接，节省资源
	// 建立连接，需要消耗的资源：内存，握手通信啥的。。。
	//
	// 1.一开始调用Open函数的时候并不会创建连接，只有真正用到的时候才会建立连接，
	// 2.连接最长存活时间通过SetConnMaxLifetime设置, 当超过这个时间连接不会被复用，（默认值为0，即永不过期）
	// 3.建立一个连接后，用完不会立即释放，而是放入连接池内，标记为空闲状态，等需要的时候直接拿出来，
	// 4.最大空闲连接数由SetMaxIdleConns设置，当所需要的连接数超过连接池中已有连接，则会建立新的连接，
	// 用完后会先检查该连接在使用过程中是否已经无效，无效则不放入连接池，直接释放
	//
	// 5.如果该连接用完后仍有效，则检查是否有等待连接的阻塞请求，有则直接给最早的请求，没有则判断连接池是否达到上限，未达到上限则放回连接池；
	// 达到上限则直接释放连接。
	// 6.总连接数=正在使用的连接数+空闲连接数； 总连接数通过SetMaxOpenConns设置，默认等于0（即没有限制），如果总连接数超过这个值，则会
	// 将连接请求放入队列中
	//
	// 1.设置最大空闲连接数, 如果不调用该函数进行设置,则默认为2, 如果设置的数量 n <= 0,
	// 则db.maxIdle为-1,即没有可以复用的连接(连接池的最大空闲连接数为0条)
	//
	// 2.当设置的最大空闲连接数超过最大连接数，则会将最大空闲连接数重置为与最大连接数相等
	// 3.如果当前连接池中的空闲连接数超过设置的最大空闲连接数，则会强制关闭
	sqlDB.SetMaxIdleConns(10)
	// 设置最大连接数, 当设置的最大空闲连接数超过最大连接数，则会将最大空闲连接数重置为与最大连接数相等, 默认为0，即没有限制
	sqlDB.SetMaxOpenConns(100)
	// 1.设置连接可复用的最大时间, 指的是连接池中的空闲连接如果超过这个时间就会断开连接, 默认为0, 即永远复用
	// 源码中maxLifetimeClosed是由于连接池中空闲连接超过了MaxLifetime而被强制关闭的连接数
	// maxIdleClosed是由于空闲连接数超过了连接池的最大空闲连接数而被强制关闭的连接数
	//
	// 2.核心逻辑: 假设输入参数为n, 先计算出当前时间a减去n之后的时间b,然后遍历连接池(db.freeConn)中的所有空闲连接,
	// 如果连接的创建时间在b之前,也就是说该连接从创建到现在已近经过了大于n的时间,那么会先将该连接记录在一个closing切片中,
	// 然后将记录连接池中连接的切片的最后面一个元素,移动到切片的最前面,最后去除末尾的元素,(相当于去掉首部元素然后把末尾元素移动到首部),
	// 如此循环直到连接池遍历完为止
	sqlDB.SetConnMaxLifetime(10 * time.Second)
}
