package routers

import (
	"Blog_gin/api"
	"Blog_gin/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

var store = cookie.NewStore([]byte("HyvCD89g3VDJ9646BFGEh37GFJ"))

func (router RouterGroup) UserRouter() {
	app := api.ApiGroupApp.UserApi
	//使用session中间件
	router.Use(sessions.Sessions("sessionid", store))
	//邮箱登录
	router.POST("email_login", app.EmailLoginView)
	//QQ登录
	router.POST("login", app.QQLoginView)
	//获取跳转QQ登录的跳转链接
	router.GET("qq_login_path", app.QQLoginLinkView)
	//用户权限的变更
	router.PUT("user_role", middleware.JwtAdmin(), app.UserUpdateView)
	//用户修改密码
	router.PUT("user_password", middleware.JwtAuth(), app.UserUpdatePassword)
	//用户退出登录
	router.POST("logout", middleware.JwtAuth(), app.LogoutView)
	//管理员删除用户
	router.DELETE("users", middleware.JwtAdmin(), app.UserRemoveView)
	//用户绑定邮箱
	router.POST("user_bind_email", middleware.JwtAuth(), app.UserBindEmailView)
	//用户创建
	router.POST("users", app.UserCreateView)
	//用户列表
	router.GET("users", middleware.JwtAuth(), app.UserListView)
	//用户信息
	router.GET("user_info", middleware.JwtAuth(), app.UserInfoView)
	//用户修改信息
	router.PUT("user_info", middleware.JwtAuth(), app.UserUpdateNickName)

}
