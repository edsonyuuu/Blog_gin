package menu_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"github.com/gin-gonic/gin"
)

// MenuDetailView
// @Summary 获取菜单详情
// @Description 根据菜单ID获取菜单详情和相关横幅信息
// @Tags 菜单管理
// @Produce json
// @Param id path int true "菜单ID"
// @Success 200 {object} res.Response{data=MenuResponse}
// @Router /menu/{id} [get]
func (MenuApi) MenuDetailView(c *gin.Context) {
	id := c.Param("id")
	var menuModel models.MenuModel
	err := global.DB.Take(&menuModel, id).Error
	if err != nil {
		res.FailWithMessage("菜单不存在", c)
		return
	}

	//查连接表
	var menuBanners []models.MenuBannerModel
	global.DB.Preload("BannerModel").Order("sort desc").Find(&menuBanners, "menu_id = ?", id)
	var banners = make([]Banner, 0)
	for _, banner := range menuBanners {
		if menuModel.ID != banner.MenuID {
			continue
		}
		banners = append(banners, Banner{
			ID:   banner.BannerID,
			Path: banner.BannerModel.Path,
		})
	}
	menuResponse := MenuResponse{
		MenuModel: menuModel,
		Banners:   banners,
	}
	res.OkWithData(menuResponse, c)

}
