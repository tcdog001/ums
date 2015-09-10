package models

import (
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

func (this *UserInfo) Register() error {	
	this.LastRegisterTime = time.Now()
	
	return DbEntryRegister(this)
}

func (this *UserInfo) Init() {
	this.UserName = cutLastChar(this.UserName)
}