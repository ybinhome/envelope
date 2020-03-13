package infra

import "github.com/tietang/props/kvs"

const (
	KeyProps = "_config"
)

// 基础资源上下文结构体
type StarterContext map[string]interface{}

func (s StarterContext) Props() kvs.ConfigSource {
	p := s[KeyProps]
	if p == nil {
		panic("配置还没有初始化")
	}
	return p.(kvs.ConfigSource)
}

// 基础资源启动器接口
type Starter interface {
	// 1. 系统启动，初始化一些基础资源
	Init(StarterContext)
	// 2. 系统基础资源的安装
	Setup(StarterContext)
	// 3. 启动基础资源
	Start(StarterContext)
	// 4. 资源停止和销毁
	Stop(StarterContext)
	// 启动器是否可阻塞
	StartBlocking() bool
}

// 不是所有基础资源都需要实现 Starter 接口的 5 种方法，例如：配置资源只需要实现 Init 方法即可，而数据库资源需要实现如上 5 个方法，因此，可以定义一个空实现，所有的基础资源嵌套此空实现，从而不同的基础资源按需
//   实现需要的方法即可，无需全部实现。

// 基础空启动器实现 - 定义一个忽略的变量，验证 BaseStarter 是否实现了 Start 接口的所有方法，无报错即实现了所有方法
var _ Starter = new(BaseStarter)

// 基础空启动器实现 - 结构体
type BaseStarter struct {
}

// 基础空启动器实现 - 方法实现
func (b *BaseStarter) Init(ctx StarterContext)  {}
func (b *BaseStarter) Setup(ctx StarterContext) {}
func (b *BaseStarter) Start(ctx StarterContext) {}
func (b *BaseStarter) Stop(ctx StarterContext)  {}
func (b *BaseStarter) StartBlocking() bool      { return false }

// 上面定义了基础资源的生命周期，还需要一个注册器，将所有基础资源管理起来

// 启动注册器 - 结构器 (全局仅需要一个注册器，因此首字母小写即可，无需导出)
type starterRegister struct {
	starters []Starter
}

// 启动注册器 - 注册方法实现
func (r *starterRegister) Register(s Starter) {
	r.starters = append(r.starters, s)
}

// 启动注册器 - 返回所有基础资源实现
func (r *starterRegister) AllStarters() []Starter {
	return r.starters
}

var StarterRegister *starterRegister = new(starterRegister)

func Register(s Starter) {
	StarterRegister.Register(s)
}

// 系统基础资源的启动管理
func SystemRun() {
	ctx := StarterContext{}
	// 1. 初始化基础资源
	for _, starter := range StarterRegister.AllStarters() {
		starter.Init(ctx)
	}

	// 2. 安装基础资源
	for _, starter := range StarterRegister.AllStarters() {
		starter.Setup(ctx)
	}

	// 3. 启动基础资源
	for _, starter := range StarterRegister.AllStarters() {
		starter.Start(ctx)
	}
}
