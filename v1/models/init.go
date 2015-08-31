package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	. "asdf"
)

var localSwitch bool = true

func init() {
	//初始化log
	log.log = logs.NewLogger(10000)
	log.log.SetLogger("console", "") // use file
	SetLogger(&log)
	
	radParamInit()
	
	orm.RegisterModel(new(Userinfo), new(Userstatus),new(Userrecord))
	//register mysql driver
	err := orm.RegisterDriver("mysql", orm.DR_MySQL)
	if err != nil {
		beego.Critical(err)
	}
	//register default database
	if !localSwitch {
		orm.RegisterDataBase("default", "mysql", "autelan:Autelan1202@tcp(rdsrenv7vrenv7v.mysql.rds.aliyuncs.com:3306)/umsdb?charset=utf8&&loc=Asia%2FShanghai")
	} else {
		dbIp := beego.AppConfig.String("DbIp")
		dbPort := beego.AppConfig.String("DbPort")
		dbName := beego.AppConfig.String("DbName")
		dbUser := beego.AppConfig.String("DbUser")
		dbPassword := beego.AppConfig.String("DbPassword")

		dbUrl := dbUser + ":" + dbPassword + "@tcp(" + dbIp + ":" + dbPort + ")/" + dbName + "?charset=utf8&loc=Asia%2FShanghai"
		beego.Debug("dbUrl=", dbUrl)

		err = orm.RegisterDataBase("default", "mysql", dbUrl)
		if err != nil {
			beego.Critical(err)
		}
	}

	orm.SetMaxIdleConns("default", 30)
	orm.SetMaxOpenConns("default", 30)
}

func CheckDatabase() bool {
	o := orm.NewOrm()
	err := o.Using("default")
	if err != nil {
		return false
	} else {
		return true
	}
}
