package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)





type ChatResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}

// MessageAction 
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")
    action_type := c.Query("action_type")

	rdb := c.MustGet("rdb").(*redis.Client)
    
	if action_type == "1" {
		data, err := rdb.HGet(context.Background(),"user",token).Result() //查询是否存在该用户
	if err == nil { //如果存在该用户
		var user User
        err = json.Unmarshal([]byte(data), &user) //取出用户数据
		if err != nil{
			panic(err)
		}
		User_ID := strconv.FormatInt(user.Id,10)
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(user.Id, int64(userIdB)) //相当于聊天窗口id

	    _,err = rdb.Do(context.Background(),"Select",1).Result() //切换到第二个数据库存储
		if err != nil {
			panic(err)
		   } 
        err := rdb.HSetNX(context.Background(),User_ID , toUserId, 0).Err() //初始化MsgID存储在redis
        if err != nil {
	      panic(err)
         } 
         messageId, err := rdb.HIncrBy(context.Background(),User_ID, toUserId, 1).Result() //messageIdSequence自增1然后传递回来
        if err != nil {
	      panic(err)
        } 
		fmt.Println("此时消息的id为",messageId)
	
		timestamp := time.Now().Unix() // 获取时间的整数表示
		curMessage := Message{
			ID:         messageId,
			ToUserID: int64(userIdB),
			FromUserID: user.Id,
			Content:    content,
			CreateTime: timestamp,
		}
		NewMessage,err := json.Marshal(curMessage)
		if err != nil{
			panic(err)
		}
		MsgID:= strconv.FormatInt(messageId, 10)
        err = rdb.HSetNX(context.Background(),chatKey,MsgID,string(NewMessage)).Err() //将信息存入redis中
		if err !=nil{
			panic(err)
		}
		
		_,err = rdb.Do(context.Background(),"Select",0).Result() //切换到第一个数据库存储
		if err !=nil{
			panic(err)
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		_,err := rdb.Do(context.Background(),"Select",0).Result() //切换到第一个数据库存储
		if err !=nil{
			panic(err)
		}
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
 
	}
	return
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	rdb := c.MustGet("rdb").(*redis.Client)
	token := c.Query("token")
	toUserId := c.Query("to_user_id")

    _,err := rdb.Do(context.Background(),"Select",0).Result() //切换到第一个数据库存储
		if err !=nil{
			panic(err)
		}

	data,err := rdb.HGet(context.Background(),"user",token).Result() //拿出发信息的用户数据
	if err != nil{
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"}) //没有用户数据
		panic(err)
	}
	var user User
    err =json.Unmarshal([]byte(data),&user)
	if err !=nil{
		panic(err)
	}
	userIdB, _ := strconv.Atoi(toUserId)
    _,err = rdb.Do(context.Background(),"Select",1).Result() //切换到第二个数据库存储
	if err !=nil{
		panic(err)
	}
	chatKey := genChatKey(user.Id, int64(userIdB))
    Message_List_string,err:=rdb.HGetAll(context.Background(),chatKey).Result() //返回消息列表
	if err !=nil{
		panic(err)
	}
	var Message_List []Message
	for _,Msg:= range Message_List_string{
		var Msg_One Message
		err =json.Unmarshal([]byte(Msg),&Msg_One)
		if err != nil{
			panic(err)
		}
		Message_List = append(Message_List, Msg_One)
	}
	_,err = rdb.Do(context.Background(),"Select",0).Result() //切换到第一个数据库存储
	if err !=nil{
		panic(err)
	}

	fmt.Println("我得到的列表",Message_List)
    c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0,
		StatusMsg: "获取成功",}, MessageList: Message_List})
}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
