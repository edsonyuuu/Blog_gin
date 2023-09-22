package advert_api

import (
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/common"
	"github.com/gin-gonic/gin"
	"strings"
)

// AdvertListView 广告列表
// @Tags 广告管理
// @Summary 广告列表
// @Description 广告列表
// @Param data query models.PageInfo    false  "查询参数"
// @Produce json
// @Success 200 {object} res.Response{data=res.ListResponse[models.AdvertModel]}
// @Router /api/adverts [get]
func (AdvertApi) AdvertListView(c *gin.Context) {
	var cr models.PageInfo
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	// 判断 Referer 是否包含admin，如果是，就全部返回，不是，就返回is_show=true
	referer := c.GetHeader("Gvb_referer")
	isShow := true
	if strings.Contains(referer, "admin") {
		//管理员进入
		isShow = false
	}
	list, count, _ := common.ComList(models.AdvertModel{IsShow: isShow}, common.Option{
		PageInfo: cr,
		Debug:    true,
	})

	res.OkWithList(list, count, c)

}
