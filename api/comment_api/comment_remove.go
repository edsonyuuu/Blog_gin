package comment_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/redis_ser"
	"Blog_gin/utils"
	"Blog_gin/utils/jwts"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommentRemoveView 删除评论
// @Tags 评论管理
// @Summary 删除评论
// @Description 删除评论
// @param token header string true "token"
// @param id path int true "id"
// @Produce json
// @Success 200 {object} res.Response{}
// @Router /api/comments/{id} [delete]
func (CommentApi) CommentRemoveView(c *gin.Context) {

	var cr CommentIDRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	var commentModel models.CommentModel
	err = global.DB.Take(&commentModel, cr.ID).Error
	if err != nil {
		res.FailWithMessage("评论不存在", c)
		return
	}

	if !(commentModel.UserID == claims.UserID || claims.Role == 1) {
		res.FailWithMessage("权限错误,不可删除", c)
		return
	}
	//统计评论下的子评论数
	subCommentList := FindSubCommentCount(commentModel)
	count := len(subCommentList) + 1
	redis_ser.NewCommentCount().SetCount(commentModel.ArticleID, -count)

	//判断是否是子评论
	if commentModel.ParentCommentID != nil {
		//找到父评论，减掉对应的评论数
		global.DB.Model(&models.CommentModel{}).Where("id = ?", *commentModel.ParentCommentID).
			Update("comment_count", gorm.Expr("comment_count - ?", count))

	}

	//删除子评论以及当前评论
	var deleteCommentIDList []uint
	for _, model := range subCommentList {
		deleteCommentIDList = append(deleteCommentIDList, model.ID)
	}

	//反转，一个一个删除
	utils.Reverse(deleteCommentIDList)
	deleteCommentIDList = append(deleteCommentIDList, commentModel.ID)
	for _, id := range deleteCommentIDList {
		global.DB.Model(models.CommentModel{}).Delete("id = ?", id)
	}

	res.OkWithMessage(fmt.Sprintf("共删除 %d 条评论", len(deleteCommentIDList)), c)
	return

}
