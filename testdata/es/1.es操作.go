package main

import (
	"Blog_gin/core"
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

func main() {

}

var ES *elastic.Client

func init() {
	core.InitConf()
	core.InitLogger()
	host := "http://127.0.0.1:9200"
	var err error
	ES, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetBasicAuth("", ""))
	if err != nil {
		fmt.Println("es err:", err)
	}
}

type DemoModel struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Index     string `json:"indexs"`
	Content   string `json:"content"`
	UserId    uint   `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

func (DemoModel) CreateIndex() {
	//idex := DemoModel{}.Index
	//查索引是否存在

}

func (demo DemoModel) IndexExists() bool {
	exist, err := ES.IndexExists(demo.Index).Do(context.Background())
	if err != nil {
		return exist
	}
	return exist

}
