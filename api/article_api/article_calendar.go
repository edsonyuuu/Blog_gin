package article_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"time"
)

type CalendarResponse struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type BucketsType struct {
	Buckets []struct {
		KeyAsString string `json:"key_as_string"`
		Key         int64  `json:"key"`
		DocCount    int    `json:"doc_count"`
	} `json:"buckets"`
}

var DateCount = map[string]int{}

// ArticleCalendarView
// @Summary 文章日历视图
// @Description 获取文章的日历视图，统计每天发布的文章数量
// @Tags 文章管理
// @Accept json
// @Produce json
// @Success 200 {object} res.Response{data=[]CalendarResponse}
// @Router /api/articles/calendar [get]
func (ArticleApi) ArticleCalendarView(c *gin.Context) {
	//时间聚合
	agg := elastic.NewDateHistogramAggregation().Field("created_at").CalendarInterval("day")

	//时间段搜索
	//今天开始到去年的今天
	now := time.Now()
	aYearAgo := now.AddDate(-1, 0, 0)

	format := "2006-01-02 15:04:05"

	query := elastic.NewRangeQuery("created_at").Gte(aYearAgo.Format(format)).Lte(now.Format(format))

	result, err := global.ESClient.Search(models.ArticleModel{}.Index()).Query(query).Aggregation("calendar", agg).Size(0).Do(context.Background())
	if err != nil {
		global.Log.Errorf("article calendar err:%+v\n", err.Error())
		res.FailWithMessage("查询失败", c)
		return
	}

	var data BucketsType

	_ = json.Unmarshal(result.Aggregations["canlendar"], &data)

	var resList = make([]CalendarResponse, 0)
	for _, bucket := range data.Buckets {
		time, _ := time.Parse(format, bucket.KeyAsString)
		DateCount[time.Format("2006-01-02")] = bucket.DocCount
	}

	days := int(now.Sub(aYearAgo).Hours() / 24)
	for i := 0; i <= days; i++ {
		day := aYearAgo.AddDate(0, 0, i).Format("2006-01-02")

		count, _ := DateCount[day]
		resList = append(resList, CalendarResponse{
			Date:  day,
			Count: count,
		})
	}
	res.OkWithData(resList, c)
}
