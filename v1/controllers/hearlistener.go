package controllers

import (
	"github.com/astaxie/beego"
	"time"
	"ums/v1/models"
)

const (
	GC_INTERVAL      = 1  //(Minute)每隔一分钟清理一次超时user
	TIMEOUT_INTERVAL = 10 //(Minute)超时时间为10分钟
)

type UserListener struct {
	LastAliveTime time.Time
}

var Listener map[string]UserListener

func init() {
	Listener = make(map[string]UserListener)

	go func() {
		for {
			for k, v := range Listener {
				beego.Debug("key=", k, "v=", v)
				if time.Now().Sub(v.LastAliveTime) >= time.Duration(TIMEOUT_INTERVAL)*time.Minute {

					//stop redius??

					user := models.Userstatus{
						Usermac: k,
					}
					models.DelUserStatusByMac(&user)
					delete(Listener, k)
				}
			}
			beego.Debug("Listener Gc running...")
			time.Sleep(GC_INTERVAL * time.Minute)
		}
	}()
}
