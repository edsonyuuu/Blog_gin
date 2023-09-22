package api

import (
	"Blog_gin/api/advert_api"
	"Blog_gin/api/images_api"
	"Blog_gin/api/menu_api"
	"Blog_gin/api/message_api"
	"Blog_gin/api/settings_api"
	"Blog_gin/api/tag_api"
	"Blog_gin/api/user_api"
)

type ApiGroup struct {
	SettingsApi settings_api.SettingsApi
	ImagesApi   images_api.ImagesApi
	UserApi     user_api.UserApi
	TagApi      tag_api.TagApi
	AdvertApi   advert_api.AdvertApi
	MenuApi     menu_api.MenuApi
	MessageApi  message_api.MessageApi
}

var ApiGroupApp = new(ApiGroup)
