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

var listener map[string]time.Time

func run() {
	for {
		for k, v := range listener {
			beego.Debug("key=", k, "v=", v)
			
			if time.Now().Sub(v) >= time.Duration(TIMEOUT_INTERVAL)*time.Minute {
				//stop redius??
				user := &models.UserStatus{
					UserMac: k,
				}
				user.DelUserStatusByMac()
				delete(listener, k)
			}
		}
		
		beego.Debug("Listener Gc running...")
		time.Sleep(GC_INTERVAL * time.Minute)
	}
}

func init() {
	listener = make(map[string]time.Time)

	go run()
}

func addListener(mac string) {
	listener[mac] = time.Now()
}

func delListener(mac string) {
	delete(listener, mac)
}