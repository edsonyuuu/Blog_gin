package user_ser

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/models/ctype"
	"Blog_gin/utils"
	"Blog_gin/utils/pwd"
	"errors"
)

func (UserService) CreateUser(userName, nickName, password string, role ctype.Role, email string, ip string) error {
	//判断用户名是否存在
	var userModel models.UserModel
	err := global.DB.Take(&userModel, "user_name=?", userName).Error
	if err == nil {
		//即存在
		return errors.New("用户名已存在")
	}
	//对密码进行哈希操作
	hashPwd := pwd.HashPwd(password)

	//头像
	//
	avatar := "/uploads/avatar/default.jpg"
	addr := utils.GetAddr(ip)
	//入库
	err = global.DB.Create(&models.UserModel{
		NickName:   nickName,
		UserName:   userName,
		Password:   hashPwd,
		Avatar:     avatar,
		Email:      email,
		Addr:       addr,
		IP:         ip,
		Role:       role,
		SignStatus: ctype.SignEmail,
	}).Error

	if err != nil {
		global.Log.Error(err)
		return err
	}
	return nil
}
