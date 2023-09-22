package user_api

import (
	"Blog_gin/global"
	"Blog_gin/res"
	"Blog_gin/service"
	"Blog_gin/utils/jwts"
	"github.com/gin-gonic/gin"
)

// LogoutView 处理注销功能。
// @Summary 注销用户
// @Description 根据提供的用户信息和令牌注销用户。
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param token header string true "访问令牌"
// @Success 200 {object} res.Response{}
// @Router /api/logout [post]
func (UserApi) LogoutView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)
	token := c.Request.Header.Get("token")
	err := service.ServiceApp.UserService.Logout(claims, token)

	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("注销失败", c)
		return
	}
	res.OkWithMessage("注销成功", c)
}
