package gorpc

import (
	"github.com/ybinhome/envelope/services"
)

type EnvelopeRpc struct {
}

// GoRPC 接口规范：
// 1. 入参和出参都要作为方法的参数；
// 2. 方法必须有两个参数，并且是可导出类型；
// 3. 第二个参数 (返回值) 必须是指针类型；
// 4. 方法返回值要返回 error 类型；
// 5. 方法必须是可导出的；

func (e *EnvelopeRpc) SendOut(in services.RedEnvelopeSendingDTO, out *services.RedEnvelopeActivity) error {
	s := services.GetRedEnvelopeService()
	a1, err := s.SendOut(in)
	// 对于出参的引用，只能拷贝，不能修改，修改会报错
	a1.CopyTo(out)
	return err
}

func (e *EnvelopeRpc) Receive(in services.RedEnvelopeReceiveDTO, out *services.RedEnvelopeItemDTO) error {
	s := services.GetRedEnvelopeService()
	a, err := s.Receive(in)
	// 对于出参的引用，只能拷贝，不能修改，修改会报错
	a.CopeTo(out)
	return err
}

// 通过 rpc client 访问 rpc server 中注册的 api 见 example/rpcClient
