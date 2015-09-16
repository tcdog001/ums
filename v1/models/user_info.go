package models

import (
	"time"
)

type UserInfo struct {
	UserName string `orm:"pk"`
	//Registered       bool
	LastRegisterTime time.Time
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

func (this *UserInfo) Init() {
	this.UserName = CutLastChar(this.UserName)
}

func (this *UserInfo) Get() error {
	return dbEntryGet(nil, this)
}

func (this *UserInfo) Exist() bool {
	return dbEntryExist(nil, this)
}

func (this *UserInfo) Insert() error {
	return dbEntryInsert(nil, this)
}

func (this *UserInfo) Update() error {
	return dbEntryUpdate(nil, this)
}

func (this *UserInfo) Delete() error {
	return dbEntryDelete(nil, this)
}

func (this *UserInfo) IsRegistered() bool {
	return nil == this.Get()
}

func (this *UserInfo) Register() error {
	//this.Registered = true
	this.LastRegisterTime = time.Now()
	return dbEntryRegister(nil, this)
}

func (this *UserInfo) UnRegister() error {
	if err := this.Get(); nil != err {
		return err
	}

	//this.Registered = false

	return this.Update()
}
