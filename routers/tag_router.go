package routers

import "Blog_gin/api"

func (router RouterGroup) TagRouter() {
	app := api.ApiGroupApp.TagApi
	//创建标签
	router.POST("tags", app.TagCreateView)
	//标签列表
	router.GET("tags", app.TagListView)
	//标签名称列表
	router.GET("tag_names", app.TagNameListView)
	//更新标签
	router.PUT("tags/:id", app.TagUpdateView)
	//删除标签
	router.DELETE("tags", app.TagRemoveView)

}
