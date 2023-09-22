package menu_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"github.com/gin-gonic/gin"
)

type MenuDetailRequest struct {
	Path string `json:"path" form:"path"`
}

// MenuDetailByPathView 菜单详情,根据路径查
// @Tags 菜单管理
// @Summary 菜单详情,根据路径查
// @Description 菜单详情,根据路径查
// @Param data query MenuDetailRequest  true  "路径参数"
// @Produce json
// @Success 200 {object} res.Response{data=MenuResponse}
// @Router /api/menus/detail [get]
func (MenuApi) MenuDetailByPathView(c *gin.Context) {
	var cr MenuDetailRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	var menuModel models.MenuModel
	err = global.DB.Take(&menuModel, "path = ?", cr.Path).Error
	if err != nil {
		res.FailWithMessage("菜单不存在", c)
		return
	}
	//查连接表
	var menuBanners []models.MenuBannerModel
	global.DB.Preload("BannerModel").Order("sort desc").Find(&menuBanners, "menu_id = ?", menuModel.ID)
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
