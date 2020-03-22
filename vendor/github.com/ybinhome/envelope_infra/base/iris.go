package base

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	irisRecover "github.com/kataras/iris/v12/middleware/recover"
	"github.com/sirupsen/logrus"
	"github.com/ybinhome/envelope_infra"
	"time"
)

var irisApplication *iris.Application

// iris applicaiton 的暴漏函数
func Iris() *iris.Application {
	Check(irisApplication)
	return irisApplication
}

type IrisServerStarter struct {
	infra.BaseStarter
}

func (i *IrisServerStarter) Init(ctx infra.StarterContext) {
	// 创建	iris application 实例
	irisApplication = initIris()
	// 日志组件的配置和扩展，iris 内部使用自己的 go log 组件，日志输出格式和 logrus 不一样，为了后续便于分析，我们统一日志输出格式采用 logrus
	//    获取 iris 的日志对象
	logger := irisApplication.Logger()
	//    使用 logger 的 Install 方法，安装 logrus 日志组件，将日志输出到 logrus 中
	logger.Install(logrus.StandardLogger())
}

func (i *IrisServerStarter) Start(ctx infra.StarterContext) {
	// 将注册的路由信息打印至控制台
	//    通过 Iris 的 GetRoutes 方法获取注册的路由信息
	routers := Iris().GetRoutes()
	for _, r := range routers {
		// 通过 iris route 的 Trace 方法获取路由信息
		logrus.Info(r.Trace())
	}

	// 启动 iris web 服务器
	//    获取 web 服务器端口
	port := ctx.Props().GetDefault("app.server.port", "18080")
	//    启动 iris web 服务器
	Iris().Run(iris.Addr(":" + port))
}

func (i *IrisServerStarter) StartBlocking() bool {
	return true
}

func initIris() *iris.Application {
	// 创建 iris 实例
	app := iris.New()
	// 设置 recovery 中间件，recover 中间件和 recover 函数冲突，因此我们重命名为 irisRecover
	app.Use(irisRecover.New())
	// 配置日志组件
	cfg := logger.Config{
		// 设置日志中记录的信息
		Status: true,
		IP:     true,
		Method: true,
		Path:   true,
		Query:  true,
		// 格式化日志输出
		LogFunc: func(endTime time.Time, latency time.Duration, status, ip, method, path string, message interface{}, headerMessage interface{}) {
			// 通过 iris 对象中的 Logger 方法的 Infof 来格式化输出
			app.Logger().Infof("| %s | %s | %s | %s | %s | %s | %s | %s |",
				endTime.Format("2006-01-02.15:04:05.000000"),
				latency.String(), status, ip, method, path, headerMessage, message,
			)
		},
	}
	// 将日志配置传递给日志中间件
	app.Use(logger.New(cfg))

	return app
}
