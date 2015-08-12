package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"sms/sms_fx"
	"ums/v1/models"
	//"strings"
)

type RegisterData struct {
	Code int64 `json:"code"`
}

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

	var ret RegisterData

	beego.Info("request body=", string(this.Ctx.Input.RequestBody))

	var account models.Userinfo

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &account)
	if err != nil {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	account.Init_obj()

	//check with sms webserver
	webserver := beego.AppConfig.String("WbServer")
	res, err := sms_fx.SendCreateAccount(webserver, account.Username, 10)
	if err != nil || res.Result != true {
		beego.Debug("error:Check with sms server failed!")
		ret.Code = -3
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	//注册account到数据库
	if !models.RegisterUserinfo(&account) {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	beego.Info("insert table useraccount success!")

	//返回给设备处理结果
	ret.Code = 0
	writeContent, _ := json.Marshal(ret)
	this.Ctx.WriteString(string(writeContent))
	return
}


