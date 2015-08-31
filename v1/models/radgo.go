package models

import (
	"asdf"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"radgo"
	"strconv"
)

//************************************************
//以下为实现radgo的logger的接口
//************************************************
type mylog struct {
	log *logs.BeeLogger
}

var log mylog

func (me *mylog) Emerg(format string, v ...interface{}) {
	me.log.Emergency(format, v...)
}

func (me *mylog) Alert(format string, v ...interface{}) {
	me.log.Alert(format, v...)
}

func (me *mylog) Crit(format string, v ...interface{}) {
	me.log.Critical(format, v...)
}

func (me *mylog) Error(format string, v ...interface{}) {
	me.log.Error(format, v...)
}

func (me *mylog) Warning(format string, v ...interface{}) {
	me.log.Warning(format, v...)
}
func (me *mylog) Notice(format string, v ...interface{}) {
	me.log.Notice(format, v...)
}

func (me *mylog) Info(format string, v ...interface{}) {
	me.log.Informational(format, v...)
}

func (me *mylog) Debug(format string, v ...interface{}) {
	me.log.Debug(format, v...)
}

//************************************************************
//以下为实现radgo的IAuth IAcct接口
//************************************************************
type RadUserstatus struct {
	User *Userstatus
}

// IAuth
func (user *RadUserstatus) UserPassword() []byte {

	return []byte(user.User.Authcode)
}

// IAcct
func (user *RadUserstatus) SSID() []byte {

	return []byte(user.User.Ssid)
}

// IAcct
func (user *RadUserstatus) DevMac() []byte {
	return user.User.devmac[:]
}

// IAcct
func (user *RadUserstatus) SessionId() []byte {
	// return user session
	// when new user, use ClientSessionId init session
	return []byte(user.User.radSession)
}

// IAcct
func (user *RadUserstatus) UserName() []byte {
	return []byte(user.User.Username)
}

// IAcct
func (user *RadUserstatus) UserMac() []byte {
	return user.User.usermac[:]
}

// IAcct
func (user *RadUserstatus) UserIp() uint32 {
	return uint32(user.User.userip)
}

// IAcct
func (user *RadUserstatus) AcctInputOctets() uint32 {
	var flow uint32
	flow = (uint32)(user.User.Flowup & 0xffffffff)
	return flow
}

// IAcct
func (user *RadUserstatus) AcctOutputOctets() uint32 {
	var flow uint32
	flow = (uint32)(user.User.Flowdown & 0xffffffff)
	return flow
}

// IAcct
func (user *RadUserstatus) AcctInputGigawords() uint32 {
	var flow uint32
	flow = (uint32)(user.User.Flowup >> 32)
	return flow
}

// IAcct
func (user *RadUserstatus) AcctOutputGigawords() uint32 {
	var flow uint32
	flow = uint32(user.User.Flowdown >> 32)
	return flow
}

// IAcct
func (user *RadUserstatus) AcctTerminateCause() uint32 {
	return radgo.DeauthReason(user.User.DeauthReason).TerminateCause()
}

// IAcct
func (user *RadUserstatus) GetClass() []byte {
	return user.User.radClass
}

// IAcct
func (user *RadUserstatus) SetClass(class []byte) {
	user.User.radClass = class
}

// IParam
func (user *RadUserstatus) Secret() []byte {
	secret := beego.AppConfig.String("RadSecret")
	fmt.Println("RadSecret:", secret)
	return []byte(secret)
}

// IParam
func (user *RadUserstatus) NasIdentifier() []byte {
	//passwd for redius in configure
	Identifier := beego.AppConfig.String("NasIdentifier")
	return []byte(Identifier)
}

// IParam
func (user *RadUserstatus) NasIpAddress() uint32 {
	ip := beego.AppConfig.String("NasIpAddress")
	uip := uint32(asdf.IpAddressFromString(ip))
	return uip
}

// IParam
func (user *RadUserstatus) NasPort() uint32 {
	port := beego.AppConfig.String("NasPort")
	uport, err := strconv.Atoi(port)
	if err == nil {
		return uint32(uport)
	} else {
		return 0
	}
}

// IParam
func (user *RadUserstatus) NasPortType() uint32 {
	return uint32(radgo.AnptIeee80211)
}

// IParam
func (user *RadUserstatus) ServiceType() uint32 {
	return uint32(radgo.AstLogin)
}

// IParam
func (user *RadUserstatus) Server() string {
	server := beego.AppConfig.String("RadServer")
	return server
}

// IParam
func (user *RadUserstatus) AuthPort() string {
	port := beego.AppConfig.String("AuthPort")
	return port
}

// IParam
func (user *RadUserstatus) AcctPort() string {
	port := beego.AppConfig.String("AcctPort")
	return port
}

// IParam
func (user *RadUserstatus) Timeout() int {
	time := beego.AppConfig.String("RadTimeout")
	utime, err := strconv.Atoi(time)
	if err == nil {
		return utime
	} else {
		return 0
	}
}
