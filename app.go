package main

import (
	core2 "Blog_gin/core"
	_ "Blog_gin/docs"
	"Blog_gin/flag"
	"Blog_gin/global"
	"Blog_gin/routers"
)

func main() {
	//读取配置文件
	core2.InitConf()
	//初始化日志
	global.Log = core2.InitLogger()
	//连接数据库
	global.DB = core2.InitGorm()

	core2.InitAddrDB()
	defer global.AddrDB.Close()
	//命令行参数绑定
	option := flag.Parse()
	if flag.IsWebStop(option) {
		flag.SwitchOption(option)
		return
	}
	//连接redis
	global.Redis = core2.ConnectRedis()
	//连接es
	//global.ESClient = core2.EsConnect()

	router := routers.InitRouter()

	addr := global.Config.System.Addr()
	global.Log.Infof("gvb_server运行在：%s", addr)
	err := router.Run(addr)
	if err != nil {
		global.Log.Fatalf(err.Error())
	}
}
