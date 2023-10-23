package flag

import (
	"Blog_gin/core"
	"Blog_gin/global"
	sys_flag "flag"
	"github.com/fatih/structs"
)

type Option struct {
	DB   bool
	User string //  -u  admin 创建一个admin用户，  -u user
	Es   string //es
}

//go run app.go -u admin
//go run app.go -db

// Parse 解析命令
func Parse() Option {
	db := sys_flag.Bool("db", false, "初始化数据库")
	user := sys_flag.String("u", "", "创建用户")
	es := sys_flag.String("es", "", "es操作")
	//解析命令行参数写入注册得flag中
	sys_flag.Parse()
	return Option{
		DB:   *db,
		User: *user,
		Es:   *es,
	}
}

// IsWebStop  是否停止web项目
func IsWebStop(option Option) (f bool) {
	maps := structs.Map(&option)
	for _, v := range maps {
		switch val := v.(type) {
		case string:
			if val != "" {
				f = true
			}
		case bool:
			if val == true {
				f = true
			}
		}
	}
	return f
}

func SwitchOption(option Option) {
	if option.DB {
		Makemigrations()
	}
	if option.User == "admin" || option.User == "user" {
		CreateUser(option.User)
		return
	}

	if option.Es == "create" {
		//连接es
		global.ESClient = core.EsConnect()
		EsCreateIndex()
	}
	sys_flag.Usage()

}
