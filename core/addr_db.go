package core

import (
	"Blog_gin/global"
	"fmt"
	geoip2db "github.com/cc14514/go-geoip2-db"
	"log"
)

func InitAddrDB() {
	db, err := geoip2db.NewGeoipDbByStatik()
	if err != nil {
		log.Fatal("ip地址数据库加载失败", err)
	}
	fmt.Println("mysql连接成功")
	global.AddrDB = db
}
