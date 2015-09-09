package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type UserInfo struct {
	UserName         string    `orm:"pk";json:"username"`
	LastRegisterTime time.Time `json:"-"`
}

func (account *UserInfo) TableName() string {
	return "userinfo"
}

func (this *UserInfo) RegisterUserinfo() bool {
	o := orm.NewOrm()
	this.LastRegisterTime = time.Now()
	beego.Debug("regiteraccount table=", this.TableName())
	//查找对应的username是否存在
	
	if ok := o.QueryTable(this.TableName()).
			Filter("username", this.UserName).
			Exist(); ok {
		//account存在，则更新account信息
		//return UpdateUserinfo(account)
		return true
	} else {
		//account不存在，则插入account信息
		if _, err := o.Insert(this); err != nil {
			beego.Error(err)
			return false
		}
		return true
	}
}

func (this *UserInfo) UpdateUserinfo() bool {
	beego.Debug("UpdateUserinfo table=", this.TableName())
	
	acc := &UserInfo{}
	o := orm.NewOrm()
	
	if  err := o.QueryTable(this.TableName()).
			Filter("username", this.UserName).
			One(acc);
		err != nil {
		return false
	} else {
		acc.UserName = this.UserName
		acc.LastRegisterTime = this.LastRegisterTime

		beego.Debug("UpdateUserinfo Username =", acc.UserName)

		if _, err := o.Update(acc); err != nil {
			beego.Error(err)
			return false
		}
		return true
	}
}

func (this *UserInfo) Init() {
	this.UserName = cutLastChar(this.UserName)
}