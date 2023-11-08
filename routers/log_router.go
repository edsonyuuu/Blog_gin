package routers

import (
	"Blog_gin/api"
	"Blog_gin/middleware"
)

func (router RouterGroup) LogRouter() {
	app := api.ApiGroupApp.LogApi
	router.GET("logs", app.LogListView)
	router.DELETE("logs", middleware.JwtAdmin(), app.LogRemoveListView)
}
