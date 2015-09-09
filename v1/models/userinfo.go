package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Userinfo struct {
	Username         string    `orm:"pk";json:"username"`
	LastRegisterTime time.Time `json:"-"`
}

func (account *Userinfo) TableName() string {
	return "userinfo"
}

func (this *Userinfo) RegisterUserinfo() bool {
	o := orm.NewOrm()
	this.LastRegisterTime = time.Now()
	beego.Debug("regiteraccount table=", this.TableName())
	//查找对应的username是否存在
	exist := o.QueryTable(this.TableName()).Filter("username", this.Username).Exist()
	if exist {
		//account存在，则更新account信息
		//return UpdateUserinfo(account)
		return true
	} else {
		//account不存在，则插入account信息
		_, err := o.Insert(this)
		if err != nil {
			beego.Error(err)
			return false
		}
		return true
	}
}

func (this *Userinfo) UpdateUserinfo() bool {
	beego.Debug("UpdateUserinfo table=", this.TableName())
	o := orm.NewOrm()

	var acc Userinfo
	err := o.QueryTable(this.TableName()).Filter("username", this.Username).One(&acc)
	if err != nil {
		return false
	} else {
		acc.Username = this.Username
		acc.LastRegisterTime = this.LastRegisterTime

		beego.Debug("UpdateUserinfo Username =", acc.Username)

		_, err := o.Update(&acc)
		if err != nil {
			beego.Error(err)
			return false
		}
		return true
	}
}

func (this *Userinfo) Init() {
	Len := len(this.Username)
	b := []byte(this.Username)
	
	if b[Len-1] == '_' {
		b = b[:Len-1]
		this.Username = string(b)
	}
}