package article_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/common"
	"Blog_gin/utils/jwts"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

type CollResponse struct {
	models.ArticleModel
	CreatedAt string `json:"created_at"`
}

// ArticleCollectListView 用户收藏的文章列表
// @Tags 文章管理
// @Summary 用户收藏的文章列表
// @Description 用户收藏的文章列表
// @Param data query models.PageInfo true  "表示多个参数"
// @Param token header string true  "token"
// @Router /api/articles/collects [get]
// @Produce json
// @Success 200 {object} res.Response{data=res.ListResponse[CollResponse]}
func (ArticleApi) ArticleCollectListView(c *gin.Context) {

	var cr models.PageInfo
	_ = c.ShouldBindQuery(&cr)
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)
	var articleIDList []interface{}
	list, count, err := common.ComList(models.UserCollectModel{UserID: claims.UserID}, common.Option{
		PageInfo: cr,
	})

	var collMap = map[string]string{}
	for _, model := range list {
		articleIDList = append(articleIDList, model.ArticleID)
		collMap[model.ArticleID] = model.CreatedAt.Format("2006-01-02 15:04:05")
	}

	boolSearch := elastic.NewTermsQuery("_id", articleIDList...)

	var collList = make([]CollResponse, 0)

	//传id列表，查es
	result, err := global.ESClient.Search(models.ArticleModel{}.Index()).Query(boolSearch).Size(1000).Do(context.Background())
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(result.Hits.TotalHits.Value, articleIDList)

	for _, hit := range result.Hits.Hits {
		var article models.ArticleModel
		err = json.Unmarshal(hit.Source, &article)
		if err != nil {
			global.Log.Error(err)
			continue
		}

		article.ID = hit.Id
		article.Content = ""
		collList = append(collList, CollResponse{
			ArticleModel: article,
			CreatedAt:    collMap[hit.Id],
		})
	}
	res.OkWithList(collList, count, c)
}
