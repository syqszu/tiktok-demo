package main

import (
	"fmt"
	"time"

	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/Loginpb"
	"github.com/gin-gonic/gin"

	"gorm.io/driver/mysql" //加入mysql
	"gorm.io/gorm"         //加入grom
	"gorm.io/gorm/schema"
	//    "gorm.io/driver/sqlite"
)

func main() {
	go service.RunMessageServer()

	//连接到数据库,用户名为root，密码为空，连接到的数据库名是douyindemo
	dsn := "root:123456@tcp(127.0.0.1:3306)/douyindemo?charset=utf8mb4&parseTime=True&loc=Local" //一般数据库的连接都是127.0.0.1:3306,不用改
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用此选项后，“User”的表将为“Users”
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(db)

	sqlDB, err := db.DB()

	// SetMaxIdleConns设置空闲连接池中的最大连接数。
	sqlDB.SetMaxIdleConns(10)

	//SetMaxOpenConns设置数据库的最大打开连接数。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime设置连接可以重复使用的最长时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second) //十秒钟

	//1.主键没有    结构体添加gorm.Model
	//2.名称变成复数问题
	//AutoMigrate 会创建表、缺失的外键、约束、列和索引
	//登录注册信息
	db.AutoMigrate(&Loginpb.DouyinUserLoginRequest{})     //创建登录请求表
	db.AutoMigrate(&Loginpb.DouyinUserLoginResponse{})    //创建登录响应表
	db.AutoMigrate(&Loginpb.User{}) //创建用户信息表

	//创建接口
	r := gin.Default()

	/*代码约定

	成功:200
	错误:400

	*/

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
