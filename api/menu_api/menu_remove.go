package menu_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MenuRemoveView
// @Summary 删除菜单
// @Description 根据请求的菜单ID列表删除菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param request body models.RemoveRequest true "菜单ID列表"
// @Success 200 {object} res.Response{}
// @Router /api/menus [delete]
func (MenuApi) MenuRemoveView(c *gin.Context) {
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	var menuList []models.MenuModel
	count := global.DB.Find(&menuList, cr.IDList).RowsAffected
	if count == 0 {
		res.FailWithMessage("菜单不存在", c)
		return
	}

	//事务
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&menuList).Association("Banners").Clear()
		if err != nil {
			global.Log.Error(err)
			return err
		}
		err = tx.Delete(&menuList).Error
		if err != nil {
			global.Log.Error(err)
			return err
		}
		return nil
	})

	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("删除菜单失败", c)
		return
	}

	res.OkWithMessage(fmt.Sprintf("共删除 %d 个菜单", count), c)
}
