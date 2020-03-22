package base

import (
	"github.com/sirupsen/logrus"
	"github.com/ybinhome/envelope/infra"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

var callbacks []func()

func Register(fn func()) {
	callbacks = append(callbacks, fn)
}

type HookStarter struct {
	infra.BaseStarter
}

func (h *HookStarter) Init(ctx infra.StarterContext) {
	sigs := make(chan os.Signal)
	// 设置监听 3 和 15 两种信号
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		for {
			// 监听信号
			c := <-sigs
			// 输出日志，记录监听到的信号
			logrus.Info("notify: ", c)
			// 迭代 callbacks 中存储的 starter 的 Stop 方法清理资源函数
			for _, fn := range callbacks {
				fn()
			}
			// 清理完成后退出循环和程序
			break
			os.Exit(0)
		}
	}()
}

func (h *HookStarter) Start(ctx infra.StarterContext) {
	// 获取已注册的所有 starter
	starters := infra.GetStarters()

	//
	for _, s := range starters {
		typ := reflect.TypeOf(s)
		logrus.Infof("[ Register Notify Stop] : %s.Stop()", typ.String())
		Register(func() {
			s.Stop(ctx)
		})
	}
}
