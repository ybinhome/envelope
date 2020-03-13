package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func main() {
	app := iris.Default()
	app.Get("/hello", func(ctx iris.Context) {
		// 抛出异常，测试 OnAnyErrorCode
		panic("hello 出错了")
		ctx.WriteString("hello world! iris")
	})

	// 通过 Party 方法，将 users 和 orders 划分到 v1 分组中
	v1 := app.Party("/v1")

	// 分组单独使用中间件，调用中间件函数中需要单独使用 context.Next() 函数执行，否则会卡在使用中间件处
	v1.Use(func(context iris.Context) {
		logrus.Info("自定义中间件")
		context.Next()
	})

	// 参数化 path 匹配整型，限制最小值为 2
	v1.Get("/users/{id:uint64 min(2)}",
		func(ctx iris.Context) {
			id := ctx.Params().GetUint64Default("id", 0)
			ctx.WriteString(strconv.Itoa(int(id)))
		})

	// 参数化 path 匹配 string 类型，匹配前缀为 a_
	v1.Get("/orders/{action:string prefix(a_)}", func(ctx iris.Context) {
		a := ctx.Params().Get("action")
		ctx.WriteString(a)
	})

	// 定义全局错误信息
	app.OnAnyErrorCode(func(context iris.Context) {
		context.WriteString("看起来服务器出错了！")
	})

	// 定义特定响应码的错误信息
	app.OnErrorCode(http.StatusNotFound, func(context iris.Context) {
		context.WriteString("访问路径不存在。")
	})

	err := app.Run(iris.Addr(":8082"))
	fmt.Println(err)
}
