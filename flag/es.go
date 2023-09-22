package flag

import "Blog_gin/models"

func EsCreateIndex() {
	models.FullTextModel{}.CreateIndex()
}
