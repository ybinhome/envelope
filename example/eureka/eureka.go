package main

import (
	"github.com/tietang/go-eureka-client/eureka"
	"github.com/tietang/props/ini"
)

func main() {
	// 载入配置文件
	config := ini.NewIniFileConfigSource("example/eureka/config.ini")

	// 创建 eureka 客户端
	client := eureka.NewClient(config)

	// 启动 eureka 客户端
	client.Start()

	// 阻塞 eureka 客户端
	c := make(chan int, 1)
	<-c
}
