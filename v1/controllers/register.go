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
func (this *RegisterController) Post() {
	//解析json
	//insert regtable
	//request webserver
	
	beego.Info("request body=", string(this.Ctx.Input.RequestBody))

	code := &StatusCode{}
	account := &models.Userinfo{}

	err := json.Unmarshal(this.Ctx.Input.RequestBody, account)
	if err != nil {
		code.Write(this.Ctx, -2)
		
		return
	}
	account.Init()

	// liujf
	//	check user state from db
	//		is registered: is error, abort it
	//		not registered: go on
	
	//check with sms webserver
	webserver := beego.AppConfig.String("WbServer")
	res, err := sms_fx.SendCreateAccount(webserver, account.Username, 10)
	if err != nil || res.Result != true {
		beego.Debug("error:Check with sms server failed!")
		code.Write(this.Ctx, -3)
		
		return
	}
	//注册account到数据库
	if !account.RegisterUserinfo() {
		code.Write(this.Ctx, -2)
		return
	}
	beego.Info("insert table useraccount success!")

	//返回给设备处理结果
	code.Write(this.Ctx, 0)
	
	return
}


