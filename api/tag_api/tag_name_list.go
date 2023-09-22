package tag_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

type TagResponse struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// TagNameListView 标签名称列表
// @Tags 标签管理
// @Summary 标签名称列表
// @Description 标签名称列表
// @Produce json
// @Success 200 {object} res.Response{data=[]TagResponse}
// @Router /api/tag_names [get]
func (TagApi) TagNameListView(c *gin.Context) {
	type T struct {
		DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
		SumOtherDocCount        int `json:"sum_other_doc_count"`
		Buckets                 []struct {
			Key      string `json:"key"`
			DocCount int    `json:"doc_count"`
		} `json:"buckets"`
	}

	query := elastic.NewBoolQuery()
	agg := elastic.NewTermsAggregation().Field("tags")
	result, err := global.ESClient.Search(models.ArticleModel{}.Index()).Query(query).Aggregation("tags", agg).Size(0).Do(context.Background())
	if err != nil {
		global.Log.Errorf("tags search err:%+v\n", err.Error())
		return //这里感觉需要添加
	}

	byteData := result.Aggregations["tags"]
	var tagType T
	json.Unmarshal(byteData, &tagType)

	var tagList = make([]TagResponse, 0)
	for _, bucket := range tagType.Buckets {
		tagList = append(tagList, TagResponse{
			Label: bucket.Key,
			Value: bucket.Key,
		})
	}
	res.OkWithData(tagList, c)
}
