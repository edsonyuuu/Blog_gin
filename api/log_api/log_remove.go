package log_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/plugins/log_stash"
	"Blog_gin/res"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (LogApi) LogRemoveListView(c *gin.Context) {
	var cr models.RemoveRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	var list []log_stash.LogStashModel
	count := global.DB.Find(&list, cr.IDList).RowsAffected
	if count == 0 {
		res.FailWithMessage("日志不存在", c)
		return
	}
	global.DB.Delete(&list)
	res.OkWithMessage(fmt.Sprintf("共删除 %d 个日志", count), c)

}
