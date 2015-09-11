package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"radgo"
	mod "ums/v1/models"
)

type updateInput struct {
	UserMac  string `json:"usermac"`
	FlowUp   uint64 `json:"flowup"`
	FlowDown uint64 `json:"flowdown"`
}

type UpdateController struct {
	beego.Controller
}

func (this *UpdateController) Get() {
	this.TplNames = "home.html"
}

func (this *UpdateController) Post() {
	body := this.Ctx.Input.RequestBody
	beego.Debug("requestBody=", string(body))
	
	code := &StatusCode{}
	input := &updateInput{}
	
	err := json.Unmarshal(body, input)
	if err != nil {
		code.Write(this.Ctx, -2)
		
		return
	}
	beego.Debug("update info=", input)

	user := &mod.UserStatus{
		UserMac: input.UserMac,
	}
	
	if nil != user.Get() {
		code.Write(this.Ctx, -2)
		
		return
	}
	
	user.FlowDown 	= input.FlowDown
	user.FlowUp		= input.FlowUp
	
	//check with radius
	raduser := &mod.RadUser{
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
	if nil != user.Update() {
		code.Write(this.Ctx, -2)
		
		return
	}

	//插入listener
	mod.AddAlive(user.UserName, user.UserMac)

	//返回给设备处理结果
	code.Write(this.Ctx, 0)
}
