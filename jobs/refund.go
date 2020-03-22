package jobs

import (
	"fmt"
	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	"github.com/ybinhome/envelope/core/envelopes"
	"github.com/ybinhome/envelope/infra"
	"time"
)

type RefundExpiredJobStarter struct {
	infra.BaseStarter
	ticker *time.Ticker
	mutex  *redsync.Mutex
}

func (r *RefundExpiredJobStarter) Init(ctx infra.StarterContext) {
	d := ctx.Props().GetDurationDefault("jobs.refund.interval", time.Minute)
	r.ticker = time.NewTicker(d)

	// 使用 redsync 分布式互斥锁 - 1. 构建连接池
	maxIdle := ctx.Props().GetIntDefault("redis.maxIdle", 2)
	maxActive := ctx.Props().GetIntDefault("redis.maxActive", 5)
	idleTimeout := ctx.Props().GetDurationDefault("redis.idleTimeout", 20*time.Second)
	addr := ctx.Props().GetDefault("redis.addr", "127.0.0.1:6379")
	pools := make([]redsync.Pool, 0)
	pool := &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial("tcp", addr)
		},
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
	}
	pools = append(pools, pool)

	// 使用 redsync 分布式互斥锁 - 2. 创建 redsync 对象
	rsync := redsync.New(pools)

	// 使用 redsync 分布式互斥锁 - 3. 创建互斥锁
	// -  获取节点 ip，用于展示那个节点获得了分布式互斥锁
	ip, err := utils.GetExternalIP()
	if err != nil {
		ip = "127.0.0.1"
	}
	// - 获取分布式互斥锁
	r.mutex = rsync.NewMutex(
		// 设置锁名称
		"lock:RefundExpired",
		// 设置锁过期时间
		redsync.SetExpiry(50*time.Second),
		// 设置重试次数
		redsync.SetRetryDelay(3),
		// 设置互斥锁的值，格式为 时间: ip 地址
		redsync.SetGenValueFunc(func() (s string, err error) {
			now := time.Now()
			logrus.Infof("节点 %s 正在执行过期红包退款任务", ip)
			return fmt.Sprintf("%d: %s", now.Unix(), ip), nil
		}),
	)
}

func (r *RefundExpiredJobStarter) Start(ctx infra.StarterContext) {
	go func() {
		for {
			c := <-r.ticker.C

			err := r.mutex.Lock()
			if err == nil {
				logrus.Debug("过期红包退款开始...", c)
				// 红包业务逻辑代码
				domain := envelopes.ExpiredEnvelopeDomain{}
				domain.Expired()
			} else {
				logrus.Info("已经有节点在运行该任务了")
			}
			r.mutex.Unlock()
		}
	}()
}

func (r *RefundExpiredJobStarter) Stop(ctx infra.StarterContext) {

}
