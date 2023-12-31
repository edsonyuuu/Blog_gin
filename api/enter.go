package api

import (
	"Blog_gin/api/advert_api"
	"Blog_gin/api/article_api"
	"Blog_gin/api/chat_api"
	"Blog_gin/api/comment_api"
	"Blog_gin/api/data_api"
	"Blog_gin/api/digg_api"
	"Blog_gin/api/images_api"
	"Blog_gin/api/log_api"
	"Blog_gin/api/menu_api"
	"Blog_gin/api/message_api"
	"Blog_gin/api/news_api"
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
	ArticleApi  article_api.ArticleApi
	ChatApi     chat_api.ChatApi
	NewsApi     news_api.NewsApi
	CommentApi  comment_api.CommentApi
	DataApi     data_api.DataApi
	LogApi      log_api.LogApi
	DiggApi     digg_api.DiggApi
}

var ApiGroupApp = new(ApiGroup)
