package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	token := c.PostForm("token")

	// 验证token有效性
	user, err := ValidateToken(c, db, token)
	if err != nil {
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//数据入库
	video := Video{
		AuthorID: user.Id,
		Author:   user,                   // assuming usersLoginInfo[token] returns a User object
		PlayUrl:  "/public/" + finalName, // assuming the video can be accessed from /public/ endpoint
		// Fill the other fields as per your requirement
	}

	if err := db.Create(&video).Error; err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	token := c.PostForm("token")
	userID := c.PostForm("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID is required",
		})
		return
	}
	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	db := c.MustGet("db").(*gorm.DB)

	var videos []Video
	if err := db.Where("author_id = ?", userID).Preload("Author").Find(&videos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
