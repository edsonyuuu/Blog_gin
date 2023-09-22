package user_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/models/ctype"
	"Blog_gin/plugins/log_stash"
	"Blog_gin/res"
	"Blog_gin/utils"
	"Blog_gin/utils/jwts"
	"Blog_gin/utils/pwd"
	"fmt"
	"github.com/gin-gonic/gin"
)

type EmailLoginRequest struct {
	UserName string `json:"user_name" binding:"required" msg:"请输入用户名"`
	Password string `json:"password" binding:"required" msg:"请输入密码"`
}

// EmailLoginView 邮箱登录，返回token，用户信息需要从token中解码
//
//	@Tags			用户管理
//	@Summary		邮箱登录
//	@Description	邮箱登录，返回token，用户信息需要从token中解码
//	@Param			data	body	EmailLoginRequest	true	"查询参数"
//	@Produce		json
//	@Success		200	{object}	res.Response{}
//	@Router			/api/email_login [post]
func (UserApi) EmailLoginView(c *gin.Context) {
	var cr EmailLoginRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, &cr, c)
		return
	}

	log := log_stash.NewLogByGin(c)

	var user models.UserModel

	err = global.DB.Take(&user, "user_name = ? or email = ?", cr.UserName, cr.UserName).Error
	if err != nil {
		//没找到信息
		global.Log.Warn("用户名不存在")
		log.Warn(fmt.Sprintf("%s 用户不存在", cr.UserName))
		res.FailWithMessage("用户名或者密码错误", c)
		return
	}

	//校验密码
	isCheck := pwd.CheckPwd(user.Password, cr.Password)
	if !isCheck {
		global.Log.Warn("用户名或密码错误")
		log.Warn(fmt.Sprintf("用户名密码错误 %s %s", cr.UserName, cr.Password))
		res.FailWithMessage("用户名或者密码错误", c)
		return
	}

	//登录成功,生成token
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		NickName: user.NickName,
		Role:     int(user.Role),
		UserID:   user.ID,
		Avatar:   user.Avatar,
	})
	if err != nil {
		global.Log.Error(err)
		log.Error(fmt.Sprintf("token生成失败 %s", err.Error()))
		res.FailWithMessage("token生成失败", c)
		return
	}
	ip, addr := utils.GetAddrByGin(c)
	log = log_stash.New(ip, token)
	log.Info("登录成功")

	err = global.DB.Create(&models.LoginDataModel{
		UserID:    user.ID,
		IP:        ip,
		NickName:  user.NickName,
		Token:     token,
		Device:    "",
		Addr:      addr,
		LoginType: ctype.SignEmail,
	}).Error

	if err != nil {
		global.Log.Error(err.Error())
	}

	res.OkWithData(token, c)
}
