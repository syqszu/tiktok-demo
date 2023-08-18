package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		Author:   user,
		PlayUrl:  VIDEO_SERVER_URL + "public/" + finalName, // 视频作为静态资源通过 /public/ 访问
		// Fill the other fields as per your requirement
		CoverUrl: "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg", // TODO: 使用Ffpemg对视频切片获取封面
	}

	if err := db.Create(&video).Error; err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	VideoList = append(VideoList, video)
	//在结构体中追加元素

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	token := c.Query("token")
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID is required",
		})
		return
	}

	// 验证token有效性
	_, err := ValidateToken(c, db, token)
	if err != nil {
		return
	}

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
