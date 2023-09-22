package images_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"fmt"
	"github.com/gin-gonic/gin"
)

// ImageRemoveView 删除图片
//
//	@Tags			图片管理
//	@Summary		删除图片
//	@Description	删除图片
//	@Param			data	body	models.RemoveRequest	true	"表示多个参数"
//	@Produce		json
//	@Success		200	{object}	res.Response{}
//	@Router			/api/images [delete]
func (ImagesApi) ImageRemoveView(c *gin.Context) {
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	var imageList []models.BannerModel //banner代表的是菜单表
	count := global.DB.Find(&imageList, cr.IDList).RowsAffected
	if count == 0 {
		res.FailWithMessage("文件不存在", c)
		return
	}
	global.DB.Create(&imageList)
	res.OkWithMessage(fmt.Sprintf("共删除 %d 张图片", count), c)
}
