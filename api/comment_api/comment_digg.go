package comment_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/redis_ser"
	"fmt"
	"github.com/gin-gonic/gin"
)

type CommentIDRequest struct {
	ID uint `json:"id" form:"id" uri:"id"`
}

// CommentDigg
// @Summary 点赞评论
// @Description 点赞评论
// @Tags 评论管理
// @Accept  json
// @Produce  json
// @Param comment body CommentIDRequest true "评论参数"
// @Success 200 {object} res.Response
// @Router /api/comments/{id} [get]
func (CommentApi) CommentDigg(c *gin.Context) {
	var cr CommentIDRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	var commentModel models.CommentModel
	err = global.DB.Take(&commentModel, cr.ID).Error
	if err != nil {
		res.FailWithMessage("评论不存在", c)
		return
	}

	redis_ser.NewCommentDigg().Set(fmt.Sprintf("%d", cr.ID))

	res.OkWithMessage("评论点赞成功", c)
	return
}
