package controller

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySQL
var (
	MYSQL_HOST     string = os.Getenv("MYSQL_HOST")
	MYSQL_PORT     string = os.Getenv("MYSQL_PORT")
	MYSQL_USER     string = os.Getenv("MYSQL_USER")
	MYSQL_PASSWORD string = os.Getenv("MYSQL_PASSWORD")
)

// Redis
var (
	REDIS_HOST    string = os.Getenv("REDIS_HOST")
	REDIS_PORT    string = os.Getenv("REDIS_PORT")
	REDISCLI_AUTH string = os.Getenv("REDISCLI_AUTH")
)

// Object stroage
var (
	VIDEO_SERVER_URL string = os.Getenv("VIDEO_SERVER_URL")
)

func InitServices() (*gorm.DB, *redis.Client) {
	// 建立redis连接
	redis_db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
		Password: REDISCLI_AUTH,
		DB:       0,
	})

	// 建立mysql连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/douyindemo?charset=utf8mb4&parseTime=True&loc=Local", MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST+":"+MYSQL_PORT)

	var mysql_db *gorm.DB
	var err error

	for retries := 0; retries < 60; retries++ {
		mysql_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		fmt.Println("Failed to connect to the database. Retrying in 1 second...")
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		panic("Failed to connect to the database after 60 attempts: " + err.Error())
	}

	// 配置连接池
	sqlDB, err := mysql_db.DB()
	if err != nil {
		panic("Failed to get underlying sql.DB from GORM: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(10)                  // 设置空闲连接池中的最大连接数。
	sqlDB.SetMaxOpenConns(100)                 // 设置数据库的最大打开连接数。
	sqlDB.SetConnMaxLifetime(10 * time.Second) // 设置连接可以重复使用的最长时间：10s

	// 自动迁移数据结构
	err = mysql_db.AutoMigrate(&Comment{}, &Video{}, &User{})
	if err != nil {
		panic("Failed to migrate tables" + err.Error())
	}

	return mysql_db, redis_db
}
