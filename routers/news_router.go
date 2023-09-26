package routers

import "Blog_gin/api"

func (router RouterGroup) NewsRouter() {
	app := api.ApiGroupApp.NewsApi
	router.POST("news", app.NewsListView)
}
