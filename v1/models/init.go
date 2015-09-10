package models

import (
	. "asdf"
	"github.com/astaxie/beego/logs"
)

func init() {
	logInit()
	radParamInit()	
	dbInit()
}

func logInit() {
	log.log = logs.NewLogger(10000)
	log.log.SetLogger("console", "") // use file
	SetLogger(&log)
}