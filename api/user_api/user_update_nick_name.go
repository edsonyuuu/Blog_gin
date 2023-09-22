package user_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/utils/jwts"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"strings"
)

type UserUpdateNickNameRequest struct {
	NickName string `json:"nick_name" structs:"nick_name"`
	Sign     string `json:"sign"`
	Link     string `json:"link"`
}

// UserUpdateNickName 修改当前登录人的昵称，签名，链接
// @Tags 用户管理
// @Summary 修改当前登录人的昵称，签名，链接
// @Description 修改当前登录人的昵称，签名，链接
// @Param token header string true "token"
// @Param data body UserUpdateNickNameRequest true "昵称，签名，链接"
// @Produce json
// @Success 200 {object} res.Response{}
// @Router /api/user_info [put]
func (UserApi) UserUpdateNickName(c *gin.Context) {
	var cr UserUpdateNickNameRequest

	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}
	var newMaps = map[string]interface{}{}
	maps := structs.Map(cr)
	for k, v := range maps {
		if val, ok := v.(string); ok && strings.TrimSpace(val) != "" {
			newMaps[k] = val
		}
	}

	var user models.UserModel
	err = global.DB.Take(&user, claims.UserID).Error
	if err != nil {
		res.FailWithMessage("用户不存在", c)
		return
	}

	err = global.DB.Model(&user).Updates(newMaps).Error
	if err != nil {
		global.Log.Errorf("修改用户名失败 err:%+v\n", err.Error())
		res.FailWithMessage("修改用户名失败", c)
		return
	}

	res.OkWithMessage("修改个人信息成功", c)
}
