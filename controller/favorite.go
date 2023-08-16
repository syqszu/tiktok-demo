package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handles /douyin/favourite/action
func FavoriteAction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	token := c.Query("token")
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Invalid video_id"})
		return
	}

	// 检查token是否有效
	var user User
	if err := db.First(&user, User{Token: token}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Invalid token"})
			fmt.Println("Invalid token")
			return
		}
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "Internal error"})
		fmt.Println("Internal error")
		return
	}

	// 检查video_id是否有效
	var video Video
	if err := db.First(&video, Video{Id: videoId}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "video_id does not exist"})
			fmt.Println("Invalid video_id")
			return
		}
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "Internal error"})
		fmt.Println("Internal error")
		return
	}

	// TODO: 取消点赞

	// 记录点赞
	
	// Video ID为主键，不会记录重复点赞
	if err := db.Model(&user).Association("FavoritedVideos").Append(&video); err != nil {
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "Internal error"})
		fmt.Printf("failed to add video to favorites: %v", err)
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// Handles /douyin/favourite/list
func FavoriteList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// 获取token和user_id参数
	token := c.Query("token")
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Invalid user_id"})
		return
	}
	
	// 检查user_id是否有效
	var user User
	if err := db.Preload("FavoritedVideos").Where(User{Id: userId}).First(&user).Error; err != nil {
		// 无效user_id
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, Response{StatusCode: 2, StatusMsg: "user_id does not exist"})
			fmt.Println("user_id does not exist")
			return
		}
		// 其他错误
		c.JSON(http.StatusInternalServerError, Response{StatusCode: -1, StatusMsg: "Internal error"})
		fmt.Println("Internal error")
		return
	}

	// 检查token是否有效
	if user.Token != token {
		c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Invalid token"})
		fmt.Println("Invalid token")
		return
	}

	// 返回点赞列表
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: user.FavoritedVideos,
	})
}
