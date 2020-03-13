package main

import (
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	_ "github.com/ybinhome/envelope"
	"github.com/ybinhome/envelope/infra"
)

func main() {
	// 获取程序配置文件所在路径
	file := kvs.GetCurrentFilePath("config.ini", 1)
	// 加载并解析配置文件
	config := ini.NewIniFileConfigSource(file)
	// 构造启动管理器
	app := infra.New(config)
	// 启动应用程序
	app.Start()

	// 临时阻塞进程
	//c := make(chan int)
	//<-c
}
