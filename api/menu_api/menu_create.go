package menu_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/models/ctype"
	"Blog_gin/res"
	"github.com/gin-gonic/gin"
)

type ImageSort struct {
	ImageID uint `json:"image_id"`
	Sort    int  `json:"sort"`
}

type MenuRequest struct {
	Title         string      `json:"title" binding:"required" msg:"请完善菜单名称" structs:"title"`
	Path          string      `json:"path" binding:"required" msg:"请完善菜单路径" structs:"path"`
	Slogan        string      `json:"slogan" structs:"slogan"`
	Abstract      ctype.Array `json:"abstract" structs:"abstract"`
	AbstractTime  int         `json:"abstract_time" structs:"abstract_time"` // 切换的时间，单位秒
	BannerTime    int         `json:"banner_time" structs:"banner_time"`     // 切换的时间，单位秒
	Sort          int         `json:"sort" structs:"sort"`                   // 菜单的序号
	ImageSortList []ImageSort `json:"image_sort_list" structs:"-"`           // 具体图片的顺序
}

// MenuCreateView
// @Summary 创建菜单
// @Description 创建菜单并关联图片
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param data body MenuRequest true "菜单请求参数"
// @Success 200 {object} res.Response{}
// @Router /api/menus [post]
func (MenuApi) MenuCreateView(c *gin.Context) {
	var cr MenuRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, &cr, c)
		return
	}

	//判断
	var menuList models.MenuModel
	count := global.DB.Find(&menuList, "title = ? and path ?", cr.Title, cr.Path).RowsAffected
	if count > 0 {
		res.FailWithMessage("重复的菜单", c)
		return
	}

	//创建banner数据入库
	menuModel := models.MenuModel{
		Title:        cr.Title,
		Path:         cr.Path,
		Slogan:       cr.Slogan,
		Abstract:     cr.Abstract,
		AbstractTime: cr.AbstractTime,
		BannerTime:   cr.BannerTime,
		Sort:         cr.Sort,
	}
	err = global.DB.Create(&menuModel).Error
	if err != nil {
		global.Log.Errorf("创建banner失败 err:%+v\n", err.Error())
		return
	}
	if len(cr.ImageSortList) == 0 {
		res.OkWithMessage("菜单添加成功", c)
		return
	}

	var menuBannerList []models.MenuBannerModel

	for _, sort := range cr.ImageSortList {
		//解析判断是否真正存在这张图片
		menuBannerList = append(menuBannerList, models.MenuBannerModel{
			MenuID:   menuModel.ID,
			BannerID: sort.ImageID,
			Sort:     sort.Sort,
		})
	}
	//给第三张表入库
	err = global.DB.Create(&menuBannerList).Error
	if err != nil {
		global.Log.Errorf("bannerMenuModels 入库失败 err:%+v\n", err.Error())
		res.FailWithMessage("菜单图片关联失败", c)
		return
	}

	res.OkWithMessage("菜单添加成功", c)
}
