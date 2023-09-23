package article_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/es_ser"
	"Blog_gin/service/redis_ser"
	"Blog_gin/utils/jwts"
	"github.com/gin-gonic/gin"
	"github.com/liu-cn/json-filter/filter"
	"github.com/olivere/elastic/v7"
)

type ArticleSearchRequest struct {
	models.PageInfo
	Tag    string `json:"tag" form:"tag"`
	IsUser bool   `json:"is_user" form:"is_user"`
}

// ArticleListView 文章列表
// @Tags 文章管理
// @Summary 文章列表
// @Description 文章列表
// @Param data query ArticleSearchRequest    false  "表示多个参数"
// @Param token header string    false  "token"
// @Produce json
// @Success 200 {object} res.Response{data=res.ListResponse[models.ArticleModel]}
// @Router /api/articles [get]
func (ArticleApi) ArticleListView(c *gin.Context) {
	var cr ArticleSearchRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}
	boolSearch := elastic.NewBoolQuery()

	if cr.IsUser {
		token := c.GetHeader("token")
		claims, err := jwts.ParseToken(token)
		if err == nil && !redis_ser.CheckLogout(token) {
			boolSearch.Must(elastic.NewTermQuery("user_id", claims.UserID))
		}
	}

	list, count, err := es_ser.CommList(es_ser.Option{
		PageInfo: cr.PageInfo,
		Fields:   []string{"title", "content", "category"},
		Tag:      cr.Tag,
		Query:    boolSearch,
	})
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("查询失败", c)
		return
	}

	//json-filter  空值问题
	data := filter.Omit("list", list)
	_list, _ := data.(filter.Filter)
	if string(_list.MustMarshalJSON()) == "{}" {
		res.OkWithList(list, int64(count), c)
		return
	}

	res.OkWithList(data, int64(count), c)

}
