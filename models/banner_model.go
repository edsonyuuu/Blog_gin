package models

import (
	"Blog_gin/global"
	"Blog_gin/models/ctype"
	"gorm.io/gorm"
	"os"
)

type BannerModel struct {
	MODEL
	Path      string          `json:"path"`                        // 图片路径
	Hash      string          `json:"hash"`                        // 图片的hash值，用于判断重复图片
	Name      string          `gorm:"size:38" json:"name"`         // 图片名称
	ImageType ctype.ImageType `gorm:"default:1" json:"image_type"` // 图片的类型， 本地还是七牛
}

// BeforeDelete 在删除操作之前的部分操作
func (b *BannerModel) BeforeDelete(tx *gorm.DB) (err error) {
	if b.ImageType == ctype.Local {
		//代表删除本地图片，还要删除本地的存储
		err := os.Remove(b.Path[0:])
		if err != nil {
			global.Log.Error(err)
			return err
		}
	}
	return nil
}
