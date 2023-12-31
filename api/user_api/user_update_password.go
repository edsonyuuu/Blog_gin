package user_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/utils/jwts"
	"Blog_gin/utils/pwd"
	"github.com/gin-gonic/gin"
)

type UpdatePasswordRequest struct {
	OldPwd string `json:"old_pwd"  binding:"required" msg:"请输入旧密码"`
	Pwd    string `json:"pwd" binding:"required"  msg:"请输入新密码"`
}

// UserUpdatePassword 更新密码
// @Summary 更新用户密码
// @Description 更新用户的密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param request body UpdatePasswordRequest true "请求体参数"
// @Success 200 {object} res.Response{}
// @Router /api/user_password [put]
func (UserApi) UserUpdatePassword(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	var cr UpdatePasswordRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(err, &cr, c)
		return
	}

	var user models.UserModel

	err := global.DB.Take(&user, claims.UserID).Error
	if err != nil {
		res.FailWithMessage("用户不存在", c)
		return
	}
	//判断密码是否一致
	if !pwd.CheckPwd(user.Password, cr.OldPwd) {
		res.FailWithMessage("密码错误", c)
		return
	}
	hashPwd := pwd.HashPwd(cr.Pwd)
	err = global.DB.Take(&user).Update("password", hashPwd).Error
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("新密码修改失败", c)
		return
	}
	res.OkWithMessage("新密码修改成功", c)
	return

}
