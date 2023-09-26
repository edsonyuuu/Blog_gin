package comment_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/es_ser"
	"Blog_gin/service/redis_ser"
	"Blog_gin/utils/jwts"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentRequest struct {
	ArticleID       string `json:"article_id" binding:"required" msg:"请选择文章"`
	Content         string `json:"content" binding:"required" msg:"请输入评论内容"`
	ParentCommentID *uint  `json:"parent_comment_id"` //父评论id
}

// CommentCreateView
// @Summary 创建评论
// @Description 创建文章评论
// @Tags 评论管理
// @Accept  json
// @Produce  json
// @Param comment body CommentRequest true "评论参数"
// @Param token header string true "token"
// @Success 200 {object} res.Response
// @Router /api/comments [post]
func (CommentApi) CommentCreateView(c *gin.Context) {
	var cr CommentRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	//判断文章是否存在
	_, err = es_ser.CommDetail(cr.ArticleID)
	if err != nil {
		res.FailWithMessage("文章不存在", c)
		return
	}

	//判断是否为子评论
	if cr.ParentCommentID != nil {
		//给父评论数＋1
		//父评论id
		var parentComment models.CommentModel
		//寻找父评论
		err = global.DB.Take(&parentComment, cr.ParentCommentID).Error
		if err != nil {
			res.FailWithMessage("父评论不存在", c)
			return
		}
		if parentComment.ArticleID != cr.ArticleID {
			res.FailWithMessage("评论文章不一致", c)
			return
		}
		//给父评论数+1
		global.DB.Model(&parentComment).Update("comment_count", gorm.Expr("comment_count + 1"))
	}
	global.DB.Create(&models.CommentModel{
		ParentCommentID: cr.ParentCommentID,
		Content:         cr.Content,
		ArticleID:       cr.ArticleID,
		UserID:          claims.UserID,
	})

	//给文章评论数+1
	redis_ser.NewCommentCount().Set(cr.ArticleID)
	res.OkWithMessage("文章评论成功", c)
	return
}
