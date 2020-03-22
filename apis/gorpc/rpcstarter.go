package gorpc

import (
	"github.com/ybinhome/envelope/infra"
	"github.com/ybinhome/envelope/infra/base"
)

// 通过 init 方法我们无法保证向 rpc server 注册 api 时，rpc server 已经初始化完成，因此我们通过一个单独的 starter 来实现

type GoRpcApiStarter struct {
	infra.BaseStarter
}

func (g *GoRpcApiStarter) Init(ctx infra.StarterContext) {
	base.RpcRegister(new(EnvelopeRpc))
}
