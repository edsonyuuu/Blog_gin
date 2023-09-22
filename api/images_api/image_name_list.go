package images_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ImageResponse struct {
	ID   uint   `json:"id"`
	Path string `json:"path"` // 图片路径
	Name string `json:"name"` //图片名称
}

// ImageNameListView  图片名称列表
//
//	@Tags			图片管理
//	@Summary		图片名称列表
//	@Description	图片名称列表
//	@Produce		json
//	@Success		200	{object}	res.Response{data=[]ImageResponse}
//	@Router			/api/images_names [get]
func (ImagesApi) ImageNameListView(c *gin.Context) {
	var imageList []ImageResponse
	global.DB.Model(models.BannerModel{}).Select("id", "path", "name").Scan(&imageList)
	fmt.Println(imageList)
	res.OkWithData(imageList, c)
}
