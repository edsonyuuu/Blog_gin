package user_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/plugins/email"
	"Blog_gin/res"
	"Blog_gin/utils/jwts"
	"Blog_gin/utils/pwd"
	"Blog_gin/utils/random"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type BindEmailRequest struct {
	Email    string  `json:"email" binding:"required,email" msg:"邮箱非法"`
	Code     *string `json:"code"`
	Password string  `json:"password"`
}

// UserBindEmailView 用户绑定邮箱
// @Tags 用户管理
// @Summary 用户绑定邮箱
// @Description 用户绑定邮箱
// @Param data body BindEmailRequest true  "查询参数"
// @Param token header string true "token"
// @Produce json
// @Success 200 {object} res.Response{}
// @Router /api/user_bind_email [post]
func (UserApi) UserBindEmailView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)
	//用户绑定邮箱，第一次输入邮箱
	//后台给这个邮箱发送验证码
	var cr BindEmailRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, &cr, c)
		return
	}
	session := sessions.Default(c)
	if cr.Code == nil {
		//第一次后天发送验证码
		//生成四位验证码，将生成的验证码存入session
		code := random.Code(4)
		//写入session
		session.Set("valid_code", code)
		err = session.Save()
		if err != nil {
			global.Log.Error(err)
			res.FailWithMessage("session错误", c)
			return
		}
		err = email.NewCode().Send(cr.Email, "你的验证码是 "+code)
		if err != nil {
			global.Log.Error(err)
		}
		res.OkWithMessage("验证码已发送，请查收", c)
		return
	}
	//第二次,用户输入邮箱，验证码，密码
	code := session.Get("valid_code")
	//校验验证码
	if code != *cr.Code {
		res.FailWithMessage("验证码错误", c)
		return
	}

	//修改用户的邮箱
	var user models.UserModel
	err = global.DB.Take(&user, claims.UserID).Error
	if err != nil {
		res.FailWithMessage("用户不存在", c)
		return
	}
	hashPwd := pwd.HashPwd(cr.Password)
	//第一次邮箱和第二次的邮箱也要进行校验

	err = global.DB.Model(&user).Updates(map[string]any{
		"email":    cr.Email,
		"password": hashPwd,
	}).Error

	if err != nil {
		global.Log.Errorf("一致性校验出错 err:%+v\n", err.Error())
		res.FailWithMessage("绑定邮箱失败", c)
		return
	}

	//完成绑定
	res.OkWithMessage("邮箱绑定成功", c)
}
