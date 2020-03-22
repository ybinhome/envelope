package main

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/ybinhome/envelope/services"
	"net/rpc"
)

func main() {
	// 通过 rpc 包的 Dial 函数来创建连接到 rpc server 的 rpc client
	c, err := rpc.Dial("tcp", ":8082")
	if err != nil {
		logrus.Error(err)
	}

	sendout(c)
	receive(c)
}

func sendout(c *rpc.Client) {
	// 通过 rpc client 的 Call 方法来调用注册到 rpc server 中的 api，Call 方法的第一个参数是 serviceMethod，由 结构体名.方法名 组成，第二个参数为入参，第三个参数为出参
	in := services.RedEnvelopeSendingDTO{
		EnvelopeType: services.GeneralEnvelopeType,
		Username:     "测试用户",
		UserId:       "1ZTEBCNKIMBZPxEvh4BDYjEfSVU",
		Blessing:     "",
		Amount:       decimal.NewFromFloat(1),
		Quantity:     2,
	}
	out := &services.RedEnvelopeActivity{}
	err := c.Call("EnvelopeRpc.SendOut", in, out)
	if err != nil {
		logrus.Panic(err)
	}

	// 为了展示，将返回的数据打印出来
	logrus.Infof("%+v", out)
}

func receive(c *rpc.Client) {
	in := services.RedEnvelopeReceiveDTO{
		EnvelopeNo:   "1ZTEZUdnxTJsCe9kdkZEmILfjL4",
		RecvUsername: "测试用户1",
		RecvUserId:   "1ZTEZQl7j482rp93dgtRQ3eGSAP",
		AccountNo:    "",
	}
	out := &services.RedEnvelopeItemDTO{}

	err := c.Call("EnvelopeRpc.Receive", in, out)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("%+v", out)
}
