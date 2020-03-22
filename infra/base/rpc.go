package base

import (
	"github.com/sirupsen/logrus"
	"github.com/ybinhome/envelope/infra"
	"net"
	"net/rpc"
	"reflect"
)

// 为了能够将服务注册到 rpc server，从而使 rpc server 来处理这些请求，需要将创建的 rpc server 暴露出去
var rpcServer *rpc.Server

// rpc server 获取函数
func RpcServer() *rpc.Server {
	Check(rpcServer)
	return rpcServer
}

// api 向 rpc server 注册函数
func RpcRegister(ri interface{}) {
	typ := reflect.TypeOf(ri)
	logrus.Infof("goRPC Register: %s", typ.String())
	RpcServer().Register(ri)
}

type GoRPCStarter struct {
	infra.BaseStarter
	server *rpc.Server
}

func (g *GoRPCStarter) Init(ctx infra.StarterContext) {
	// 将初始化完成的 rpc server 赋值给 rpcServer 变量，从而通过 RpcServer() 将其暴露出去
	g.server = rpc.NewServer()
	rpcServer = g.server
}

func (g *GoRPCStarter) Start(ctx infra.StarterContext) {
	port := ctx.Props().GetDefault("app.rpc.port", "8082")

	// 监听网络端口
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Info("tcp port listened for rpc:", port)

	// rpc server 使用监听段鸥处理网络连接和请求
	go g.server.Accept(listener)
}
