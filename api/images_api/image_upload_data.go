package images_api

import (
	"Blog_gin/global"
	"Blog_gin/res"
	"Blog_gin/service/image_ser"
	"Blog_gin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
	"strings"
)

// ImageUploadDataView 上传多个图片，返回图片的url
//
//	@Tags			图片管理
//	@Summary		上传单个图片，返回图片的url
//	@Description	上传单个图片，返回图片的url
//	@Accept			multipart/form-data
//	@Param			image	formData	file	true	"文件上传"
//	@Produce		json
//	@Success		200	{object}	res.Response{}
//	@Router			/api/image [post]
func (ImagesApi) ImageUploadDataView(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		res.FailWithMessage("参数校验失败", c)
		return
	}
	fileName := file.Filename
	basePath := global.Config.Upload.Path
	filePath := path.Join(basePath, fileName)
	//

	//文件白名单判断
	nameList := strings.Split(fileName, ".")             //将该部分转换为小写字母，并将结果存储在 suffix 变量中
	suffix := strings.ToLower(nameList[len(nameList)-1]) //将文件后缀名转换为小写形式，以便后续处理
	if !utils.InList(suffix, image_ser.WhiteImageList) {
		res.FailWithMessage("非法文件", c)
		return
	}

	//判断文件大小
	size := float64(file.Size) / float64(1024*1024)
	if size >= float64(global.Config.Upload.Size) {
		msg := fmt.Sprintf("图片大小超过设定大小，当前大小为:%.2fMB,设定大小为:%dMB", size, global.Config.Upload.Size)
		res.FailWithMessage(msg, c)
		return
	}

	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}
	res.OkWithData("/"+filePath, c)
	return
}
