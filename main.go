package main

import (
	"github.com/gin-gonic/gin"
	"github.com/syqszu/tiktok-demo/controller"
	"github.com/syqszu/tiktok-demo/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
}

func main() {

	controller.InitializeDB()
	dsn := "root:123456@tcp(127.0.0.1:3306)/douyin"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// 自动创建表格
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("Failed to create tables")
	}

	go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
