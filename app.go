package envelope

import (
	"github.com/ybinhome/envelope/infra"
	"github.com/ybinhome/envelope/infra/base"
)

// 通过统一的 init 函数来手动注册，从而实现手动管理程序启动顺序，防止乱序
// 如果程序对启动顺序无感知，可以在 starter 文件中直接使用 init 函数初始化

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&base.IrisServerStarter{})
}
