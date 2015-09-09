package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Userrecord struct {
	Id 		   uint64 `orm:"auto"`
	Username   string
	Usermac    string
	Devmac     string
	Authtime   time.Time
	Deauthtime time.Time
}

func (this *Userrecord) TableName() string {
	return "Userrecord"
}

func (this *Userrecord) RegisterUserrecord() bool {
	o := orm.NewOrm()
	_, err := o.Insert(this)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}
