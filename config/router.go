package config

import (
	"feeddemo/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	//不用拦截的接口组
	apiRouter := r.Group("/feed")

	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.RelationFollowList)
	apiRouter.GET("/relation/follower/list/", controller.RelationFollowerList)
	apiRouter.POST("/blog/publish/", controller.PublishBlog)
	apiRouter.DELETE("/blog/delete/", controller.DeleteBlog)
	apiRouter.GET("/blog/getNewList/", controller.GetNewFeeds)
	apiRouter.GET("/blog/getHotList/", controller.GetHotFeeds)
	apiRouter.GET("/blog/getUserBlogList/", controller.GetUserPublishList)
	apiRouter.GET("/blog/getRelationBlogList/", controller.GetRelationFeedList)
	apiRouter.GET("/blog/getRelationBlogList1/", controller.GetRelationFeedList1)
	apiRouter.PUT("/blog/topBlog/", controller.TopBlog)
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)
}
