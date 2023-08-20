package controller

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// FFmpeg转码操作...
func transcodeVideo(finalName string) {
	//视频转码压缩
	cmd := exec.Command("FFmpeg/ffmpeg.exe", "-i", "public/"+finalName, "-c:v", "libx264", "-crf", "23", "-preset", "medium", "-c:a", "aac", "-b:a", "128k", "public/"+"new"+finalName)
	err := cmd.Run() //运行
	if err != nil {
		fmt.Println(err)
	}
	os.Remove("public/" + finalName)                                //移除原文件
	err = os.Rename("public/"+"new"+finalName, "public/"+finalName) //重命名转码文件
	if err != nil {
		fmt.Println("Error renaming file:", err)
		return
	}

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
	// FFmpeg命令截图
	cmd := exec.Command("FFmpeg/ffmpeg.exe", "-i", "public/"+finalName, "-ss", "1", "-vframes", "1", "img/"+finalName+".jpg")
	err = cmd.Run() //运行
	if err != nil {
		fmt.Println(err)
	}
	//视频转码压缩
	go transcodeVideo(finalName) //异步操作防止程序返回错误

	// 数据入库
	video := Video{
		AuthorID: user.Id,
		Author:   user,
		PlayUrl:  VIDEO_SERVER_URL + "static/" + finalName, // 视频作为静态资源通过 /static/ 访问
		// Fill the other fields as per your requirement
		CoverUrl:   VIDEO_SERVER_URL + "img/" + finalName + ".jpg", // TODO: 使用Ffpemg对视频切片获取封面
		UploadTime: time.Now().Unix(),
	}

	if err := db.Create(&video).Error; err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 返回成功信息
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
