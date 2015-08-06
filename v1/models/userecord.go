package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Userecord struct {
	Username   string
	Usermac    string
	Devmac     string
	Authtime   time.Time
	Deauthtime time.Time
}

func (user *Userecord) TableName() string {
	return "userecord"
}
func RegisterUserecord(user *Userecord) bool {
	o := orm.NewOrm()
	_, err := o.Insert(user)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}
