package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"radgo"
	"ums/v1/models"
)

type UpdateController struct {
	beego.Controller
}

func (this *UpdateController) Get() {
	this.TplNames = "home.html"
}

func (this *UpdateController) Post() {
	beego.Debug("requestBody=", string(this.Ctx.Input.RequestBody))
	
	code := &StatusCode{}
	info := &models.UserUpdate{}
	
	err := json.Unmarshal(this.Ctx.Input.RequestBody, info)
	if err != nil {
		code.Write(this.Ctx, -2)
		
		return
	}
	beego.Debug("update info=", info)

	user := &models.UserStatus{
		UserMac:  info.UserMac,
	}
	
	// liujf: don't check exist, just check one
	if !user.Exist() {
		beego.Info("UserStatus had been deleted when update come")
		code.Write(this.Ctx, -4)
		
		return
	}
	
	if nil != user.One() {
		code.Write(this.Ctx, -2)
		
		return
	}
	
	user.FlowDown 	= info.FlowDown
	user.FlowUp		= info.FlowUp
	
	//check with radius
	raduser := &models.RadUserstatus{
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
	if !user.Update() {
		code.Write(this.Ctx, -2)
		
		return
	}

	//插入listener
	addAlive(user.UserMac)

	//返回给设备处理结果
	code.Write(this.Ctx, 0)
}
