package routers

import "Blog_gin/api"

func (router RouterGroup) DiggRouter() {
	app := api.ApiGroupApp.DiggApi
	router.POST("digg/article", app.DiggArticleView)
}
