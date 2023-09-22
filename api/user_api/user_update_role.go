package user_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/models/ctype"
	"Blog_gin/res"
	"github.com/gin-gonic/gin"
)

type UserRole struct {
	Role     ctype.Role `json:"role" binding:"required,oneof=1 2 3 4"  msg:"权限参数错误"`
	NickName string     `json:"nick_name"` //防止用户昵称非法
	UserID   uint       `json:"user_id"   binding:"required"  msg:"用户id错误"`
}

// UserUpdateView  用户权限变更
// @Summary 修改用户权限
// @Description 根据用户ID修改用户的权限和昵称
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body UserRole true "用户权限信息"
// @Success 200 {object}  res.Response{}
// @Router /api/user_role [put]
func (UserApi) UserUpdateView(c *gin.Context) {
	var cr UserRole

	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(err, &cr, c)
		return
	}

	var user models.UserModel

	err := global.DB.Take(&user, cr.UserID).Error
	if err != nil {
		res.FailWithMessage("用户id错误,用户不存在", c)
		return
	}
	err = global.DB.Model(&user).Updates(map[string]any{
		"role":      cr.Role,
		"nick_name": cr.NickName,
	}).Error
	if err != nil {
		res.FailWithMessage("修改权限失败", c)
		return
	}
	res.OkWithMessage("修改权限成功", c)

}
