package flag

import (
	"Blog_gin/global"
	models2 "Blog_gin/models"
	"Blog_gin/plugins/log_stash"
)

func Makemigrations() {
	var err error
	//SetupJoinTable,设置表之间的多对多关系
	//global.DB.SetupJoinTable(&models2.UserModel{}, "CollectsModels", &models2.UserCollectModel{})
	global.DB.SetupJoinTable(&models2.MenuModel{}, "Banners", &models2.MenuBannerModel{})
	//生成表结构
	err = global.DB.Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(
			&models2.BannerModel{},
			&models2.TagModel{},
			&models2.MessageModel{},
			&models2.AdvertModel{},
			&models2.UserModel{},
			&models2.CommentModel{},
			&models2.ArticleModel{},
			&models2.UserCollectModel{},
			&models2.MenuModel{},
			&models2.MenuBannerModel{},
			&models2.FadeBackModel{},
			&models2.LoginDataModel{},
			&models2.ChatModel{},
			&log_stash.LogStashModel{},
		)
	if err != nil {
		global.Log.Error("[ error ] 生成数据库表结构失败")
		return
	}
	global.Log.Info("[ success ] 生成数据库表结构成功！")
}
