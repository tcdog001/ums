package models

import (
	"github.com/astaxie/beego"
	"time"
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
				user := &UserStatus{
					UserMac: k,
				}

				dbEntryDelete(nil, user)
				delete(alive, k)
			}
		}
		
		beego.Debug("Listener Gc running...")
		time.Sleep(GC_INTERVAL * time.Minute)
	}
}

func aliveInit() {
	alive = make(map[string]time.Time)

	go run()
}

func AddAlive(mac string) {
	alive[mac] = time.Now()
}

func DelAlive(mac string) {
	delete(alive, mac)
}