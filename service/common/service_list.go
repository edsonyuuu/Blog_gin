package common

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"fmt"
	"gorm.io/gorm"
)

type Option struct {
	models.PageInfo
	Debug bool
	Likes []string
}

// ComList 每页限制数、偏移量、排序方式，并将查询结果存储到 list 变量中。
func ComList[T any](model T, option Option) (list []T, count int64, err error) {
	DB := global.DB
	if option.Debug {
		DB = global.DB.Session(&gorm.Session{Logger: global.MysqlLog})
	}

	if option.Sort == "" {
		option.Sort = "created_at desc" //默认按照时间往前排
	}
	DB = DB.Where(model)
	for index, column := range option.Likes {
		if index == 0 {
			DB.Where(fmt.Sprintf("%s like ?", column), fmt.Sprintf("%%%s%%", option.Key))
			continue
		}
		DB.Or(fmt.Sprintf("%s like ?", column), fmt.Sprintf("%%%s%%", option.Key))
	}
	count = DB.Find(&list).RowsAffected
	//这里的query会受到上边
	query := DB.Where(model)
	offset := (option.Page - 1) * option.Limit
	if offset < 0 {
		offset = 0
	}

	err = query.Limit(option.Limit).Offset(offset).Order(option.Sort).Find(&list).Error
	return list, count, err

}
