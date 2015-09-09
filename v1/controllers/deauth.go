package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"radgo"
	"ums/v1/models"
	"time"
)

type DeauthController struct {
	beego.Controller
}

func (this *DeauthController) Get() {
	this.TplNames = "home.html"
}

func (this *DeauthController) Post() {
	code := &StatusCode{}

	body := this.Ctx.Input.RequestBody
	beego.Info("request body=", string(body))

	luser := &models.UserStatus{}
	if err := json.Unmarshal(body, luser); nil!=err {
		code.Write(this.Ctx, -2)
		
		return
	}

	user := &models.UserStatus{
		UserMac: luser.UserMac,
	}
	
	if exist := user.IsFindUserstatusByMac(); !exist {
		code.Write(this.Ctx, -4)
		
		return
	}

	//check with redius
	if err := user.FindUserstatusByMac(); err != nil {
		code.Write(this.Ctx, -2)
		
		return
	}
	
	raduser := &models.RadUserstatus{
		User: user,
	}
	
	if err, aerr := radgo.ClientAcctStop(raduser); err != nil {
		beego.Debug("error:Failed when check with radius!")
		code.Write(this.Ctx, -3)
		
		return
	} else if aerr != nil {
		beego.Debug("error:Radius failed!")
		code.Write(this.Ctx, -3)
		
		return
	}
	beego.Debug("Redius stop success!")

	if ok := user.DelUserStatusByMac(); !ok {
		code.Write(this.Ctx, -2)
		
		return
	}
	
	//del from listener
	delListener(user.UserMac)
	
	//生成用户记录
	record := &models.UserRecord {
		UserName : user.UserName,
		UserMac : user.UserMac,
		DevMac : user.DevMac,
		AuthTime : user.AuthTime,
		DeauthTime : time.Now(),
	}
	record.RegisterUserrecord()
	
	//返回成功
	code.Write(this.Ctx, 0)
	
	return
}
