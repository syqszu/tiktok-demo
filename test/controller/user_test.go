package controller_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
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
	db, mockDb, mock := SetUpGormDb()
	rdb, _ := redismock.NewClientMock()
	r := SetUpGinServer(db, rdb)
	defer mockDb.Close()
	defer rdb.Close()

	// 构建测试数据
	username := "TestName"
	password := "TestPassword"

	// 设置 mock 数据库操作
	mock.ExpectBegin()

	mock.ExpectQuery("SELECT (.+) FROM `users` WHERE (.+) FOR UPDATE").
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"Id", "Name", "Token"}))

	mock.ExpectExec("INSERT INTO `users` (.+) VALUES (.+)").
		WithArgs(username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// 发送请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/douyin/user/register/?username=%s&password=%s", username, password), nil)
	r.ServeHTTP(w, req)

	// 校验响应
	assert.Equal(t, http.StatusOK, w.Code)

	var resp controller.UserLoginResponse // 解析响应JSON
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %s", err)
	}

	assert.Equal(t, int64(1), resp.UserId) // id = 1
	assert.NotEmpty(t, resp.Token)
}
