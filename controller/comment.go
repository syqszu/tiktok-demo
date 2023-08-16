package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	db := c.MustGet("db").(*gorm.DB)
	videoID, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)

	if user, exist := usersLoginInfo[token]; exist {
		switch actionType {
		case "1":
			text := c.Query("comment_text")
			if text == "" {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Comment text cannot be empty"})
				return
			}
			comment := Comment{
				UserID:     user.Id,
				Content:    text,
				VideoID:    int64(videoID),
				CreateDate: time.Now().Format("01-02"),
			}
			if err := db.Create(&comment).Error; err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
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
			c.JSON(http.StatusOK, Response{StatusCode: 0})
		default:
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Invalid action type"})
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	token := c.Query("token")
	videoID := c.Query("video_id")
	db := c.MustGet("db").(*gorm.DB)

	// Verify user
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	var comments []Comment
	if err := db.Where("video_id = ?", videoID).Preload("User").Find(&comments).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}
