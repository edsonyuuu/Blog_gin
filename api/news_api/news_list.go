package news_api

import (
	"Blog_gin/res"
	"Blog_gin/service/redis_ser"
	"Blog_gin/utils/requests"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type Params struct {
	ID   string `json:"id" form:"id"`
	Size int    `json:"size" form:"size"`
}

type Header struct {
	SignatureKey string `form:"signaturekey" structs:"signaturekey"`
	Version      string `form:"version" structs:"version"`
	UserAgent    string `form:"User-Agent" structs:"User-Agent"`
}

type NewResponse struct {
	Code int                 `json:"code"`
	Data []redis_ser.NewData `json:"data"`
	Msg  string              `json:"msg"`
}

const newAPI = "https://api.codelife.cc/api/top/list"
const timeout = 2 * time.Second

// NewsListView
// @Summary 获取新闻列表
// @Description 获取新闻列表
// @Tags 新闻管理
// @Accept json
// @Produce json
// @Param ID formData string false   "ID"
// @Param Size formData int false   "Size"
// @Param signaturekey formData string false "签名密钥"
// @Param version formData string false "版本号"
// @Param User-Agent formData string false "User-Agent"
// @Success 200 {object} res.Response
// @Router /api/news/list [post]
func (NewsApi) NewsListView(c *gin.Context) {

	var cr Params
	var headers Header
	err := c.ShouldBindJSON(&cr)
	err = c.ShouldBindJSON(&headers)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}
	if cr.Size == 0 {
		cr.Size = 1
	}

	key := fmt.Sprintf("%s-%d", cr.ID, cr.Size)
	newsData, _ := redis_ser.GetNews(key)
	if len(newsData) != 0 {
		res.OkWithData(newsData, c)
		return
	}
	httpResponse, err := requests.Post(newAPI, cr, structs.Map(headers), timeout)

	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}
	var response NewResponse
	byteData, err := io.ReadAll(httpResponse.Body)
	err = json.Unmarshal(byteData, &response)
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}
	if response.Code != 200 {
		res.FailWithMessage(response.Msg, c)
		return
	}

	res.OkWithData(response.Data, c)
	redis_ser.SetNews(key, response.Data)
	return

}
