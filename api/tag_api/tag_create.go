package tag_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"github.com/gin-gonic/gin"
)

type TagRequest struct {
	Title string `json:"title" binding:"required" msg:"请输入标题" structs:"title"` //标题
}

// TagCreateView 添加标签
// @Tags 标签管理
// @Summary 创建标签
// @Description 创建标签
// @Param data body TagRequest    true  "表示多个参数"
// @Produce json
// @Success 200 {object} res.Response{}
// @Router /api/tags [post]
func (TagApi) TagCreateView(c *gin.Context) {
	var cr TagRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(err, &cr, c)
		return
	}
	//
	var tag models.TagModel
	err := global.DB.Take(&tag, "title = ?", cr.Title).Error
	if err == nil {
		res.FailWithMessage("该标签已存在", c)
		return
	}

	err = global.DB.Create(&models.TagModel{
		Title: cr.Title,
	}).Error
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("添加标签失败", c)
		return
	}

	res.OkWithMessage("添加标签成功", c)

}
