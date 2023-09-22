package redis_ser

import (
	"Blog_gin/global"
	"Blog_gin/utils"
	"time"
)

const prefix = "logout_"

// Logout  针对注销的操作
func Logout(token string, diff time.Duration) error {
	err := global.Redis.Set(prefix+token, "", diff).Err()
	global.Log.Errorf("service_usr 注销登录错误 err:%+v\n", err.Error())
	return err
}

func CheckLogout(token string) bool {
	keys := global.Redis.Keys(prefix + "*").Val()
	if utils.InList(prefix+token, keys) {
		return true
	}
	return false
}
