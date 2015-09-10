package models

import (
	. "asdf"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type IDbOpt interface {
	TableName() string
	KeyName() string
	Key() string
}

var ormer orm.Ormer
var localSwitch = true

func dbInit() {
	orm.RegisterModel(new(UserInfo), new(UserStatus),new(UserRecord))
	
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
	
	ormer = orm.NewOrm()
}

func dbError(count int64, err error) error {
	if 0==count {
		beego.Error(ErrNoExist)
		
		return ErrNoExist
	} else if count > 1 {
		beego.Error(ErrTooMore)
		
		return ErrTooMore
	} else if nil != err {
		beego.Error(err)
		
		return err
	}
	
	return nil
}

func DbEntryGet(e IDbOpt, one interface{}) error {
	return ormer.QueryTable(e.TableName()).
			Filter(e.KeyName(), e.Key()).
			One(one)
}

func DbEntryPull(e IDbOpt) error {
	return ormer.QueryTable(e.TableName()).
			Filter(e.KeyName(), e.Key()).
			One(e)
}

func DbEntryExist(e IDbOpt) bool {
	return ormer.QueryTable(e.TableName()).
			Filter(e.KeyName(), e.Key()).
			Exist()
}

func DbEntryInsert(e IDbOpt) error {
	beego.Debug("insert", e.TableName(), "entry", e)
		
	count, err := ormer.Insert(e)
	
	return dbError(count, err)
}

func DbEntryUpdate(e IDbOpt) error {
	beego.Debug("update", e.TableName(), "entry", e)
		
	count, err := ormer.Update(e)
	
	return dbError(count, err)
}

func DbEntryDelete(e IDbOpt) error {
	beego.Debug("delete", e.TableName(), 
		e.KeyName(), "=", e.Key())
	
	count, err := ormer.Delete(e)
	
	return dbError(count, err)
}

func DbEntryRegister(e IDbOpt) error {
	if !DbEntryExist(e) {
		return DbEntryInsert(e)
	}
	
	return DbEntryUpdate(e)
}
