package tag_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
)

// TagUpdateView 更新标签
// @Tags 标签管理
// @Summary 更新标签
// @Description 更新标签
// @Param data body TagRequest    true  "标签的一些参数"
// @Produce json
// @Success 200 {object} res.Response{data=string}
// @Router /api/tags/:id [put]
func (TagApi) TagUpdateView(c *gin.Context) {
	id := c.Param("id")

	var cr TagRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(err, &cr, c)
		return
	}

	var tag models.TagModel
	err := global.DB.Take(&tag, id).Error
	if err != nil {
		res.FailWithMessage("标签不存在", c)
		return
	}

	//结构体转map的第三方包
	maps := structs.Map(&cr)
	err = global.DB.Model(&tag).Updates(maps).Error
	if err != nil {
		global.Log.Errorf("Tag Update Error :%+v\n", err.Error())
		res.FailWithMessage("标签修改失败", c)
		return
	}
	res.OkWithMessage("标签修改成功", c)
}
