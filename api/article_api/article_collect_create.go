package article_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/es_ser"
	"Blog_gin/utils/jwts"
	"github.com/gin-gonic/gin"
)

// ArticleCollectCreateView 收藏文章
// @Summary 收藏文章
// @Description 收藏文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param token header string    false  "token"
// @Param	data	query models.ESIDRequest	false	"请求参数"
// @Success 200 {object} res.Response{}
// @Router /api/articles/collects [post]
func (ArticleApi) ArticleCollectCreateView(c *gin.Context) {
	var cr models.ESIDRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	model, err := es_ser.CommDetail(cr.ID)
	if err != nil {
		res.FailWithMessage("文章不存在", c)
		return
	}

	var coll models.UserCollectModel
	err = global.DB.Take(&coll, "user_id = ? and article_id = ?", claims.UserID, cr.ID).Error
	var num = -1
	if err != nil {
		global.DB.Create(&models.UserCollectModel{
			UserID:    claims.UserID,
			ArticleID: cr.ID,
		})
		//文章的收藏数+1
		num = 1
	}

	//取消收藏 -1
	global.DB.Delete(&coll)

	err = es_ser.ArticleUpdate(cr.ID, map[string]any{
		"collects_count": model.CollectsCount + num,
	})
	if num == 1 {
		res.OkWithMessage("收藏文章成功", c)
		return
	} else {
		res.OkWithMessage("取消收藏成功", c)
		return
	}
}
