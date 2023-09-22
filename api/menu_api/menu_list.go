package menu_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"github.com/gin-gonic/gin"
)

type Banner struct {
	ID   uint   `json:"id"`
	Path string `json:"path"`
}

type MenuResponse struct {
	models.MenuModel
	Banners []Banner `json:"banners"`
}

// MenuListView
// @Summary 菜单列表
// @Description 获取菜单列表及对应的横幅信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} res.Response{}
// @Router /api/menus [get]
func (MenuApi) MenuListView(c *gin.Context) {
	var menuList []models.MenuModel
	var menuIDlist []uint

	global.DB.Order("sort desc").Find(&menuList).Select("id").Scan(&menuIDlist)
	//查连接表
	var menuBanners []models.MenuBannerModel
	global.DB.Preload("BannerModel").Order("sort desc").Find(&menuBanners, "menu_id in ?", menuIDlist)
	var menus = make([]MenuResponse, 0)

	for _, model := range menuList {
		//model为一个菜单
		var banners = make([]Banner, 0)
		for _, banner := range menuBanners {
			if model.ID != banner.MenuID {
				continue
			}
			banners = append(banners, Banner{
				ID:   banner.BannerID,
				Path: banner.BannerModel.Path,
			})
		}
		menus = append(menus, MenuResponse{
			MenuModel: model,
			Banners:   banners,
		})
	}

	res.OkWithList(menus, int64(len(menus)), c)
}
