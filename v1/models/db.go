package models

import (
	. "asdf"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type IDbEntry interface {
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

		dbUrl := dbUser + ":" + dbPassword + 
			"@tcp(" + dbIp + ":" + dbPort + ")/" + 
			dbName + "?charset=utf8&loc=Asia%2FShanghai"
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

func dbOrmer(o orm.Ormer) orm.Ormer {
	if nil==o {
		return ormer
	} else {
		return o
	}
}

func dbEntryGet(orm orm.Ormer, e IDbEntry, one interface{}) error {
	beego.Debug("get", e.TableName(), "by entry", e)
	err := dbOrmer(orm).QueryTable(e.TableName()).
			Filter(e.KeyName(), e.Key()).
			One(one)
	beego.Debug("get", e.TableName(), "new entry", one)
	
	return err
}

func dbEntryPull(orm orm.Ormer, e IDbEntry) error {	
	beego.Debug("before pull", e.TableName(), "entry", e)
	err := dbOrmer(orm).QueryTable(e.TableName()).
			Filter(e.KeyName(), e.Key()).
			One(e)
	beego.Debug("after pull", e.TableName(), "entry", e)
	
	return err
}

func dbEntryExist(orm orm.Ormer, e IDbEntry) bool {
	ok := dbOrmer(orm).QueryTable(e.TableName()).
			Filter(e.KeyName(), e.Key()).
			Exist()
	
	beego.Debug(e.TableName(), "entry", e, "exist", ok)
	
	return ok
}

func dbEntryInsert(orm orm.Ormer, e IDbEntry) error {
	count, err := dbOrmer(orm).Insert(e)
	
	beego.Debug("insert", e.TableName(), 
		"entry", e, 
		"count", count,
		"error", err)
	
	return dbError(count, err)
}

func dbEntryUpdate(orm orm.Ormer, e IDbEntry) error {
	count, err := dbOrmer(orm).Update(e)
	
	beego.Debug("update", e.TableName(), 
		"entry", e, 
		"count", count,
		"error", err)
	
	return dbError(count, err)
}

func dbEntryDelete(orm orm.Ormer, e IDbEntry) error {
	count, err := dbOrmer(orm).Delete(e)
	
	beego.Debug("delete", e.TableName(), 
		"entry", e, 
		"count", count,
		"error", err)
	
	return dbError(count, err)
}

func dbEntryRegister(orm orm.Ormer, e IDbEntry) error {
	if !dbEntryExist(orm, e) {
		return dbEntryInsert(orm, e)
	}
	
	return dbEntryUpdate(orm, e)
}
