package digg_api

import (
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/redis_ser"
	"github.com/gin-gonic/gin"
)

func (DiggApi) DiggArticleView(c *gin.Context) {
	var cr models.ESIDRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}
	// 对长度校验
	// 查es
	redis_ser.NewDigg().Set(cr.ID)
	res.OkWithMessage("文章点赞成功", c)
}
