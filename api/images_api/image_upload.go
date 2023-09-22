package images_api

import (
	"Blog_gin/global"
	"Blog_gin/res"
	"Blog_gin/service"
	"Blog_gin/service/image_ser"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/fs"
	"os"
)

// ImageUploadView 处理上传多个图片,返回图片的url
//
//	@Tags			图片管理
//	@Summary		上传多个图片，返回图片的url
//	@Description	上传多个图片，返回图片的url
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			images	formData	file	true	"图片文件上传"
//	@Success		200		{object}	res.Response{}
//	@Router			/api/images [post]
func (ImagesApi) ImageUploadView(c *gin.Context) {
	//上传多个图片

	form, err := c.MultipartForm()
	if err != nil {
		fmt.Printf("接收参数失败 err:%+v\n", err.Error())
		res.FailWithMessage(err.Error(), c)
		return
	}

	fileList, ok := form.File["images"]

	if !ok {
		fmt.Println("文件不存在")
		res.FailWithMessage("不存在的文件", c)
		return
	}

	//判断路径是否存在
	basePath := global.Config.Upload.Path
	_, err = os.ReadDir(basePath)
	if err != nil {
		//递归创建
		//目录路径和权限模式,该函数的作用是在文件系统中创建一个或多个嵌套的目录
		err = os.MkdirAll(basePath, fs.ModePerm)
		if err != nil {
			global.Log.Error(err)
			fmt.Printf("递归步骤出现问题 err:%+v\n", err.Error())
		}
	}

	//不存在的情况
	var resList []image_ser.FileUploadResponse

	for _, file := range fileList {

		//上传文件
		//调用service.ServiceApp.ImageService.ImageUploadService(file)方法进行文件上传操作，
		//并将返回的结果保存在serviceRes变量中。
		serviceRes := service.ServiceApp.ImageService.ImageUploadService(file)
		if !serviceRes.IsSuccess {
			resList = append(resList, serviceRes)
			continue
		}

		//成功
		if !global.Config.QiNiu.Enable {
			//本地保存
			err = c.SaveUploadedFile(file, serviceRes.FileName)
			if err != nil {
				fmt.Printf("本地保存失败 err:%+v\n", err.Error())
				global.Log.Error(err)
				serviceRes.Msg = err.Error()
				serviceRes.IsSuccess = false
				resList = append(resList, serviceRes)
				continue
			}
		}
		resList = append(resList, serviceRes)
	}
	res.OkWithData(resList, c)
}
