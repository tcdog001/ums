package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"radgo"
	mod "ums/v1/models"
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

	juser := &mod.UserStatus{}
	if err := json.Unmarshal(body, juser); nil!=err {
		code.Write(this.Ctx, -2)
		
		return
	}

	user := &mod.UserStatus{
		UserMac: juser.UserMac,
	}
	
	if nil != user.Get() {
		code.Write(this.Ctx, -2)
		
		return
	}
		
	//check with redius
	raduser := &mod.RadUser{
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

	if nil != user.Delete() {
		code.Write(this.Ctx, -2)
		
		return
	}
	//del from listener
	mod.DelAlive(user.UserMac)
	
	//生成用户记录
	record := &mod.UserRecord {
		UserName : user.UserName,
		UserMac : user.UserMac,
		DevMac : user.DevMac,
		AuthTime : user.AuthTime,
		DeauthTime : time.Now(),
	}
	record.Register()
	
	//返回成功
	code.Write(this.Ctx, 0)
	
	return
}
