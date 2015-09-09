package models

import (
	. "asdf"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"radgo"
	"time"
)

type UserStatus struct {
	UserName     string    `json:"username"`
	UserIp       string    `json:"userip"`
	UserMac      string    `orm:"pk";json:"usermac"`
	DevMac       string    `json:"devmac"`
	Ssid         string    `json:"ssid"`
	AuthCode     string    `json:"authcode"`
	FlowUp       uint64    `json:"flowup"`
	FlowDown     uint64    `json:"flowdown"`
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

func (this *UserStatus) Init() {
	Mac(this.usermac[:]).FromString(this.UserMac)
	Mac(this.devmac[:]).FromString(this.DevMac)
	this.RadSession = radgo.NewSessionId(this.usermac[:], this.devmac[:])
	this.userip = IpAddressFromString(this.UserIp)
	
	this.AuthCode = cutLastChar(this.AuthCode)
	this.Ssid = cutLastChar(this.Ssid)
}

func (this *UserStatus) TableName() string {
	return "userstatus"
}

func (this *UserStatus) KeyName() string {
	return "usermac"
}

func (this *UserStatus) Key() string {
	return this.UserMac
}

func (this *UserStatus) Register() bool {
	o := orm.NewOrm()
	this.AuthTime = time.Now()
	beego.Debug("userstatus table=", this.TableName())
	//查找对应的mac地址是否存在，存在的话要求状态为离线
	
	if 	ok := EntryExist(o, this); ok {
		//用户存在，则更新用户信息
		var u UserStatus
		
		if 	err := EntryOne(o, this, &u);
			err != nil {
			beego.Debug("get item in userstatus failed!!")
			
			return false
		} else {
			u.UserName = this.UserName
			u.UserMac = this.UserMac
			u.AuthCode = this.AuthCode
			u.DevMac = this.DevMac
			
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

func (this *UserStatus) Update() bool {
	beego.Debug("Update userstatus table=", this.TableName())
	o := orm.NewOrm()

	var u UserStatus
	
	if 	err := EntryOne(o, this, &u); 
		err != nil {
		return false
	} else {
		u.FlowUp = this.FlowUp
		u.FlowDown = this.FlowDown
		beego.Debug("Update userstatus usermac = ", u.UserMac)

		if _, err := o.Update(&u); err != nil {
			beego.Error(err)
			return false
		}
		
		return true
	}
}

func (this *UserStatus) Exist() bool {
	return EntryExist(orm.NewOrm(), this)
}

func (this *UserStatus) One() error {
	return EntryOne(orm.NewOrm(), this, this)
}

func (this *UserStatus) Delete() bool {
	o := orm.NewOrm()

	var u UserStatus
	
	if 	err := EntryOne(o, this, &u); 
		err != nil {
		beego.Error(err)
		return false
	} else if _, err := o.Delete(&u); err != nil {
		beego.Error(err)
		return false
	}

	return true
}
