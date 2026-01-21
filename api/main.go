package main

import (
	"fmt"
	"log"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/tianyuliang/standardization-backend/api/internal/config"
	"github.com/tianyuliang/standardization-backend/api/internal/handler"
	"github.com/tianyuliang/standardization-backend/api/internal/svc"
)

var configFile = fmt.Sprintf("api/etc/%s.yaml", "api")

func main() {
	// 1. 加载配置
	var c config.Config
	conf.MustLoad(configFile, &c)

	// 2. 初始化数据库连接
	db, err := initDB(&c)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	c.DB = db

	// 3. 初始化服务上下文
	svcCtx := svc.NewServiceContext(c)

	// 4. 创建 HTTP 服务器
	server := rest.MustNewServer(c.RestConf)

	// 5. 注册路由
	handler.RegisterHandlers(server, svcCtx)

	// 6. 启动服务
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// initDB 初始化数据库连接
func initDB(c *config.Config) (*gorm.DB, error) {
	// 从配置中读取数据库连接信息
	// 注意：实际项目中应该从环境变量或配置文件读取
	// 这里使用硬编码作为示例，需要替换为实际配置
	dsn := "root:password@tcp(localhost:3306)/idrm?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// TODO: 从配置文件读取连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}
