package article_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type CategoryResponse struct {
	Label string `json:"label"` // 标签
	Value string `json:"value"`
}

// ArticleCategoryListView 文章分类列表
// @Tags 文章管理
// @Summary 文章分类列表
// @Description 文章分类列表
// @Produce json
// @Success 200 {object} res.Response{data=[]CategoryResponse}
// @Router /api/categorys [get]
func (ArticleApi) ArticleCategoryListView(c *gin.Context) {

	type T struct {
		DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
		SumOtherDocCount        int `json:"sum_other_doc_count"`
		Buckets                 []struct {
			Key      string `json:"key"`
			DocCount int    `json:"doc_count"`
		} `json:"buckets"`
	}
	//DocCountErrorUpperBound（整数类型）：表示文档计数错误的上限。
	//SumOtherDocCount（整数类型）：表示其他文档计数的总和。
	//Buckets（结构体切片）：表示一个名为"Buckets"的结构体切片，其中每个结构体具有两个字段：
	//Key（字符串类型）：表示键值。
	//DocCount（整数类型）：表示文档计数。
	//
	//创建了一个新的聚合（aggregation）对象，具体类型为Terms Aggregation（术语聚合）。该聚合使用category字段进行分组。
	agg := elastic.NewTermsAggregation().Field("category")
	//Search(models.ArticleModel{}.Index())：创建一个搜索请求，指定了要搜索的索引，这里使用了models.ArticleModel{}.Index()来获取索引名称。
	//Query(elastic.NewBoolQuery())：设置搜索请求的查询条件，这里使用了空的布尔查询，表示匹配所有文档。
	//Aggregation("category", agg)：将之前创建的聚合对象agg添加到搜索请求中，以对搜索结果进行聚合操作。"category"是聚合操作的名称，可以用于后续获取聚合结果。
	//Size(0)：设置搜索请求的大小为0，表示只返回聚合结果，不返回匹配的文档。
	//Do(context.Background())：执行搜索请求，并返回搜索结果。context.Background()用于提供上下文。
	result, err := global.ESClient.Search(models.ArticleModel{}.Index()).Query(elastic.NewBoolQuery()).Aggregation("category", agg).Size(0).Do(context.Background())
	if err != nil {
		logrus.Error(err)
		return
	}

	byteData := result.Aggregations["category"]
	var categoryType T
	//反序列化
	_ = json.Unmarshal(byteData, &categoryType)
	var categoryList = make([]CategoryResponse, 0)
	for _, bucket := range categoryType.Buckets {
		categoryList = append(categoryList, CategoryResponse{
			Label: bucket.Key,
			Value: bucket.Key,
		})
	}
	res.OkWithData(categoryList, c)

}
