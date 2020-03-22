package base

import (
	"github.com/kataras/iris/v12"
	"github.com/tietang/go-eureka-client/eureka"
	"github.com/ybinhome/envelope/infra"
	"time"
)

type EurekaStarter struct {
	infra.BaseStarter
	client *eureka.Client
}

func (e *EurekaStarter) Init(ctx infra.StarterContext) {
	// 通过配置文件中的 eureka 相关配置项实例化 eureka client
	e.client = eureka.NewClient(ctx.Props())
	// 启动 eureka client
	e.client.Start()
}

func (e *EurekaStarter) Setup(ctx infra.StarterContext) {
	// 配置服务在 eureka 中的 /info 页面信息，展示服务名称和注册时间
	info := make(map[string]interface{})
	info["startTime"] = time.Now()
	info["appName"] = ctx.Props().GetDefault("app.name", "envelope")
	Iris().Get("/info", func(context iris.Context) {
		context.JSON(info)
	})

	// 配置服务在 eureka 中的 /health 页面信息，展示 eureka client 的将康状态信息
	Iris().Get("/health", func(context iris.Context) {
		health := eureka.Health{
			Details: make(map[string]interface{}),
		}
		health.Status = eureka.StatusUp
		context.JSON(health)
	})
}
