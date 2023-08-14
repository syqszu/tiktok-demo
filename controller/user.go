package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
	},
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

type UserRegisterResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

func GenerateToken(username string, password string) (string, error) {
	// 使用 bcrypt 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// 生成 token
	token := username + string(hashedPassword)
	return token, err
}

func ReturnError(c *gin.Context, msg string, errCode int32) {
	c.JSON(http.StatusBadRequest, UserRegisterResponse{
		Response: Response{StatusCode: errCode, StatusMsg: msg},
	})
	return
}

// user/register 处理用户注册
func Register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// 从参数中获取用户名和密码
	username := c.Query("username")
	password := c.Query("password")

	// 校验是否已存在该用户名
	var count int64
	db.Where(User{Name: username}).Count(&count)
	if count != 0 {
		ReturnError(c, "用户已存在", 1)
		return
	}

	// 生成用户token
	token, err := GenerateToken(username, password)
	if err != nil {
		ReturnError(c, "注册失败，请重试", 2)
		fmt.Errorf("Generate token error: %v", err)
		return
	}

	// 在数据库中注册新用户
	newUser := User{
		Name:  username,
		Token: token,
	}
	result := db.Create(&newUser)
	if result.Error != nil {
		ReturnError(c, "注册失败，请重试", 2)
		fmt.Errorf("Insert user error: %v", result.Error)
		return
	}
	fmt.Printf("Created user with ID = %d", newUser.Id)

	// 更新用户信息到内存映射
	usersLoginInfo[newUser.Name] = newUser

	// 返回注册成功响应
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   newUser.Id,
		Token:    token,
	})
}

// Login 处理用户登录
func Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

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
	db := c.MustGet("db").(*gorm.DB)

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
