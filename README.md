# douyin-demo





### FFmpeg

视频使用了FFmpeg截取了第一秒第一帧作为封面

publish.go

```go
// FFmpeg命令截图
	cmd := exec.Command("FFmpeg/ffmpeg.exe", "-i", "public/"+finalName, "-ss" , "1" , "-vframes", "1", "img/"+finalName+".jpg")
	err = cmd.Run() //运行
	if err != nil {
		fmt.Println(err)
	}
```

异步操作对视频进行压缩转码，使得视频适合播放：

```go
func transcodeVideo(finalName string) {
     //视频转码压缩
	cmd := exec.Command("ffmpeg", "-i", "public/"+finalName, "-c:v","libx264", "-crf" ,"23" ,"-preset" ,"medium", "-c:a" ,"aac", "-b:a", "128k" ,"public/"+"new"+finalName)
	err := cmd.Run() //运行
	if err != nil {
		fmt.Println(err)
	}
	os.Remove("public/"+finalName) //移除原文件
	err = os.Rename("public/"+"new"+finalName,"public/"+finalName) //重命名转码文件
	if err != nil {
        fmt.Println("Error renaming file:", err)
        return
    }
```

