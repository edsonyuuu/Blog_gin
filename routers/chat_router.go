package routers

import "Blog_gin/api"

func (router RouterGroup) ChatRouter() {
	app := api.ApiGroupApp.ChatApi
	//群聊
	router.GET("chat_groups", app.ChatGroupView)
	//群聊记录
	router.GET("chat_groups_records", app.ChatListView)
}
