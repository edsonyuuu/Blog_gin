package article_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

// FullTextSearchView 全文搜索
// @Summary 全文搜索
// @Description 根据关键词进行全文搜索，返回匹配结果和高亮内容
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param request body models.PageInfo true "请求参数"
// @Success 200 {object} res.Response
// @Router /search [get]
func (ArticleApi) FullTextSearchView(c *gin.Context) {
	var cr models.PageInfo
	_ = c.ShouldBindQuery(&cr)
	//这里创建了一个布尔查询boolQuery，用于构建搜索查询的布尔逻辑。
	boolQuery := elastic.NewBoolQuery()

	//如果查询参数cr.Key不为空，则创建一个多字段匹配查询MultiMatchQuery，并将它添加到boolQuery中。这个查询会在"title"和"body"字段上进行全文匹配。
	if cr.Key != "" {
		boolQuery.Must(elastic.NewMultiMatchQuery(cr.Key, "title", "body"))
	}

	//这里使用全局的Elasticsearch客户端global.ESClient进行搜索操作。使用models.FullTextModel{}.Index()指定了要搜索的索引。
	//使用之前构建的boolQuery作为查询条件，并设置了结果的大小为100条。同时也进行了高亮设置，将"body"字段进行高亮显示
	result, err := global.ESClient.Search(models.FullTextModel{}.Index()).
		Query(boolQuery).Highlight(elastic.NewHighlight().Field("body")).Size(100).Do(context.Background())
	if err != nil {
		return
	}
	//result.Hits.TotalHits.Value表示搜索结果的总命中数。
	count := result.Hits.TotalHits.Value

	//这里使用循环遍历搜索结果的每个命中项。
	//通过json.Unmarshal将命中项的原始数据解析为models.FullTextModel类型的结构体。
	//如果命中项有高亮字段"body"，则将高亮内容赋值给model.Body。最后将model添加到fullTextList中。
	fullTextList := make([]models.FullTextModel, 0)
	for _, hit := range result.Hits.Hits {
		var model models.FullTextModel
		_ = json.Unmarshal(hit.Source, &model)

		body, ok := hit.Highlight["body"]
		if ok {
			model.Body = body[0]
		}

		fullTextList = append(fullTextList, model)

	}
	res.OkWithList(fullTextList, count, c)
}
