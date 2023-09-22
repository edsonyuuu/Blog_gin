package routers

import "Blog_gin/api"

func (router RouterGroup) ImagesRouter() {
	app := api.ApiGroupApp.ImagesApi
	//多图片上传
	router.POST("images", app.ImageUploadView)
	//单图片上传
	router.POST("image", app.ImageUploadDataView)
	//查看图片，图片进行了分页
	router.GET("images", app.ImageListView)
	//查看图片名称列表
	router.GET("images_names", app.ImageNameListView)
	//删除图片
	router.DELETE("images", app.ImageRemoveView)
	//更改图片
	router.PUT("images", app.ImageUpdateView)

}
