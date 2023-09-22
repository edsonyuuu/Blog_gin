package menu_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"github.com/gin-gonic/gin"
)

type MenuNameResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Path  string `json:"path"`
}

// MenuNameList 是获取菜单名称列表的处理函数。
// @Summary 获取菜单名称列表
// @Description 获取菜单名称列表
// @Tags 菜单管理
// @Produce json
// @Success 200 {object} res.Response{data=MenuNameResponse}
// @Router /api/menu_names [get]
func (MenuApi) MenuNameList(c *gin.Context) {
	var menuNameList []MenuNameResponse
	global.DB.Model(models.MenuModel{}).Select("id", "title", "path").Scan(&menuNameList)
	res.OkWithData(menuNameList, c)
}
