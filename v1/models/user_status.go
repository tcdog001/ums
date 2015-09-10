package models

import (
	. "asdf"
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

func (this *UserStatus) Register() error {
	this.AuthTime = time.Now()
	
	//查找对应的mac地址是否存在，存在的话要求状态为离线
	//用户不存在，则插入用户信息
	
	return DbEntryRegister(this)
}

