package article_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/models/ctype"
	"Blog_gin/res"
	"Blog_gin/service/es_ser"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"time"
)

type ArticleUpdateRequest struct {
	Title    string   `json:"title"`     // 文章标题
	Abstract string   `json:"abstract"`  // 文章简介
	Content  string   `json:"content"`   // 文章内容
	Category string   `json:"category"`  // 文章分类
	Source   string   `json:"source"`    // 文章来源
	Link     string   `json:"link"`      // 原文链接
	BannerID uint     `json:"banner_id"` // 文章封面id
	Tags     []string `json:"tags"`      // 文章标签
	ID       string   `json:"id"`
}

// ArticleUpdateView 修改文章
// @Summary 修改文章
// @Description 修改文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param	data	query ArticleUpdateRequest	false	"请求参数"
// @Success 200 {object} res.Response{}
// @Router /api/articles [put]
func (ArticleApi) ArticleUpdateView(c *gin.Context) {
	var cr ArticleUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		global.Log.Error(err)
		res.FailWithError(err, &cr, c)
		return
	}

	var bannerUrl string
	if cr.BannerID != 0 {
		err = global.DB.Model(models.BannerModel{}).Where("id = ? ", cr.BannerID).Select("path").Scan(&bannerUrl).Error
		if err != nil {
			res.FailWithMessage("banner不存在", c)
			return
		}
	}
	article := models.ArticleModel{
		UpdatedAt: time.Now().Format("2006-01-02 15:04-05"),
		Title:     cr.Title,
		Keyword:   cr.Title,
		Abstract:  cr.Abstract,
		Content:   cr.Content,
		Category:  cr.Category,
		Source:    cr.Source,
		Link:      cr.Link,
		BannerID:  cr.BannerID,
		BannerUrl: bannerUrl,
		Tags:      cr.Tags,
	}
	maps := structs.Map(&article)
	var dataMap = map[string]any{}

	//去掉空值
	for key, value := range maps {
		switch val := value.(type) {
		case string:
			if val == "" {
				continue
			}
		case int:
			if val == 0 {
				continue
			}

		case uint:
			if val == 0 {
				continue
			}
		case ctype.Array:
			if len(val) == 0 {
				continue
			}

		case []string:
			if len(val) == 0 {
				continue
			}

		}
		dataMap[key] = value
	}

	err = article.GetDataByID(cr.ID)
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("文章不存在", c)
		return
	}

	err = es_ser.ArticleUpdate(cr.ID, dataMap)
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("文章更新失败", c)
		return
	}

	//更新成功，同步到全文搜索中
	newArticle, _ := es_ser.CommDetail(cr.ID)
	if article.Content != newArticle.Content || article.Title != newArticle.Title {
		es_ser.DeleteFullTextByArticleID(cr.ID)
		es_ser.AsyncArticleByFullText(cr.ID, article.Title, newArticle.Content)
	}

	res.OkWithMessage("更新成功", c)

}
