package models

import (
	"github.com/astaxie/beego"
	"time"
)

const (
	GC_INTERVAL      = 1  //(Minute)每隔一分钟清理一次超时user
	TIMEOUT_INTERVAL = 10 //(Minute)超时时间为10分钟
)

type aliveCache struct {
	HitTime 	time.Time
	UserName	string
}

var alive map[string]*aliveCache

func run() {
	for {
		for k, v := range alive {
			beego.Debug("key=", k, "v=", v)
			
			if time.Now().Sub(v.HitTime) >= time.Duration(TIMEOUT_INTERVAL)*time.Minute {
				//step 1: unregister user
				info := &UserInfo{
					UserName: v.UserName,
				}
				info.UnRegister()
				
				//step 2: stop redius
				
				//step 3: delete user
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
	alive = make(map[string]*aliveCache)

	go run()
}

func getAlive(mac string) *aliveCache {
	cache := alive[mac]
	if nil == cache {
		cache = &aliveCache{}
		alive[mac] = cache
	}
	
	return cache
}

func AddAlive(name string, mac string) {
	cache := getAlive(mac)
	
	cache.HitTime = time.Now()
	cache.UserName = name
}

func DelAlive(mac string) {
	delete(alive, mac)
}