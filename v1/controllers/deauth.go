package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"radgo"
	"ums/v1/models"
)

type DeauthRet struct {
	Code int64 `json:"code"`
}

type DeauthController struct {
	beego.Controller
}

func (this *DeauthController) Get() {
	this.TplNames = "home.html"
}
func (this *DeauthController) Post() {
	ret := DeauthRet{}

	beego.Info("request body=", string(this.Ctx.Input.RequestBody))

	locuser := models.Userstatus{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &locuser)
	if err != nil {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	user := models.Userstatus{
		Usermac: locuser.Usermac,
	}
	exist := models.FindUserstatusByMac(&user)
	if !exist {
		ret.Code = -4
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	//stop redius??
	radusr := models.RadUserstatus{
		User: &user,
	}
	_, err1 := radgo.ClientAcctStop(&radusr)
	if err1 != nil {
		ret.Code = -3
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	err2 := models.DelUserStatusByMac(&user)
	if !err2 {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	//del from listener
	delete(Listener, user.Usermac)

	//返回成功
	ret.Code = 0
	writeContent, _ := json.Marshal(ret)
	this.Ctx.WriteString(string(writeContent))
	beego.Info(string(writeContent))
	return
}
