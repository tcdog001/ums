package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"sms/sms_fx"
	mod "ums/v1/models"
	//"strings"
)

type registerInput struct {
	UserName string    `json:"username"`
}

func (this *registerInput) Init() {
	this.UserName = mod.CutLastChar(this.UserName)
}

type RegisterController struct {
	beego.Controller
}

func (this *RegisterController) Get() {
	this.TplNames = "home.html"
}

var webserver = beego.AppConfig.String("WbServer")

func (this *RegisterController) Post() {
	body := this.Ctx.Input.RequestBody
	beego.Info("request body=", string(body))
	
	//step 1: get input
	code := &StatusCode{}
	input := &registerInput{}
	if err := json.Unmarshal(body, input); err != nil {
		code.Write(this.Ctx, -2)
		
		return
	}
	input.Init()
	
	//step 2: have registered ?
	info := &mod.UserInfo{
		UserName: input.UserName,
	}

	exist := false
	if nil == info.Get() { // exist
		exist = true
		
		if info.Registered { // have registered
			code.Write(this.Ctx, -3)
		
			return
		}
	}
	
	//step 3: register to sms webserver
	res, err := sms_fx.SendCreateAccount(webserver, input.UserName, 10)
	if err != nil || (nil!=res && !res.Result) {
		beego.Debug("error:Check with sms server failed!")
		code.Write(this.Ctx, -4)
		
		return
	}
	
	//step 4: register to db
	if nil != info.Register(exist) {
		code.Write(this.Ctx, -5)
		
		return
	}
	
	//step 5: output
	code.Write(this.Ctx, 0)
	
	return
}


