package base

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"github.com/ybinhome/envelope/infra"
	"sync"
)

// 配置文件会贯穿程序运行全程，因此需要对外暴漏，props kvs.ConfigSource 会在程序启动时进行初始化

var props kvs.ConfigSource

func Props() kvs.ConfigSource {
	return props
}

// 实现 props 的 starter

type PropsStarter struct {
	infra.BaseStarter
}

func (p *PropsStarter) Init(ctx infra.StarterContext) {
	props = ctx.Props()
	logrus.Info("初始化配置")
	// 读取红包系统账户
	GetSystemAccount()
}

// 读取一次系统红包账户的配置信息
type SystemAccount struct {
	AccountNo   string
	AccountName string
	UserId      string
	Username    string
}

var systemAccount *SystemAccount
var systemAccountOnce sync.Once

func GetSystemAccount() *SystemAccount {
	systemAccountOnce.Do(func() {
		systemAccount = new(SystemAccount)
		err := kvs.Unmarshal(Props(), systemAccount, "system.account")
		if err != nil {
			panic(err)
		}
	})
	return systemAccount
}

// 读取红包活动相关配置
func GetEnvelopeActivityLink() string {
	link := Props().GetDefault("envelope.link", "/v1/envelope/link")
	return link
}
func GetEnvelopeDomain() string {
	domain := Props().GetDefault("envelope.domain", "http://localhost")
	return domain
}
