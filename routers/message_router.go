package routers

import (
	"Blog_gin/api"
	"Blog_gin/middleware"
)

func (router RouterGroup) MessageRouter() {
	app := api.ApiGroupApp.MessageApi
	//消息创建
	router.POST("messages", middleware.JwtAuth(), app.MessageCreateView)
	//查看所有消息列表
	router.GET("messages_all", middleware.JwtAuth(), app.MessageListAllView)
	//用户与他人的消息列表
	router.GET("messages", middleware.JwtAuth(), app.MessageListView)
	//消息记录
	router.GET("messages_record", middleware.JwtAuth(), app.MessageRecordView)

}
