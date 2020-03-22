package infra

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"reflect"
)

// 应用程序启动管理器
type BootApplication struct {
	config         kvs.ConfigSource
	starterContext StarterContext
}

func New(config kvs.ConfigSource) *BootApplication {
	b := &BootApplication{
		config:         config,
		starterContext: StarterContext{},
	}

	b.starterContext[KeyProps] = config
	return b
}

func (b *BootApplication) Start() {
	// 1. 初始化 starter
	b.init()
	// 2. 安装 starter
	b.setup()
	// 3. 启动 starter
	b.start()
}

func (b *BootApplication) init() {
	log.Info("Initializing starters...")
	for _, starter := range StarterRegister.AllStarters() {
		v := reflect.TypeOf(starter)
		log.Debug("Initializing: type=%s", v.String())
		starter.Init(b.starterContext)
	}
}

func (b *BootApplication) setup() {
	log.Info("Setup starters...")
	for _, starter := range StarterRegister.AllStarters() {
		starter.Setup(b.starterContext)
	}
}

func (b *BootApplication) start() {
	log.Info("Starting starters...")
	for i, starter := range StarterRegister.AllStarters() {
		if starter.StartBlocking() {
			// 如果可阻塞的 starter 是最后一个 starter，直接启动并阻塞
			if i+1 == len(StarterRegister.AllStarters()) {
				starter.Start(b.starterContext)
			} else { // 可阻塞的 starter 不是最后一个 starter，使用 goroutine 来异步启动，防止阻塞后面的 starter
				go starter.Start(b.starterContext)
			}
		} else {
			starter.Start(b.starterContext)
		}
	}
}
