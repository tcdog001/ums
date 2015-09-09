package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type UserRecord struct {
	Id 		   uint64 `orm:"auto"`
	UserName   string
	UserMac    string
	DevMac     string
	AuthTime   time.Time
	DeauthTime time.Time
}

func (this *UserRecord) TableName() string {
	return "Userrecord"
}

func (this *UserRecord) RegisterUserrecord() bool {
	o := orm.NewOrm()
	
	if _, err := o.Insert(this); err != nil {
		beego.Error(err)
		return false
	}
	
	return true
}
