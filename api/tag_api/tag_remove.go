package tag_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"fmt"
	"github.com/gin-gonic/gin"
)

// TagRemoveView 批量删除标签
// @Tags 标签管理
// @Summary 批量删除标签
// @Description 批量删除标签
// @Param data body models.RemoveRequest    true  "标签id列表"
// @Produce json
// @Success 200 {object} res.Response{data=string}
// @Router /api/tags [delete]
func (TagApi) TagRemoveView(c *gin.Context) {
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	var tagList []models.TagModel
	count := global.DB.Find(&tagList, cr.IDList).RowsAffected
	if count == 0 {
		res.FailWithMessage("标签不存在", c)
		return
	}

	//标签下文章如何处理？先留下这个问题
	global.DB.Delete(&tagList)
	res.OkWithMessage(fmt.Sprintf("共删除 %d 个标签", count), c)

}
