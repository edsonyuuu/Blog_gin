package routers

import "Blog_gin/api"

func (router RouterGroup) MenuRouter() {
	app := api.ApiGroupApp.MenuApi
	//创建菜单
	router.POST("menus", app.MenuCreateView)
	//查看菜单
	router.GET("menus", app.MenuListView)
	//通过菜单名字来查看
	router.GET("menu_names", app.MenuNameList)
	//更新菜单
	router.PUT("menus/:id", app.MenuUpdateView)
	//删除菜单
	router.DELETE("menus", app.MenuRemoveView)
	//详细菜单信息
	router.GET("menus/detail", app.MenuDetailByPathView)
	//详细菜单信息通过id
	router.GET("menus/:id", app.MenuDetailView)
}
