package images_api

import (
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/common"
	"github.com/gin-gonic/gin"
)

// ImageListView 图片列表
//
//	@Tags			图片管理
//	@Summary		图片列表
//	@Description	图片列表
//	@Param			data	query	models.PageInfo	false	"查询参数"
//	@Produce		json
//
// Success 200 {object} res.Response{data=res.ListResponse[models.ImageModel]}
//
//	@Router			/api/images [get]
func (ImagesApi) ImageListView(c *gin.Context) {
	var cr models.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}
	list, count, err := common.ComList(models.BannerModel{}, common.Option{
		PageInfo: cr,
		Debug:    true,
	})

	res.OkWithList(list, count, c)
	return

}
