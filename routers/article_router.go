package routers

import (
	"Blog_gin/api"
	"Blog_gin/middleware"
)

func (router RouterGroup) ArticleRouter() {
	app := api.ApiGroupApp.ArticleApi
	//新建文章
	router.POST("articles", middleware.JwtAuth(), app.ArticleCreateView)
	//查看文章列表
	router.GET("articles", middleware.JwtAuth(), app.ArticleListView)
	//查看文章目录
	router.GET("categorys", app.ArticleCategoryListView)
	//文章标题详情
	router.GET("articles/detail", app.ArticleDetailByTitleView)
	//文章日历列表
	router.GET("articles/calendar", app.ArticleCalendarView)
	//文章标签列表
	router.GET("articles/tags", app.ArticleTagListView)
	//更新文章
	router.PUT("articles", app.ArticleUpdateView)
	//删除文章
	router.DELETE("articles", app.ArticleRemoveView)
	//用户收藏文章
	router.POST("articles/collects", middleware.JwtAuth(), app.ArticleCollectCreateView)
	//用户查看收藏文章列表
	router.GET("articles/collects", middleware.JwtAuth(), app.ArticleCollectListView)
	//文章收藏批量删除
	router.DELETE("articles/collects", middleware.JwtAuth(), app.ArticleCollBatchView)
	//全文搜索
	router.GET("articles/text", app.FullTextSearchView)
	//通过id搜索文章内容
	router.GET("articles/content/:id", app.ArticleContentByIDView)
	//文章详情
	router.GET("articles/:id", app.ArticleDetailView)
}
