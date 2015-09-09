package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"radgo"
	"time"
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
	info := &models.Userupdate{}
	
	err := json.Unmarshal(this.Ctx.Input.RequestBody, info)
	if err != nil {
		code.Write(this.Ctx, -2)
		
		return
	}
	beego.Debug("update info=", info)

	user := &models.Userstatus{
		Usermac:  info.Usermac,
		Flowup:   info.Flowup,
		Flowdown: info.Flowdown,
	}
	
	exist := user.IsFindUserstatusByMac()
	if !exist {
		beego.Info("Userstatus had been deleted when update come")
		code.Write(this.Ctx, -4)
		
		return
	}
	//check with radius
	err1 := user.FindUserstatusByMac()
	if err1 != nil {
		code.Write(this.Ctx, -2)
		
		return
	}
	
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
	err3 := user.UpdateUserstatusBymac()
	if !err3 {
		code.Write(this.Ctx, -2)
		
		return
	}

	//插入listener
	Listener[user.Usermac] = time.Now()
	//返回给设备处理结果
	code.Write(this.Ctx, 0)
}
