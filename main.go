package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syqszu/tiktok-demo/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
}

var db *gorm.DB

func main() {
	// 建立数据库连接
	dsn := "root:123456@tcp(127.0.0.1:3306)/douyindemo?charset=utf8mb4&parseTime=True&loc=Local" // TODO: 从配置文件中读取
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
	fmt.Printf("Connected to the database %s", db)

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get underlying sql.DB from GORM: " + err.Error())
	}
	defer sqlDB.Close()
	sqlDB.SetMaxIdleConns(10)                  // 设置空闲连接池中的最大连接数。
	sqlDB.SetMaxOpenConns(100)                 // 设置数据库的最大打开连接数。
	sqlDB.SetConnMaxLifetime(10 * time.Second) // 设置连接可以重复使用的最长时间：10s

	// 自动迁移数据结构
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("Failed to migrate table: users\n" + err.Error())
	}

	go service.RunMessageServer()

	// 启动 HTTP 服务

	/**
	 * HTTP 响应代码约定：
	 * 成功:200
	 * 输入错误:400
	 * 内部错误：500
	 */

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}) // 注册数据库连接中间件

	initRouter(r) // 初始化路由

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
