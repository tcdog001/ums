package models

import (
	. "asdf"
	"radgo"
	"time"
)

type UserStatus struct {
	UserName     string
	UserIp       string
	UserMac      string    `orm:"pk"`
	DevMac       string
	Ssid         string
	AuthCode     string
	
	FlowUp       uint64
	FlowDown     uint64
	AuthTime     time.Time `orm:"type(datetime)"`
	
	// radgo.DeauthReason
	Reason 		 int
	
	// radius state, save in db
	RadSession   	[]byte
	RadClass 		[]byte
	RadChallenge	[]byte
	
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
	
	this.AuthCode = CutLastChar(this.AuthCode)
	this.Ssid = CutLastChar(this.Ssid)
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

func (this *UserStatus) Get() error {
	return dbEntryGet(nil, this)
}

func (this *UserStatus) Exist() bool {
	return dbEntryExist(nil, this)
}

func (this *UserStatus) Insert() error {
	return dbEntryInsert(nil, this)
}

func (this *UserStatus) Update() error {
	return dbEntryUpdate(nil, this)
}

func (this *UserStatus) Delete() error {
	return dbEntryDelete(nil, this)
}

func (this *UserStatus) Register() error {
	this.AuthTime = time.Now()
	
	return dbEntryRegister(nil, this)
}

func (this *UserStatus) RadUser() *RadUser {
	return &RadUser{
		User: this,
	}
}