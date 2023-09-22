package user_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/models/ctype"
	"Blog_gin/plugins/qq"
	"Blog_gin/res"
	"Blog_gin/utils"
	"Blog_gin/utils/jwts"
	"Blog_gin/utils/pwd"
	"Blog_gin/utils/random"
	"github.com/gin-gonic/gin"
)

func (UserApi) QQLoginView(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		res.FailWithMessage("没有code", c)
		return
	}

	qqInfo, err := qq.NewQQLogin(code)
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}

	openID := qqInfo.OpenID
	//根据openID判断用户是否存在

	var user models.UserModel
	err = global.DB.Take(&user, "token = ?", openID).Error
	ip, addr := utils.GetAddrByGin(c)
	if err != nil {
		//不存在就注册
		hashPwd := pwd.HashPwd(random.RandString(16))
		user = models.UserModel{
			NickName:   qqInfo.NickName,
			UserName:   qqInfo.NickName,
			Password:   hashPwd,
			Avatar:     qqInfo.Avatar,
			Addr:       addr,
			Token:      openID,
			IP:         ip,
			Role:       ctype.PermissionUser,
			SignStatus: ctype.SignQQ,
		}
		err = global.DB.Create(&user).Error
		if err != nil {
			global.Log.Error(err)
			res.FailWithMessage("注册失败", c)
			return
		}

	}

	//登录操作
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		NickName: user.NickName,
		Role:     int(user.Role),
		UserID:   user.ID,
		Avatar:   user.Avatar,
	})
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("token生成失败", c)
		return
	}
	err = global.DB.Create(&models.LoginDataModel{
		UserID:    user.ID,
		IP:        ip,
		NickName:  user.NickName,
		Token:     token,
		Device:    "",
		Addr:      addr,
		LoginType: ctype.SignQQ,
	}).Error

	if err != nil {
		global.Log.Errorf("创建QQ登录记录失败 err:%+v\n", err.Error())
	}
	res.OkWithData(token, c)
}
