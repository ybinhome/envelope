package infra

import "fmt"

type ConfigStarter struct {
	BaseStarter
}

func (c *ConfigStarter) Init(ctx StarterContext) {
	fmt.Println("配置初始化")
}

func (c *ConfigStarter) Setup(ctx StarterContext) {
	fmt.Println("配置初安装")
}

func (c *ConfigStarter) Start(ctx StarterContext) {
	fmt.Println("配置启动")
}

func init() {
	Register(&ConfigStarter{})
}
