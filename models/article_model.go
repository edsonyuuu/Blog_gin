package models

import (
	"Blog_gin/global"
	"Blog_gin/models/ctype"
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type ArticleModel struct {
	ID        string `json:"id" structs:"id"`                 // es的id
	CreatedAt string `json:"created_at" structs:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at" structs:"updated_at"` // 更新时间

	Title    string `json:"title" structs:"title"`                // 文章标题
	Keyword  string `json:"keyword,omit(list)" structs:"keyword"` // 关键字
	Abstract string `json:"abstract" structs:"abstract"`          // 文章简介
	Content  string `json:"content,omit(list)" structs:"content"` // 文章内容

	LookCount     int `json:"look_count" structs:"look_count"`         // 浏览量
	CommentCount  int `json:"comment_count" structs:"comment_count"`   // 评论量
	DiggCount     int `json:"digg_count" structs:"digg_count"`         // 点赞量
	CollectsCount int `json:"collects_count" structs:"collects_count"` // 收藏量

	UserID       uint   `json:"user_id" structs:"user_id"`               // 用户id
	UserNickName string `json:"user_nick_name" structs:"user_nick_name"` //用户昵称
	UserAvatar   string `json:"user_avatar" structs:"user_avatar"`       // 用户头像

	Category string `json:"category" structs:"category"` // 文章分类
	Source   string `json:"source" structs:"source"`     // 文章来源
	Link     string `json:"link" structs:"link"`         // 原文链接

	BannerID  uint   `json:"banner_id" structs:"banner_id"`   // 文章封面id
	BannerUrl string `json:"banner_url" structs:"banner_url"` // 文章封面

	Tags ctype.Array `json:"tags" structs:"tags"` // 文章标签
}

func (ArticleModel) Index() string {
	return "article_index"
}

func (ArticleModel) Mapping() string {
	return `
{
  "settings": {
    "index":{
      "max_result_window": "100000"
    }
  }, 
  "mappings": {
    "properties": {
      "title": { 
        "type": "text"
      },
      "keyword": { 
        "type": "keyword"
      },
      "abstract": { 
        "type": "text"
      },
      "content": { 
        "type": "text"
      },
      "look_count": {
        "type": "integer"
      },
      "comment_count": {
        "type": "integer"
      },
      "digg_count": {
        "type": "integer"
      },
      "collects_count": {
        "type": "integer"
      },
      "user_id": {
        "type": "integer"
      },
      "user_nick_name": { 
        "type": "keyword"
      },
      "user_avatar": { 
        "type": "keyword"
      },
      "category": { 
        "type": "keyword"
      },
      "source": { 
        "type": "keyword"
      },
      "link": { 
        "type": "keyword"
      },
      "banner_id": {
        "type": "integer"
      },
      "banner_url": { 
        "type": "keyword"
      },
      "tags": { 
        "type": "keyword"
      },
      "created_at":{
        "type": "date",
        "null_value": "null",
        "format": "[yyyy-MM-dd HH:mm:ss]"
      },
      "updated_at":{
        "type": "date",
        "null_value": "null",
        "format": "[yyyy-MM-dd HH:mm:ss]"
      }
    }
  }
}
`
}

// 查看索引是否存在
func (a ArticleModel) IndexExists() bool {
	exists, err := global.ESClient.IndexExists(a.Index()).Do(context.Background())
	if err != nil {
		logrus.Error(err.Error())
		return exists
	}
	return exists
}

//创建索引

func (a ArticleModel) CreateIndex() error {
	if a.IndexExists() {
		//索引存在的情况
		a.RemoveIndex()
	}

	//没有索引的情况，创建索引
	createindex, err := global.ESClient.CreateIndex(a.Index()).BodyString(a.Mapping()).Do(context.Background())
	if err != nil {
		logrus.Error("创建索引失败")
		logrus.Error(err.Error())
		return err
	}

	if !createindex.Acknowledged {
		logrus.Error("创建索引失败")
		return err
	}
	logrus.Infof("索引 %s 创建成功", a.Index())
	return nil
}

func (a ArticleModel) RemoveIndex() error {
	logrus.Info("索引存在,删除索引~")
	//删除索引操作
	indexDelete, err := global.ESClient.DeleteIndex(a.Index()).Do(context.Background())
	if err != nil {
		logrus.Error("删除索引失败")
		logrus.Error(err.Error())
		return err
	}

	if !indexDelete.Acknowledged {
		logrus.Error("删除索引失败")
		return err
	}
	logrus.Info("索引删除成功")
	return nil
}

// 添加的方法
func (a *ArticleModel) Create() (err error) {
	indexResposne, err := global.ESClient.Index().Index(a.Index()).BodyJson(a).
		Refresh("true").Do(context.Background())

	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	a.ID = indexResposne.Id
	return nil
}

// 判断是否存在该文章
func (a ArticleModel) IsExistData() bool {
	res, err := global.ESClient.Search(a.Index()).
		Query(elastic.NewTermQuery("keyword", a.Title)).
		Size(1).
		Do(context.Background())
	if err != nil {
		logrus.Error(err.Error())
		return false
	}
	if res.Hits.TotalHits.Value > 0 {
		return true
	}
	return false
}

func (a *ArticleModel) GetDataByID(id string) error {
	res, err := global.ESClient.Get().Index(a.Index()).Id(id).Do(context.Background())

	if err != nil {
		return err
	}
	err = json.Unmarshal(res.Source, a)
	return err
}

//方法内部的代码逻辑如下：
//
//global.ESClient.Get().Index(a.Index()).Id(id).Do(context.Background())：
//这是一个使用Elasticsearch客户端的操作，它执行了一个GET请求，通过指定索引和文档ID来获取数据。
//global.ESClient是一个全局的Elasticsearch客户端实例，.Get()表示执行GET请求，
//.Index(a.Index())指定了索引，.Id(id)指定了文档ID，最后.Do(context.Background())执行请求并返回结果。

// res, err := ...：将请求的结果赋值给res变量，并将可能出现的错误赋值给err变量。
//
// if err != nil { return err }：如果在请求过程中出现了错误，直接返回该错误。
//
// err = json.Unmarshal(res.Source, a)：
//将获取到的数据解析为JSON格式，并将解析后的结果赋值给a，即ArticleModel结构体实例。这里使用了json.Unmarshal函数进行解析。
// return err：返回可能出现的解析错误。
