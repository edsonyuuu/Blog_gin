package user_api

import (
	"Blog_gin/global"
	"Blog_gin/res"
	"github.com/gin-gonic/gin"
)

// QQLoginLinkView 获取qq登录的跳转链接
//
//	@Tags			用户管理
//	@Summary		获取qq登录的跳转链接
//	@Description	获取qq登录的跳转链接，data就是qq的跳转地址
//	@Produce		json
//	@Success		200	{object}	res.Response{}
//	@Router			/api/qq_login_path [get]
func (UserApi) QQLoginLinkView(c *gin.Context) {
	path := global.Config.QQ.GetPath()
	if path == "" {
		res.FailWithMessage("未配置qq登录地址", c)
		return
	}

	res.OkWithData(path, c)
	return

}
