package controller_test

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/syqszu/tiktok-demo/controller"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetUpGormDb() (
	db *gorm.DB,
	mockDb *sql.DB,
	mock sqlmock.Sqlmock) {

	var err error

	// 创建 sqlmock 对象
	mockDb, mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("Failed to create sqlmock: %s", err)
		return
	}

	// GORM 在开始时调用 sql.DB.QueryRow("SELECT VERSION()") 来检查数据库连接
	mock.ExpectQuery("^SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.32"))

	// 创建 GORM 对象
	db, err = gorm.Open(mysql.New(mysql.Config{
		Conn: mockDb,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to create gorm database: %s", err)
		return
	}
	return
}

func SetUpGinServer(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("rdb", rdb)
		c.Next()
	}) // 依赖项
	controller.InitRouter(r)
	return r
}

func TestRegister(t *testing.T) {
	// 初始化
	db, mockDb, _ := SetUpGormDb()
	rdb, _ := redismock.NewClientMock()
	r := SetUpGinServer(db, rdb)
	defer mockDb.Close()
	defer rdb.Close()
	
	// 构建测试数据
	

	// 发送请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/douyin/user/register", nil)
	r.ServeHTTP(w, req)
}
