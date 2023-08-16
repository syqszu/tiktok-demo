package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var usersLoginInfo = map[string]User{

}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User     User `json:"user"`
	IsFollow bool `json:"is_follow,omitempty"`
}

type UserRegisterResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

func ReturnError(c *gin.Context, msg string, errCode int32) {
	c.JSON(http.StatusBadRequest, UserRegisterResponse{
		Response: Response{StatusCode: errCode, StatusMsg: msg},
	})
}

// user/register 处理用户注册
func Register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// 从参数中获取用户名和密码
	username := c.Query("username")
	password := c.Query("password")

	// 生成用户token
	tokenBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ReturnError(c, "注册失败，请重试", 2)
		fmt.Printf("Generate token error: %v", err)
		return
	}
	token := string(tokenBytes)

	// 尝试在数据库中注册新用户
	
	newUser := User{
		Name:  username,
		Token: token,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// 校验是否已存在该用户名
		var users []User
    	result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&users, User{Name: username})
		if result.Error != nil {
			fmt.Printf("Read user error: %v", result.Error)
			return result.Error
		}
		if len(users) > 0 {
			fmt.Println("User already exists")
			return errors.New("User already exists")
		}

		// 在数据库中注册新用户
		result = tx.Create(&newUser)
		if result.Error != nil {
			fmt.Printf("Insert user error: %v", result.Error)
			return result.Error
		}
		fmt.Printf("Created user with ID = %d", newUser.Id)

		return nil
	})
	if err != nil {
		ReturnError(c, err.Error(), 3)
		return
	}

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

	// 校验用户信息

	// 校验用户名
	if username == "" {
		ReturnError(c, "用户名不能为空", 1)
		return
	}

	var user User
	result := db.Where(User{Name: username}).First(&user)
	if result.Error != nil {
		ReturnError(c, "用户不存在", 2)
		return
	}

	// 校验密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Token), []byte(password))
	if err != nil {
		ReturnError(c, "密码错误", 3)
		return
	}

	// 返回登录成功响应
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    user.Token,
	})
}

// UserInfo 获取用户信息
func UserInfo(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// 获取用户ID
	userId := c.Query("user_id")

	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		ReturnError(c, "无效用户ID", 1)
		return
	}

	// 查询用户信息
	var user User
	result := db.Where(User{Id: id}).First(&user)
	if result.Error != nil {
		ReturnError(c, "用户不存在", 2)
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User: User{
			Id:              user.Id,
			Name:            user.Name,
			FollowCount:     user.FollowCount,
			FollowerCount:   user.FollowerCount,
			Avatar:          user.Avatar,
			BackgroundImage: user.BackgroundImage,
			Signature:       user.Signature,
			TotalFavorited:  user.TotalFavorited,
			WorkCount:       user.WorkCount,
			FavoriteCount:   user.FavoriteCount,
		},
		IsFollow: true,
	})
}
