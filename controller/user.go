package controller

import (
	"database/sql"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(0) //变成从0开始

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

var (
	db *sql.DB // MySQL数据库连接
)

// InitializeDB 初始化数据库连接
func InitializeDB() {
	var err error

	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/douyin")
	if err != nil {
		panic(err)
	}
}

// Register 处理用户注册
func Register(c *gin.Context) {
	// 从查询参数中获取用户名和密码
	username := c.Query("username")
	password := c.Query("password")

	// 生成用户token
	token := username + password
	//先判断数据库中是否存在该数据
	var List Loginpb.DouyinUserLoginResponse
	db.Where("token = ? ", token).Find(&List)
	if List.Token != nil { //如果数据库中存在该数据，将数据传入结构体
		if _, exist := usersLoginInfo[token]; !exist { //如果结构体中不存在该用户
			var dataUsers Loginpb.User
			db.Where("id = ? ", token).Find(&dataUsers)
			NewUser := User{
				Id:   *List.UserId,
				Name: *dataUsers.Name,
			}
			usersLoginInfo[*List.Token] = NewUser //将数据传入结构体
		}
	}

	// 检查用户是否已存在
	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户已存在"},
		})
		return
	}

	// 使用 bcrypt 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 2, StatusMsg: "注册失败，请重试"},
		})
		return
	}

	// 在数据库中插入用户信息，同时存储哈希后的密码
	result, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 2, StatusMsg: "注册失败，请重试"},
		})
		return
	}

	// 获取插入的用户ID
	userID, _ := result.LastInsertId()

	// 更新用户信息到内存映射
	atomic.AddInt64(&userIdSequence, 1)
	newUser := User{
		Id:   userID,
		Name: username,
	}
	usersLoginInfo[token] = newUser

	// 返回注册成功响应
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   userID,
		Token:    token,
	})
}

// Login 处理用户登录
func Login(c *gin.Context) {
	// 从查询参数中获取用户名和密码
	username := c.Query("username")
	password := c.Query("password")

	// 生成用户token
	token := username + password

	// 在数据库中查询用户
	var userID int64
	err := db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", username, password).Scan(&userID)
	if err != nil {
		var dataList Loginpb.DouyinUserLoginResponse
		db.Where("token = ?", token).Find(&dataList)

		if dataList.Token != nil { //如果数据库中存在该用户，如果临时结构体不存在该用户的话，将该用户添加进结构体
			var dataUser Loginpb.User
			db.Where("id = ?", *dataList.UserId).Find(&dataUser)
			if _, exist := usersLoginInfo[token]; !exist { //如果结构体中不存在该用户
				usersLoginInfo[token] = User{
					Id:   *dataList.UserId,
					Name: *dataUser.Name,
				}
				c.JSON(http.StatusOK, UserLoginResponse{ //返回用户id和Token
					Response: Response{StatusCode: 0},
					UserId:   *dataList.UserId,
					Token:    *dataList.Token,
				})
			}
		}
		//结构体操作
		if user, exist := usersLoginInfo[token]; exist {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "用户不存在或密码错误"},
			})
			return
		}

		// 更新用户信息到内存映射
		usersLoginInfo[token] = User{
			Id: userID,
		}

		// 返回登录成功响应
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userID,
			Token:    token,
		})
	}
}

// UserInfo 获取用户信息
func UserInfo(c *gin.Context) {
	// 从查询参数中获取用户token
	token := c.Query("token")

	// 从内存映射中获取用户信息
	user, exist := usersLoginInfo[token]
	if !exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
		return
	}

	// 构建响应消息
	response := UserResponse{
		Response: Response{StatusCode: 0},
		User: User{
			Id:              user.Id,
			Name:            user.Name,
			FollowCount:     user.FollowCount,
			FollowerCount:   user.FollowerCount,
			IsFollow:        user.IsFollow,
			Avatar:          user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature:       user.Signature,
			TotalFavorited:  user.TotalFavorited,
			WorkCount:       user.WorkCount,
			FavoriteCount:   user.FavoriteCount,
		},
	}

	// 返回响应
	c.JSON(http.StatusOK, response)
}
