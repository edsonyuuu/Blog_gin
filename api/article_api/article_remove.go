package article_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/es_ser"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

type IDListRequest struct {
	IDList []string `json:"id_list"`
}

// ArticleRemoveView 删除文章
// @Summary 删除文章
// @Description 删除文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param	data	query IDListRequest	false	"请求参数"
// @Success 200 {object} res.Response{}
// @Router /api/articles [delete]
func (ArticleApi) ArticleRemoveView(c *gin.Context) {
	var cr IDListRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		global.Log.Error(err)
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	//文章删除了，顺带把这个收藏的文章也给删除了
	//用户收藏表新增一个字段表示收藏表
	bulkService := global.ESClient.Bulk().Index(models.ArticleModel{}.Index()).Refresh("true")

	for _, id := range cr.IDList {
		req := elastic.NewBulkDeleteRequest().Id(id)
		bulkService.Add(req)
		go es_ser.DeleteFullTextByArticleID(id)
	}

	result, err := bulkService.Do(context.Background())
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("删除失败", c)
		return
	}

	res.OkWithMessage(fmt.Sprintf("成功删除 %d 篇文章", len(result.Succeeded())), c)

}
