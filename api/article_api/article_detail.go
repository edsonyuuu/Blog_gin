package article_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/es_ser"
	"Blog_gin/service/redis_ser"
	"Blog_gin/utils"
	"Blog_gin/utils/jwts"
	"github.com/gin-gonic/gin"
)

type ArticleDetailRequest struct {
	Title string `json:"title" form:"title"`
}

// ArticleDetailByTitleView 根据文章标题查看详情
// @Tags 文章管理
// @Summary 根据文章标题查看详情
// @Description 根据文章标题查看详情
// @Produce json
// @Param data query ArticleDetailRequest    false  "表示多个参数"
// @Success 200 {object} res.Response{}
// @Router /api/articles/detail [get]
func (ArticleApi) ArticleDetailByTitleView(c *gin.Context) {
	var cr ArticleRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	model, err := es_ser.CommDetailByKeyword(cr.Title)
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}

	res.OkWithData(model, c)

}

type ArticleDetailResponse struct {
	models.ArticleModel
	IsCollect bool `json:"is_collect"` //用户是否收藏
}

// ArticleDetailView 文章详情
// @Tags 文章管理
// @Summary 文章详情
// @Description 文章详情
// @Param id path string true  "id"
// @Produce json
// @Success 200 {object} res.Response{data=models.ArticleModel}
// @Router /api/articles/{id} [get]
func (ArticleApi) ArticleDetailView(c *gin.Context) {
	var cr models.ESIDRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}
	redis_ser.NewArticleLook().Set(cr.ID)

	model, err := es_ser.CommDetail(cr.ID)
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}

	isCollect := IsUserArticleColl(c, model.ID)
	var articleDetail = ArticleDetailResponse{
		ArticleModel: model,
		IsCollect:    isCollect,
	}

	res.OkWithData(articleDetail, c)

}

func IsUserArticleColl(c *gin.Context, articleID string) (isCollect bool) {
	//判断用户是否正常登录
	token := c.GetHeader("token")
	if token == "" {
		return
	}
	claims, err := jwts.ParseToken(token)
	if err != nil {
		return
	}
	//判断是否在redis中
	keys := global.Redis.Keys("logout_*").Val()
	if utils.InList("logout_"+token, keys) {
		return
	}
	var count int64
	global.DB.Model(models.UserCollectModel{}).Where("user_id = ? and article_id = ?", claims.UserID, articleID).Count(&count)
	if count == 0 {
		return
	}
	return true
}
