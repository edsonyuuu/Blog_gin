package service

import (
	"Blog_gin/service/image_ser"
	"Blog_gin/service/user_ser"
)

type ServiceGroup struct {
	ImageService image_ser.ImageService
	UserService  user_ser.UserService
}

var ServiceApp = new(ServiceGroup)
