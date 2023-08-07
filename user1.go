package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 创建一个映射，存储用户登录信息，使用用户名作为键，User 结构体作为值
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

// 记录用户 ID 序列的变量
var userIdSequence = int64(1)
var db *sql.DB

// UserLoginResponse 定义用户登录响应的结构体
type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

// UserResponse 定义用户响应的结构体
type UserResponse struct {
	Response
	User User `json:"user"`
}

// Register 处理用户注册的函数
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	var count int
	// 查询数据库中是否存在相同的 token
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE token = ?", token).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "内部服务器错误"},
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户已存在"},
		})
		return
	}

	// 在数据库中插入新用户记录
	result, err := db.Exec("INSERT INTO users (name, token) VALUES (?, ?)", username, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "内部服务器错误"},
		})
		return
	}

	lastInsertID, _ := result.LastInsertId()

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   lastInsertID,
		Token:    token,
	})
}

// Login 处理用户登录的函数
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	var id int64
	// 查询数据库获取用户 ID
	err := db.QueryRow("SELECT id FROM users WHERE token = ?", token).Scan(&id)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   id,
		Token:    token,
	})
}

// UserInfo 获取用户信息的函数
func UserInfo(c *gin.Context) {
	token := c.Query("token")

	var user User
	// 查询数据库获取用户信息
	err := db.QueryRow("SELECT id, name FROM users WHERE token = ?", token).Scan(&user.Id, &user.Name)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User:     user,
	})
}
