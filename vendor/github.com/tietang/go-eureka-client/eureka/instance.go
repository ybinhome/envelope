package eureka

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
)

func (c *Client) SetCurrentInstanceInfo(ins *InstanceInfo) {
	c.InstanceInfo = ins
}
func (c *Client) Start() {

	go c.run()
	go c.hook()
}

func (c *Client) run() {
	ins := c.InstanceInfo
	if ins == nil {
		log.Error("no instance info")
		return
	}
	timer := time.NewTicker(10 * time.Second)
	isInit := false
	lastFailBeatSeconds := 0
	lastFailBeatTimes := 0
	for {
		select {
		case <-timer.C:
			go func() {
				applications, errForGet := c.GetApplications() // Retrieves all applications from eureka server(s)
				if errForGet == nil {
					c.Applications = applications
				}
				if !isInit {
					errForReg := c.RegisterInstance(ins.AppName, ins) // Register new instance in your eureka(s)
					if errForReg == nil {
						isInit = true
					}
					err := c.UpdateStatus(ins.AppName, ins.InstanceId, StatusUp)
					if err == nil {
						isInit = true
					}
				} else {
					errForBeat := c.SendHeartbeat(ins.App, ins.InstanceId) // say to eureka that your app is alive (here you must send heartbeat before 30 sec)
					if errForBeat != nil {
						lastFailBeatTimes++
						nowSeconds := time.Now().Second()
						if lastFailBeatSeconds == 0 {
							lastFailBeatSeconds = nowSeconds
						}
						if nowSeconds-lastFailBeatSeconds >= 30 || lastFailBeatTimes >= 3 {
							isInit = false
						}
					}
				}
			}()
		}
	}
}

func (client *Client) hook() {
	hook := utils.NewHook()
	handler := func(s os.Signal, arg interface{}) {
		log.WithFields(log.Fields{
			"signal": s,
		}).Info("handle signal: ")
		client.UnregisterInstance(client.InstanceInfo.AppName, client.InstanceInfo.InstanceId)

		os.Exit(1)
	}
	//Interrupt Signal = syscall.SIGINT interrupt
	//Kill      Signal = syscall.SIGKILL killed
	//syscall.SIGTERM terminated

	hook.Register(os.Interrupt, handler)
	hook.Register(os.Kill, handler)
	hook.Register(syscall.SIGTERM, handler)

	for {
		c := make(chan os.Signal)
		signal.Notify(c)
		sig := <-c

		err := hook.Handle(sig, nil)
		if err != nil {
			log.WithFields(log.Fields{
				"signal": sig,
			}).Info("unknown signal received: : ")
			//			os.Exit(1)
		}
	}
}
