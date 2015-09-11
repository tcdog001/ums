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
	//解析json
	//insert regtable
	//request webserver
	body := this.Ctx.Input.RequestBody
	beego.Info("request body=", string(body))

	code := &StatusCode{}
	input := &registerInput{}
	if err := json.Unmarshal(body, input); err != nil {
		code.Write(this.Ctx, -2)
		
		return
	}
	
	input.Init()
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
	
	//check with sms webserver
	res, err := sms_fx.SendCreateAccount(webserver, input.UserName, 10)
	if err != nil || (nil!=res && !res.Result) {
		beego.Debug("error:Check with sms server failed!")
		code.Write(this.Ctx, -4)
		
		return
	}
	
	//注册account到数据库
	if nil != info.Register(exist) {
		code.Write(this.Ctx, -5)
		
		return
		
	}
	
	//返回给设备处理结果
	code.Write(this.Ctx, 0)
	
	return
}


