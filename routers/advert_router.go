package routers

import "Blog_gin/api"

func (router RouterGroup) AdvertRouter() {
	app := api.ApiGroupApp.AdvertApi
	//添加广告
	router.POST("adverts", app.AdvertCreateView)
	//查看广告列表
	router.GET("adverts", app.AdvertListView)
	//更新广告列表
	router.PUT("adverts/:id", app.AdvertUpdateView)
	//删除广告
	router.DELETE("adverts", app.AdvertRemoveView)
}
