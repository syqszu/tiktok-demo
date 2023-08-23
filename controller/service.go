package controller

//Mysql,Redis,FTP,服务地址设置

//视频url
var VIDEO_SERVER_URL string = "http://192.168.31.158:23333/" // 填写本机IP地址与服务器端口

//Mysql
var (
	DB_SERVER   string = "127.0.0.1:3306"
	DB_USER     string = "root"
	DB_PASSWORD string = ""
)

//Redis
var (
	RDB_Addr     string = "localhost:6379"
	RDB_PASSWORD string = "" // no password set
	RDB_DB       int    = 0  // use default DB
)

//FTP文件传输服务器
var (
	FtpAddress  string = "192.168.31.158:21"
	FtpUsername string = "23563"
	FtpPassword string = "z1o2m3b4i5e6"
)
