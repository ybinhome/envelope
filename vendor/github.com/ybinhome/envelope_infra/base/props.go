package base

import (
	"fmt"
	"github.com/tietang/props/kvs"
	"github.com/ybinhome/envelope_infra"
)

// 配置文件会贯穿程序运行全程，因此需要对外暴漏，props kvs.ConfigSource 会在程序启动时进行初始化

var props kvs.ConfigSource

func Props() kvs.ConfigSource {
	return props
}

// 实现 props 的 starter

type PropsStarter struct {
	infra.BaseStarter
}

func (p *PropsStarter) Init(ctx infra.StarterContext) {
	props = ctx.Props()
	fmt.Println("初始化配置")
}
