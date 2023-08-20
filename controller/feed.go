package controller

import (
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

	latest_time, err := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Invalid timestamp: latest_time"})
		return
	}

	var videos []Video
	result := db.
		Where("upload_time < ?", latest_time).Preload("Author").// 投稿时间过滤
		Order("upload_time desc").             // 按照投稿时间倒序
		Limit(30).                             // 限制返回30条
		Find(&videos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  "Unable to get videos" + err.Error(),
		})
	}
    
	// Use demo video if there is no video in database
	if len(videos) == 0 {
		videos = append(videos, DemoVideo)
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  videos[len(videos)-1].UploadTime, // 返回最旧视频的投稿时间
	})
}
