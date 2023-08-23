package controller

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")
	r.Static("/img", "./img")
	// apis
	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", Feed)
	apiRouter.GET("/user/", UserInfo)
	apiRouter.POST("/user/register/", Register)
	apiRouter.POST("/user/login/", Login)
	apiRouter.POST("/publish/action/", Publish)
	apiRouter.GET("/publish/list/", PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", FavoriteAction)
	apiRouter.GET("/favorite/list/", FavoriteList)
	apiRouter.POST("/comment/action/", CommentAction)
	apiRouter.GET("/comment/list/", CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", RelationAction)
	apiRouter.GET("/relation/follow/list/", FollowList)
	apiRouter.GET("/relation/follower/list/", FollowerList)
	apiRouter.GET("/relation/friend/list/", FriendList)
	apiRouter.GET("/message/chat/", MessageChat)
	apiRouter.POST("/message/action/", MessageAction)
}
