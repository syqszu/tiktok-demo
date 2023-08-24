# Build
$PROJECT_NAME = "tiktok-demo"
go build

# Set environment variables

# MySQL
$env:MYSQL_HOST = ""
$env:MYSQL_PORT = ""
$env:MYSQL_USER = ""
$env:MYSQL_PASSWORD = ""

# Redis
$env:REDIS_HOST = ""
$env:REDISCLI_AUTH = ""
$env:REDIS_PORT = ""

# Object Storage
$env:VIDEO_SERVER_URL = "http://192.168.3.85:8080/" # 填写HTTP服务端IP地址与端口

# Run tiktok-demo.exe
& ".\$PROJECT_NAME.exe"