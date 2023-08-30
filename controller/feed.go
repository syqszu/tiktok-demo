package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// 可选参数 latest_time
	var latest_time int64
	latest_time_str := c.Query("latest_time")
	if latest_time_str == "" {
		latest_time = -1
	} else {
		var err error
		latest_time, err = strconv.ParseInt(c.Query("latest_time"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Invalid timestamp: latest_time"})
			return
		}
	}

	var videos []Video
	var result *gorm.DB
	if latest_time == -1 {
		// 如果没有传入 latest_time 参数，则返回最新的30条视频
		result = db.
			Preload("Author").
			Order("upload_time desc"). // 按照投稿时间倒序
			Limit(30).                 // 限制返回30条
			Find(&videos)
	} else {
		// 如果传入了 latest_time 参数，则返回 latest_time 之前的最新的30条视频
		result = db.
			Where("upload_time < ?", latest_time).Preload("Author"). // 投稿时间过滤
			Order("upload_time desc").                               // 按照投稿时间倒序
			Limit(30).                                               // 限制返回30条
			Find(&videos)
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  "Unable to get videos" + result.Error.Error(),
		})
	}

	// 尝试获取用户信息，确认获取的视频是否点过赞
	token := c.Query("token")
	var user User
	if err := db.Preload("FavoritedVideos").First(&user, User{Token: token}).Error; err == nil {
		// 用户已登录
		fmt.Println(len(user.FavoritedVideos))
		for i, v := range videos {
			for _, f := range user.FavoritedVideos {
				if v.Id == f.Id {
					videos[i].IsFavorite = true
				}
			}
		}
	}

	nextTime := int64(0)
	if len(videos) > 0 {
		nextTime = videos[len(videos)-1].UploadTime
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  nextTime, // 返回最旧视频的投稿时间
	})
}
