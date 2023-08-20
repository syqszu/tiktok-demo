package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syqszu/tiktok-demo/controller"
	"github.com/syqszu/tiktok-demo/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/go-redis/redis/v8"
)


func main() {
	//建立redis连接
    rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// 建立数据库连接
	dsn := controller.DB_USER + ":" + controller.DB_PASSWORD + "@tcp(" + controller.DB_SERVER + ")/douyindemo?charset=utf8mb4&parseTime=True&loc=Local" // TODO: 从配置文件中读取
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
	fmt.Println("Connected to the database")
    
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
	err = db.AutoMigrate(&controller.Comment{})
	if err != nil {
		panic("Failed to migrate table: comments\n" + err.Error())
	}
	err = db.AutoMigrate(&controller.Video{})
	if err != nil {
		panic("Failed to migrate table: videos\n" + err.Error())
	}
	err = db.AutoMigrate(&controller.User{})
	if err != nil {
		panic("Failed to migrate table: users\n" + err.Error())
	}

	go service.RunMessageServer()

	// 启动 HTTP 服务

	/**
	 * HTTP 响应代码约定
	 * 成功：200
	 * 失败：400
	 */

	r := gin.Default()

	// get ip address of this server
    

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("rdb", rdb)
		c.Next()
	}) // 注册数据库连接中间件

	initRouter(r) // 初始化路由

	// TODO: read from config
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
