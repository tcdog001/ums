package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	//"radgo"
)

//************************************************
//以下为实现radgo的logger的接口
//************************************************
type mylog struct {
	log *logs.BeeLogger
}

var log mylog

func (me *mylog) Emerg(format string, v ...interface{}) {
	me.log.Emergency(format, v)
}

func (me *mylog) Alert(format string, v ...interface{}) {
	me.log.Alert(format, v)
}

func (me *mylog) Crit(format string, v ...interface{}) {
	me.log.Critical(format, v)
}

func (me *mylog) Error(format string, v ...interface{}) {
	me.log.Error(format, v)
}

func (me *mylog) Warning(format string, v ...interface{}) {
	me.log.Warning(format, v)
}
func (me *mylog) Notice(format string, v ...interface{}) {
	me.log.Notice(format, v)
}

func (me *mylog) Info(format string, v ...interface{}) {
	me.log.Informational(format, v)
}

func (me *mylog) Debug(format string, v ...interface{}) {
	me.log.Debug(format, v)
}

//************************************************************
//以下为实现radgo的IAuth IAcct接口
//************************************************************
type RadUserstatus struct {
	User *Userstatus
}

func (user *RadUserstatus) UserPassword() []byte {

	return nil
}

// IAcct
func (user *RadUserstatus) SessionId() []byte {
	sessionid := beego.AppConfig.String("RadSeesionId")
	return []byte(sessionid)
}

// IAcct
func (user *RadUserstatus) UserName() []byte {
	return []byte(user.User.Username)
}

// IAcct
func (user *RadUserstatus) UserMac() []byte {
	return []byte(user.User.Usermac)
}

// IAcct
func (user *RadUserstatus) UserIp() uint32 {
	return 0
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
	return 0 //reason
}

func (user *RadUserstatus) Secret() []byte {
	//passwd for redius in configure
	return nil
}
func (user *RadUserstatus) NasIdentifier() []byte {
	//passwd for redius in configure
	return nil
}
func (user *RadUserstatus) NasIpAddress() uint32 {
	//ums ip
	return 0
}
func (user *RadUserstatus) NasPort() uint32 {
	//ums port
	return 0
}
func (user *RadUserstatus) NasPortType() uint32 {
	return 0
}
func (user *RadUserstatus) NasPortId() uint32 {
	return 0
}
func (user *RadUserstatus) ServiceType() uint32 {
	return 0
}
func (user *RadUserstatus) Server() string {
	//redius 地址 from configure
	return ""
}
func (user *RadUserstatus) AuthPort() string {
	//reidus port for configure
	return ""
}
func (user *RadUserstatus) AcctPort() string {
	//reidus port for configure
	return ""
}
func (user *RadUserstatus) Timeout() int {
	return 0
}
