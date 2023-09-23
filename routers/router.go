package routers

import (
	"Blog_gin/global"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"net/http"
)

type RouterGroup struct {
	*gin.RouterGroup
}

func InitRouter() *gin.Engine {
	//设置模式
	gin.SetMode(global.Config.System.Env)
	router := gin.Default()
	router.StaticFS("uploads", http.Dir("uploads"))
	router.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	//路由群
	apiRouterGroup := router.Group("api")

	routerGroupApp := RouterGroup{apiRouterGroup}
	//系统设置
	routerGroupApp.SettingsRouter()
	//图片管理
	routerGroupApp.ImagesRouter()
	//用户管理
	routerGroupApp.UserRouter()
	//标签管理
	routerGroupApp.TagRouter()
	//广告模块
	routerGroupApp.AdvertRouter()
	//菜单模块
	routerGroupApp.MenuRouter()
	//消息管理
	routerGroupApp.MessageRouter()
	//文章管理
	routerGroupApp.ArticleRouter()
	//文章点赞

	//评论管理

	//新闻管理

	//聊天管理

	//日志管理

	//数据管理

	return router
}
