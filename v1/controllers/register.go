package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"sms/sms_fx"
	"ums/v1/models"
	//"strings"
)

type RegisterController struct {
	beego.Controller
}

func (this *RegisterController) Get() {
	this.TplNames = "home.html"
}

var webserver = beego.AppConfig.String("WbServer")

func (this *RegisterController) Post() {
	//解析json
	//insert regtable
	//request webserver
	body := this.Ctx.Input.RequestBody
	
	beego.Info("request body=", string(body))

	code := &StatusCode{}
	info := &models.UserInfo{}

	if err := json.Unmarshal(body, info); err != nil {
		code.Write(this.Ctx, -2)
		
		return
	}
	info.Init()

	// liujf
	//	check user state from db
	//		is registered: is error, abort it
	//		not registered: go on
	
	//check with sms webserver
	res, err := sms_fx.SendCreateAccount(webserver, info.UserName, 10)
	if err != nil || (nil!=res && !res.Result) {
		beego.Debug("error:Check with sms server failed!")
		code.Write(this.Ctx, -3)
		
		return
	}
	//注册account到数据库
	if !info.Register() {
		code.Write(this.Ctx, -2)
		return
	}
	beego.Info("insert table useraccount success!")

	//返回给设备处理结果
	code.Write(this.Ctx, 0)
	
	return
}


