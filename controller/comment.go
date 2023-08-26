package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// Handles /douyin/comment/action
func CommentAction(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	token := c.Query("token")
	actionType := c.Query("action_type")
	videoID, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)

	// 验证token有效性
	user, err := ValidateToken(c, db, token)
	if err != nil {
		return
	}

	switch actionType {
	case "1":
		text := c.Query("comment_text")
		if text == "" {
			c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Comment text cannot be empty"})
			return
		}
		comment := Comment{
			UserID:     user.Id,
			Content:    text,
			VideoID:    int64(videoID),
			CreateDate: time.Now().Format("01-02"),
		}
		if err := db.Create(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		var video Video
		if err := db.First(&video, Video{Id: videoID}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{StatusCode: -1, StatusMsg: "Internal error"})
			fmt.Printf("comment failed to find video 56: %v", err)
			return
		}
		//在视频中记录
		video.CommentCount = video.CommentCount + 1
		if err := db.Model(&Video{}).Where("Id = ?", video.Id).Updates(&video).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{StatusCode: -1, StatusMsg: "Internal error"})
			fmt.Printf("comment failed to update video : %v", err)
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0},
			Comment:  comment,
		})
	case "2":
		commentID := c.Query("comment_id")
		var comment Comment
		if err := db.Where("id = ?", commentID).Delete(&comment).Error; err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		var video Video
		if err := db.First(&video, Video{Id: videoID}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{StatusCode: -1, StatusMsg: "Internal error"})
			fmt.Printf("comment failed to find video 80: %v", err)
			return
		}
		//在视频中记录
		video.CommentCount = video.CommentCount - 1
		if err := db.Model(&Video{}).Where("Id = ?", video.Id).Updates(&video).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{StatusCode: -1, StatusMsg: "Internal error"})
			fmt.Printf("comment failed to update video: %v", err)
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	default:
		c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Invalid action type"})
	}

}

// Handles /douyin/comment/list
func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoID := c.Query("video_id")
	db := c.MustGet("db").(*gorm.DB)

	// 验证token有效性
	_, err := ValidateToken(c, db, token)
	if err != nil {
		return
	}

	// 从数据库读取评论，并按发布时间倒序排列
	var comments []Comment
	if err := db.Where("video_id = ?", videoID).Preload("User").Order("create_date DESC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}
