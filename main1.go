package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"

	
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)

	r.Run() 
	// 建立与数据库的连接
	db, err := sql.Open("mysql", getMySQLConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func getMySQLConnectionString() string {
	// 从提供的数据中获取相关信息
	host := "127.0.0.1"
	port := 3306
	database := "douyin"
	username := "root"
	password := "111111"
	charset := "utf8mb4"
	parseTime := true
	loc := "Local"

	// 根据获取的信息构建连接字符串
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		username, password, host, port, database, charset, parseTime, loc,
	)

	return connStr
}
