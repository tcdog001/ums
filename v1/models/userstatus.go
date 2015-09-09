package models

import (
	. "asdf"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"radgo"
	"time"
)

type Userstatus struct {
	Username     string    `json:"username"`
	Userip       string    `json:"userip"`
	Usermac      string    `orm:"pk";json:"usermac"`
	Devmac       string    `json:"devmac"`
	Ssid         string    `json:"ssid"`
	Authcode     string    `json:"authcode"`
	Flowup       uint64    `json:"flowup"`
	Flowdown     uint64    `json:"flowdown"`
	AuthTime     time.Time `orm:"type(datetime)";json:"-"`
	DeauthReason int       `json:"-"`
	
	// radius state, save in db
	RadSession   	[]byte	`json:"-"`
	RadClass 		[]byte	`json:"-"`
	RadChallenge	[]byte	`json:"-"`
	
	// cache
	devmac       	[6]byte
	usermac      	[6]byte
	userip 		 	IpAddress
}

func (this *Userstatus) Init() {
	Mac(this.usermac[:]).FromString(this.Usermac)
	Mac(this.devmac[:]).FromString(this.Devmac)
	this.RadSession = radgo.NewSessionId(this.usermac[:], this.devmac[:])
	this.userip = IpAddressFromString(this.Userip)
	
	Len := len(this.Authcode)
	var b []byte = []byte(this.Authcode)
	if b[Len-1] == '_' {
		b = b[:Len-1]
		this.Authcode = string(b)
	}
	
	Len = len(this.Ssid)
	var c []byte = []byte(this.Ssid)
	if c[Len-1] == '_' {
		c = c[:Len-1]
		this.Ssid = string(c)
	}
}

func (user *Userstatus) TableName() string {
	return "userstatus"
}

func (this *Userstatus) RegisterUserstatus() bool {
	o := orm.NewOrm()
	this.AuthTime = time.Now()
	beego.Debug("userstatus table=", this.TableName())
	//查找对应的mac地址是否存在，存在的话要求状态为离线
	exist := o.QueryTable(this.TableName()).Filter("usermac", this.Usermac).Exist()
	if exist {
		//用户存在，则更新用户信息
		var u Userstatus
		err := o.QueryTable(this.TableName()).Filter("usermac", this.Usermac).One(&u)
		if err != nil {
			beego.Debug("get item in userstatus failed!!")
			return false
		} else {
			u.Username = this.Username
			u.Usermac = this.Usermac
			u.Authcode = this.Authcode
			u.Devmac = this.Devmac
			_, err := o.Update(&u)
			if err != nil {
				beego.Error(err)
				return false
			}
		}
	} else {
		//用户不存在，则插入用户信息
		_, err := o.Insert(this)
		if err != nil {
			beego.Error(err)
			return false
		}
	}
	return true
}

func (this *Userstatus) UpdateUserstatusBymac() bool {
	beego.Debug("Update userstatus table=", this.TableName())
	o := orm.NewOrm()

	var u Userstatus
	err := o.QueryTable(this.TableName()).Filter("usermac", this.Usermac).One(&u)
	if err != nil {
		return false
	} else {
		u.Flowup = this.Flowup
		u.Flowdown = this.Flowdown
		beego.Debug("Update userstatus usermac = ", u.Usermac)

		_, err := o.Update(&u)
		if err != nil {
			beego.Error(err)
			return false
		}
		return true
	}
}

func (this *Userstatus) IsFindUserstatusByMac() bool {
	o := orm.NewOrm()
	exist := o.QueryTable(this.TableName()).Filter("usermac", this.Usermac).Exist()
	return exist
}

func (this *Userstatus) FindUserstatusByMac() error {
	o := orm.NewOrm()
	err := o.QueryTable(this.TableName()).Filter("usermac", this.Usermac).One(this)
	return err
}

func (this *Userstatus) DelUserStatusByMac() bool {
	o := orm.NewOrm()

	var u Userstatus
	err := o.QueryTable(this.TableName()).Filter("usermac", this.Usermac).One(&u)
	if err != nil {
		beego.Error(err)
		return false
	}

	_, err = o.Delete(&u)
	if err != nil {
		beego.Error(err)
		return false
	}

	return true
}
