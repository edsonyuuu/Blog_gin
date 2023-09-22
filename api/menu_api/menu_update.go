package menu_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
)

// MenuUpdateView 函数用于更新菜单视图
// @Summary 更新菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Param request body MenuRequest true "请求参数"
// @Success 200 {object} res.Response{}
// @Router /api/menus/{id} [put]
func (MenuApi) MenuUpdateView(c *gin.Context) {
	var cr MenuRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, &cr, c)
		return
	}

	id := c.Param("id")

	//把之前的banner清空
	var menuModel models.MenuModel
	err = global.DB.Take(&menuModel, id).Error
	if err != nil {
		res.FailWithMessage("菜单不存在", c)
		return
	}

	global.DB.Model(&menuModel).Association("Banners").Clear()
	//如果选择了banner，进行添加
	if len(cr.ImageSortList) > 0 {
		//操作第三张表
		var bannerList []models.MenuBannerModel
		for _, sort := range cr.ImageSortList {
			bannerList = append(bannerList, models.MenuBannerModel{
				MenuID:   menuModel.ID,
				BannerID: menuModel.ID,
				Sort:     sort.Sort,
			})
		}
		err = global.DB.Create(&bannerList).Error
		if err != nil {
			global.Log.Error(err)
			res.FailWithMessage("创建菜单图片失败", c)
			return
		}
	}

	//普通更新
	maps := structs.Map(&cr)
	err = global.DB.Model(&menuModel).Updates(maps).Error
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("修改菜单失败", c)
		return
	}
	res.OkWithMessage("修改菜单成功", c)

}
