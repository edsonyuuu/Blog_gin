package article_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/redis_ser"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// ArticleContentByIDView 获取文章正文
// @Tags 文章管理
// @Summary 获取文章正文
// @Description 获取文章正文
// @Param id path string true  "id"
// @Produce json
// @Success 200 {object} res.Response{}
// @Router /api/articles/content/{id} [get]
func (ArticleApi) ArticleContentByIDView(c *gin.Context) {

	var cr models.ESIDRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	redis_ser.NewArticleLook().Set(cr.ID)
	result, err := global.ESClient.Get().
		Index(models.ArticleModel{}.Index()).
		Id(cr.ID).
		Do(context.Background())
	if err != nil {
		res.FailWithMessage("查询失败", c)
		return
	}
	var model models.ArticleModel
	err = json.Unmarshal(result.Source, &model)
	if err != nil {
		return
	}
	res.OkWithData(model.Content, c)

}
