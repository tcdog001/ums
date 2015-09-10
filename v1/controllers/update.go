package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"radgo"
	mod "ums/v1/models"
)

type UpdateController struct {
	beego.Controller
}

func (this *UpdateController) Get() {
	this.TplNames = "home.html"
}

type devUserUpdate struct {
	UserMac  string `json:"usermac"`
	FlowUp   uint64 `json:"flowup"`
	FlowDown uint64 `json:"flowdown"`
}

func (this *UpdateController) Post() {
	body := this.Ctx.Input.RequestBody
	beego.Debug("requestBody=", string(body))
	
	code := &StatusCode{}
	info := &devUserUpdate{}
	
	err := json.Unmarshal(body, info)
	if err != nil {
		code.Write(this.Ctx, -2)
		
		return
	}
	beego.Debug("update info=", info)

	user := &mod.UserStatus{
		UserMac:  info.UserMac,
	}
	
	if nil != mod.DbEntryPull(user) {
		code.Write(this.Ctx, -2)
		
		return
	}
	
	user.FlowDown 	= info.FlowDown
	user.FlowUp		= info.FlowUp
	
	//check with radius
	raduser := &mod.RadUserstatus{
		User: user,
	}
	
	err2, res2 := radgo.ClientAcctUpdate(raduser)
	if err2 != nil {
		beego.Debug("error:Failed when check with radius!")
		code.Write(this.Ctx, -3)
		
		return
	}else if res2 != nil {
		beego.Debug("error:Radius failed!")
		code.Write(this.Ctx, -3)
		
		return
	}
	
	//update db
	if nil != mod.DbEntryUpdate(user) {
		code.Write(this.Ctx, -2)
		
		return
	}

	//插入listener
	mod.AddAlive(user.UserMac)

	//返回给设备处理结果
	code.Write(this.Ctx, 0)
}
