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
	
	this.Authcode = cutLastChar(this.Authcode)
	this.Ssid = cutLastChar(this.Ssid)
}

func (user *Userstatus) TableName() string {
	return "userstatus"
}

func (this *Userstatus) RegisterUserstatus() bool {
	o := orm.NewOrm()
	this.AuthTime = time.Now()
	beego.Debug("userstatus table=", this.TableName())
	//查找对应的mac地址是否存在，存在的话要求状态为离线
	
	if 	ok := o.QueryTable(this.TableName()).
			Filter("usermac", this.Usermac).
			Exist(); ok {
		//用户存在，则更新用户信息
		var u Userstatus
		
		if 	err := o.QueryTable(this.TableName()).
				Filter("usermac", this.Usermac).
				One(&u); 
			err != nil {
			beego.Debug("get item in userstatus failed!!")
			
			return false
		} else {
			u.Username = this.Username
			u.Usermac = this.Usermac
			u.Authcode = this.Authcode
			u.Devmac = this.Devmac
			
			if _, err := o.Update(&u); err != nil {
				beego.Error(err)
				return false
			}
		}
	} else {
		//用户不存在，则插入用户信息
		if _, err := o.Insert(this); err != nil {
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
	
	if 	err := o.QueryTable(this.TableName()).
			Filter("usermac", this.Usermac).
			One(&u); 
		err != nil {
		return false
	} else {
		u.Flowup = this.Flowup
		u.Flowdown = this.Flowdown
		beego.Debug("Update userstatus usermac = ", u.Usermac)

		if _, err := o.Update(&u); err != nil {
			beego.Error(err)
			return false
		}
		
		return true
	}
}

func (this *Userstatus) IsFindUserstatusByMac() bool {
	o := orm.NewOrm()

	return o.QueryTable(this.TableName()).
			Filter("usermac", this.Usermac).
			Exist()
}

func (this *Userstatus) FindUserstatusByMac() error {
	o := orm.NewOrm()

	return o.QueryTable(this.TableName()).
			Filter("usermac", this.Usermac).
			One(this)
}

func (this *Userstatus) DelUserStatusByMac() bool {
	o := orm.NewOrm()

	var u Userstatus
	
	if 	err := o.QueryTable(this.TableName()).
			Filter("usermac", this.Usermac).
			One(&u); 
		err != nil {
		beego.Error(err)
		return false
	} else if _, err := o.Delete(&u); err != nil {
		beego.Error(err)
		return false
	}

	return true
}
