package test_brun

import (
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"github.com/ybinhome/envelope/infra"
	"github.com/ybinhome/envelope/infra/base"
)

func init() {
	// 获取程序配置文件所在路径
	file := kvs.GetCurrentFilePath("../brun/config.ini", 1)
	// 加载并解析配置文件
	config := ini.NewIniFileConfigSource(file)

	// 注册基础资源
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})

	// 构造启动管理器
	app := infra.New(config)
	// 启动应用程序
	app.Start()
}
