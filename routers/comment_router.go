package routers

import (
	"Blog_gin/api"
	"Blog_gin/middleware"
)

func (router RouterGroup) CommentRouter() {
	app := api.ApiGroupApp.CommentApi
	//添加评论
	router.POST("comments", middleware.JwtAuth(), app.CommentCreateView)
	//查看评论列表
	router.GET("comments", app.CommentListView)
	//根据id查询评论
	router.GET("comments/:id", app.CommentDigg)
	//删除该id的评论
	router.DELETE("comments/:id", middleware.JwtAuth(), app.CommentRemoveView)
}
