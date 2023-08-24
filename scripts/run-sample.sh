#!/bin/bash

# Build
PROJECT_NAME="tiktok-demo"
go build

# Set environment variables

# MySQL
export MYSQL_HOST=""
export MYSQL_PORT=""
export MYSQL_USER=""
export MYSQL_PASSWORD=""

# Redis
export REDIS_HOST=""
export REDISCLI_AUTH=""
export REDIS_PORT=""

# Object Storage
export VIDEO_SERVER_URL="" # 填写HTTP服务端IP地址与端口

# Run tiktok-demo
./$PROJECT_NAME
