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

    var maxID int64
    // 查询最大的 ID
    result := db.Table("videos").Select("MAX(id)").Scan(&maxID)
	if result.Error != nil {
		panic(result.Error)
	}

	ServerIP := "http://45.40.228.46:23333" //视频服务器地址
	
	video := Video{
		Id: maxID+1,
		AuthorID: user.Id,
		Author:   user,                              // assuming usersLoginInfo[token] returns a User object
		PlayUrl:  ServerIP + "/public/" + finalName, // assuming the video can be accessed from /public/ endpoint
		// Fill the other fields as per your requirement
		CoverUrl: "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg", //之后使用Ffpemg对视频切片获取封面
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
