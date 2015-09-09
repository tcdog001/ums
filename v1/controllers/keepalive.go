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

var alive map[string]time.Time

func run() {
	for {
		for k, v := range alive {
			beego.Debug("key=", k, "v=", v)
			
			if time.Now().Sub(v) >= time.Duration(TIMEOUT_INTERVAL)*time.Minute {
				//stop redius??
				user := &models.UserStatus{
					UserMac: k,
				}
				user.Delete()
				delete(alive, k)
			}
		}
		
		beego.Debug("Listener Gc running...")
		time.Sleep(GC_INTERVAL * time.Minute)
	}
}

func init() {
	alive = make(map[string]time.Time)

	go run()
}

func addAlive(mac string) {
	alive[mac] = time.Now()
}

func delAlive(mac string) {
	delete(alive, mac)
}