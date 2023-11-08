package article_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/es_ser"
	"Blog_gin/utils/jwts"
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
	"strings"
	"time"
)

type ArticleRequest struct {
	Title    string   `json:"title" binding:"required" msg:"文章标题必填"`   // 文章标题
	Abstract string   `json:"abstract"`                                // 文章简介
	Content  string   `json:"content" binding:"required" msg:"文章内容必填"` // 文章内容
	Category string   `json:"category"`                                // 文章分类
	Source   string   `json:"source"`                                  // 文章来源
	Link     string   `json:"link"`                                    // 原文链接
	BannerID uint     `json:"banner_id"`                               // 文章封面id
	Tags     []string `json:"tags"`                                    // 文章标签
}

// ArticleCreateView   创建文章
// @Summary 创建文章
// @Description 创建文章
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param articleRequest body ArticleRequest true "文章请求参数"
// @Success 200 {object} res.Response{}
// @Router /api/articles [post]
func (ArticleApi) ArticleCreateView(c *gin.Context) {
	var cr ArticleRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, &cr, c)
		return
	}

	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)
	userId := claims.UserID
	userNickName := claims.NickName

	//处理内容content
	//解析处理Markdown内容的代码
	unsafe := blackfriday.Run([]byte(cr.Content))
	//判断是否存在标签
	//解析HTML文档，创建一个Document请求对象
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(unsafe)))
	fmt.Println("打印doc.text:", doc.Text())
	nodes := doc.Find("script").Nodes
	if len(nodes) > 0 {
		//存在标签
		doc.Find("script").Remove()
		converter := md.NewConverter("", true, nil)
		html, _ := doc.Html()
		markdown, _ := converter.ConvertString(html)
		cr.Content = markdown

	}
	//标题为空
	if cr.Abstract == "" {
		//截取汉字
		abs := []rune(doc.Text())
		//将内容转化为html，并且过滤掉xss，以及获取中文内容
		if len(abs) > 100 {
			cr.Abstract = string(abs[:100])
		} else {
			cr.Abstract = string(abs)
		}
	}

	//不传bannerId后台会传一张
	var bannerUrl string
	err = global.DB.Model(models.BannerModel{}).Where("id = ?", cr.BannerID).Select("path").Scan(&bannerUrl).Error
	if err != nil {
		res.FailWithMessage("banner不存在", c)
		return
	}

	//查用户头像
	var avatar string
	err = global.DB.Model(models.UserModel{}).Where("id = ?", userId).Select("avatar").Scan(&avatar).Error
	if err != nil {
		res.FailWithMessage("用户不存在", c)
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	article := models.ArticleModel{
		CreatedAt:    now,
		UpdatedAt:    now,
		Title:        cr.Title,
		Abstract:     cr.Abstract,
		Content:      cr.Abstract,
		UserID:       userId,
		UserNickName: userNickName,
		UserAvatar:   avatar,
		Category:     cr.Category,
		Source:       cr.Source,
		Link:         cr.Link,
		BannerID:     cr.BannerID,
		BannerUrl:    bannerUrl,
		Tags:         cr.Tags,
	}

	if article.IsExistData() {
		res.FailWithMessage("该文章已存在", c)
		return
	}

	err = article.Create()
	if err != nil {
		global.Log.Errorf("创建article表失败 err:%+v\n", err.Error())
		res.FailWithMessage("创建article表失败", c)
		return
	}

	//起协程同步数据到全文搜索中
	go es_ser.AsyncArticleByFullText(article.ID, article.Title, article.Content)
	res.OkWithMessage("文章发布成功", c)
}
