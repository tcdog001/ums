package models

import (
	. "asdf"
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
//以下为实现radgo的IAuth IAcct IParam接口
//************************************************************
type RadUserstatus struct {
	User *UserStatus
}

// IAuth
func (user *RadUserstatus) UserPassword() []byte {
	return []byte(user.User.AuthCode)
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
	return user.User.RadSession
}

// IAcct
func (user *RadUserstatus) UserName() []byte {
	return []byte(user.User.UserName)
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
	return uint32(user.User.FlowUp & 0xffffffff)
}

// IAcct
func (user *RadUserstatus) AcctOutputOctets() uint32 {
	return uint32(user.User.FlowDown & 0xffffffff)
}

// IAcct
func (user *RadUserstatus) AcctInputGigawords() uint32 {
	return uint32(user.User.FlowUp >> 32)
}

// IAcct
func (user *RadUserstatus) AcctOutputGigawords() uint32 {
	return uint32(user.User.FlowDown >> 32)
}

// IAcct
func (user *RadUserstatus) AcctTerminateCause() uint32 {
	return radgo.DeauthReason(user.User.DeauthReason).TerminateCause()
}

// IAcct
func (user *RadUserstatus) GetClass() []byte {
	return user.User.RadClass
}

// IAcct
func (user *RadUserstatus) SetClass(c []byte) {
	user.User.RadClass = c
}

// IAcct
func (user *RadUserstatus) GetChapChallenge() []byte {
	return user.User.RadChallenge
}

// IAcct
func (user *RadUserstatus) SetChapChallenge(c []byte) {
	user.User.RadChallenge = c
}

// IParam
func (user *RadUserstatus) Secret() []byte {
	return []byte(param.RadSecret)
}

// IParam
func (user *RadUserstatus) NasIdentifier() []byte {
	//passwd for redius in configure
	return []byte(param.NasIdentifier)
}

// IParam
func (user *RadUserstatus) NasIpAddress() uint32 {
	return uint32(param.NasIpAddress)
}

// IParam
func (user *RadUserstatus) NasPort() uint32 {
	return param.NasPort
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
func (user *RadUserstatus) AuthType() uint32 {
	return param.AuthType
}

// IParam
func (user *RadUserstatus) Server() string {
	return param.RadServer
}

// IParam
func (user *RadUserstatus) AuthPort() string {
	return param.AuthPort
}

// IParam
func (user *RadUserstatus) AcctPort() string {
	return param.AcctPort
}

// IParam
func (user *RadUserstatus) DmPort() string {
	return param.DmPort
}

// IParam
func (user *RadUserstatus) Timeout() uint32 {
	return param.RadTimeout
}

type radParam struct {
	RadSecret		string
	NasIdentifier	string
	RadServer		string
	AuthPort		string
	AcctPort 		string
	DmPort			string
	
	NasIpAddress	IpAddress
	
	NasPort 		uint32
	RadTimeout		uint32
	AuthType 		uint32
}

var param = &radParam{}

func radParamString(name string) string {
	return beego.AppConfig.String(name)
}

func radParamUint32(name string) uint32 {
	i, _ := strconv.Atoi(radParamString(name))
	
	return uint32(i)
}

func radParamIpAddress(name string) IpAddress {
	return IpAddressFromString(radParamString(name))
}

func radParamInit() {
	param.RadSecret 	= radParamString("RadSecret")
	param.NasIdentifier = radParamString("NasIdentifier")
	param.RadServer 	= radParamString("RadServer")
	param.AuthPort 		= radParamString("AuthPort")
	param.AcctPort 		= radParamString("AcctPort")
	param.DmPort 		= radParamString("DmPort")
	
	param.NasIpAddress 	= radParamIpAddress("NasIpAddress")
	
	param.AuthType 		= radParamUint32("AuthType")
	param.NasPort 		= radParamUint32("NasPort")
	param.RadTimeout	= radParamUint32("RadTimeout")
}
