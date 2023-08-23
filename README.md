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



### 使用Redis实现了关注操作，关注列表，粉丝列表

### 思路

思路是将要传递的user结构体转换成字符串数据存储在redis里，需要使用的时候再取出来。

token验证和查询用户都可以使用Redis优化数据

使用Redis的存储思路如图：

![image-20230820113047495](img/REARDONE.png)

**用用户ID可以访问用户的token进行验证，用户的token可以取出用户本体**

**关注列表用用户的ID作为key，被关注的用户ID作为field，value为关注用户的json数据**

**关注列表用用户的token作为key，粉丝用户ID作为field，value为粉丝用户的json数据**

user.go 在登录和注册的时候将用户数据注册进redis的哈希结构中，用 HSetNX（分布式锁）的方式存储，防止重复操作



### 使用FTP服务器对视频文件进行传输
