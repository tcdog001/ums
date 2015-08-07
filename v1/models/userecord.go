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

func (user *Userrecord) TableName() string {
	return "Userrecord"
}
func RegisterUserrecord(user *Userrecord) bool {
	o := orm.NewOrm()
	_, err := o.Insert(user)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}
