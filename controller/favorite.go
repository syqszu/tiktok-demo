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

	// TODO: 根据user_id判断是否有权限查看
	var user User
	token := c.Query("token")

	// 检查token是否有效
	if err := db.Preload("FavoritedVideos").Where(User{Token: token}).First(&user).Error; err != nil {
		// 无效token
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Invalid token"})
			fmt.Println("Invalid token")
			return
		}
		// 其他错误
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "Internal error"})
		fmt.Println("Internal error")
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
