package models

import (
	"github.com/astaxie/beego"
	"time"
)

type UserInfo struct {
	UserName         string    `orm:"pk";json:"username"`
	LastRegisterTime time.Time `json:"-"`
}

func (this *UserInfo) TableName() string {
	return "userinfo"
}

func (this *UserInfo) KeyName() string {
	return "username"
}

func (this *UserInfo) Key() string {
	return this.UserName
}

func (this *UserInfo) Register() bool {
	beego.Debug("regiter UserInfo table=", this.TableName())
	
	this.LastRegisterTime = time.Now()
	//查找对应的username是否存在
	
	if ok := DbEntryExist(this); ok {
		//account存在，则更新account信息
		//return UpdateUserinfo(account)
	} else {
		//account不存在，则插入account信息
		if _, err := ormer.Insert(this); err != nil {
			beego.Error(err)
			return false
		}
	}
	
	return true
}

func (this *UserInfo) Update() bool {
	beego.Debug("Update UserInfo table=", this.TableName())
	
	acc := &UserInfo{}
	
	if  err := DbEntryGet(this, acc); err != nil {
		return false
	} else {
		acc.UserName = this.UserName
		acc.LastRegisterTime = this.LastRegisterTime

		beego.Debug("Update UserInfo UserName =", acc.UserName)

		if _, err := ormer.Update(acc); err != nil {
			beego.Error(err)
			return false
		}
		return true
	}
}

func (this *UserInfo) Init() {
	this.UserName = cutLastChar(this.UserName)
}