package controller

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/syqszu/tiktok-demo/Loginpb"

	"gorm.io/driver/mysql" //加入mysql
	"gorm.io/gorm"         //加入grom
	"gorm.io/gorm/schema"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{}

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

// 注册函数
func Register(c *gin.Context) {

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

	username := c.Query("username")
	password := c.Query("password")

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

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else { //如果不存在将其加入数据库，同时放入临时结构体

		//这部步是更新id
		var result Loginpb.User
		db.Model(&Loginpb.User{}).Order("id desc").First(&result) //
		//Model是你要查询的模型名称，id是模型中的ID字段名。使用Order("id desc")可以按照ID字段的降序排序，然后使用First方法取得第一条结果，即具有最大ID值的数据。
		if result.Id != nil {
			userIdSequence = *result.Id //如果这条不是第一条数据
		}
		atomic.AddInt64(&userIdSequence, 1) //更新id

		//请求的注册数据
		request := Loginpb.DouyinUserLoginRequest{
			Username: &username,
			Password: &password,
		}
		//需要返回的注册数据
		statusMsg := "注册成功"
		var statusCode int32
		statusCode = 1
		response := Loginpb.DouyinUserLoginResponse{
			UserId:     &userIdSequence,
			Token:      &token,
			StatusCode: &statusCode,
			StatusMsg:  &statusMsg,
		}
		users := Loginpb.User{
			Id:   &userIdSequence,
			Name: &username,
		}
		db.Create(&request)  //将注册信息传进数据库
		db.Create(&response) //将注册响应信息传进数据库
		db.Create(&users)    //user
		//
		//结构体操作

		newUser := User{
			Id:   userIdSequence,
			Name: username,
		}
		usersLoginInfo[token] = newUser
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    username + password,
		})
	}
}

// 登录操作，先在数据库搜索是否存在
func Login(c *gin.Context) {
	//打开数据库
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

	username := c.Query("username")
	password := c.Query("password")

	token := username + password

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
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {

	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
