package main

import (
	"fmt"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"time"
)

func main() {
	// 使用 props 的 kvs 包 的 GetCurrentFilePath 方法修改文件查找的路径层级，默认路径是项目根目录
	file := kvs.GetCurrentFilePath("config.ini", 1)

	// 读取配置文件
	config := ini.NewIniFileConfigSource(file)

	// 使用 props ini 包中的 Get 类型方法读取配置项，方法名称带 Default 时，支持设置默认值
	fmt.Println(config.GetIntDefault("app.server.port", 18080))
	fmt.Println(config.GetDefault("app.name", "unknown"))
	fmt.Println(config.GetBoolDefault("app.enabled", false))
	fmt.Println(config.GetDurationDefault("app.time", time.Second))
}
