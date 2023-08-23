package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction 关注和取消关注操作
func RelationAction(c *gin.Context) {
    
	rdb := c.MustGet("rdb").(*redis.Client)

	token := c.Query("token")
	toUserId,err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		panic(err)
		}
	action_type := c.Query("action_type")

	var user User
    //检验令牌
	CheckToken, err := rdb.HGet(context.Background(),"user", token).Result()
	if err != nil {
		panic(err)
		}
	err = json.Unmarshal([]byte(CheckToken),&user)
    if err != nil {
	panic(err)
    }
	//自己不能关注自己
	if toUserId == user.Id{
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "You can't follow yourself"})
		return
		}	

	if token == user.Token{
		//执行关注或取消关注操作
		to_user_id := strconv.FormatInt(toUserId, 10)
		var ToToken string 
		ToToken, err = rdb.HGet(context.Background(),"token", to_user_id).Result()//拿到对方用户token
		if err != nil {
			fmt.Println(toUserId)
			panic(err)
			}
			var ToUser User
		CheckToken, err = rdb.HGet(context.Background(),"user", ToToken).Result() //使用token查user
		if err != nil {
			panic(err)
			}
		err = json.Unmarshal([]byte(CheckToken),&ToUser) //得到对方用户数据
		if err != nil {
			panic(err)
			}
		//数据转换
		idStr := strconv.FormatInt(user.Id, 10)
		ToUser_Id:= strconv.FormatInt(ToUser.Id, 10)
		NewToUser,err := json.Marshal(ToUser) //json
		NewUser,err := json.Marshal(user) //json
		if err != nil {
				panic(err)
		} 
        switch action_type{
		case "1":    //加入关注列表并加入对方的粉丝列表,储存在redis中
			err = rdb.HSetNX(context.Background(),idStr, ToUser_Id, string(NewToUser)).Err()
			if err != nil {
			panic(err)
			 }   
			 //对方的粉丝列表
			 err = rdb.HSetNX(context.Background(),ToUser.Token, idStr,string(NewUser)).Err()
			if err != nil {
			panic(err)
			 }   
	    case"2" :  //从关注列表和对方的粉丝列表中移除
		    err = rdb.HDel(context.Background(),idStr, ToUser_Id).Err()
			if err != nil {
				panic(err)
				 }
				  //对方的粉丝列表
			 err = rdb.HDel(context.Background(),ToUser.Token, idStr).Err()
			 if err != nil {
			 panic(err)
			  }   
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// 所有用户关注列表用redis进行存储
func FollowList(c *gin.Context) {
	var userlist []User
    id := c.Query("user_id")
	rdb := c.MustGet("rdb").(*redis.Client)
	data, err := rdb.HGetAll(context.Background(),id).Result()
	if err != nil {
		panic(err)
	}
	// data是一个map类型，val是string型,这里使用使用循环迭代输出
	//json字符串转成user结构体传入userlist
	if len(data) != 0{  
		for _, val := range data {
			var user User
			err := json.Unmarshal([]byte(val) , &user)
			if err != nil{
				panic(err)
			}
			userlist = append(userlist, user)
		  }
	} 
   
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userlist,
	})

}

// 粉丝列表
func FollowerList(c *gin.Context) {
	var userlist []User
	token := c.Query("token")
	rdb := c.MustGet("rdb").(*redis.Client)
	data, err := rdb.HGetAll(context.Background(),token).Result()
	if err != nil {
		panic(err)
	}
	// data是一个map类型，val是string型,这里使用使用循环迭代输出
	//json字符串转成user结构体传入userlist
	if len(data) != 0{  
		for _, val := range data {
			var user User
			err := json.Unmarshal([]byte(val) , &user)
			if err != nil{
				panic(err)
			}
			userlist = append(userlist, user)
		  }
	} 
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userlist,
	})
}

// 好友列表所有用户都有相同的好友列表
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}
