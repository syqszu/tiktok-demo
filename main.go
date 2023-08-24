package main

import (
	"github.com/gin-gonic/gin"
	"github.com/syqszu/tiktok-demo/controller"
	"github.com/syqszu/tiktok-demo/service"
)

func main() {
	// 初始化MySQL和Redis连接
	db, rdb := controller.InitServices()
	go service.RunMessageServer()

	// 启动 HTTP 服务
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("rdb", rdb)
		c.Next()
	}) // 注册数据库连接中间件

	controller.InitRouter(r) // 初始化路由

	// TODO: read from config
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
