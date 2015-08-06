package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"radgo"
	"time"
	"ums/v1/models"
)

type UpdateRetData struct {
	Code int64 `json:"code"`
}
type UpdateController struct {
	beego.Controller
}

func (this *UpdateController) Get() {
	this.TplNames = "home.html"
}
func (this *UpdateController) Post() {
	ret := UpdateRetData{}

	beego.Debug("requestBody=", string(this.Ctx.Input.RequestBody))
	upinfo := models.Userupdate{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &upinfo)
	if err != nil {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	beego.Debug("updateinfo=", upinfo)

	user := models.Userstatus{
		Usermac:  upinfo.Usermac,
		Flowup:   upinfo.Flowup,
		Flowdown: upinfo.Flowdown,
	}
	exist := models.IsFindUserstatusByMac(&user)
	if !exist {
		beego.Info("userstatus had been deleted when update come")
		ret.Code = -4
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	//update redius??
	//if failed return ??
	radusr := models.RadUserstatus{
		User: &user,
	}
	_, err2 := radgo.ClientAcctUpdate(&radusr)
	if err2 != nil {
		ret.Code = -3
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	err1 := models.UpdateUserstatusBymac(&user)
	if !err1 {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	//插入listener
	usrls := UserListener{
		LastAliveTime: time.Now(),
	}
	Listener[user.Usermac] = usrls

	//返回给设备处理结果
	ret.Code = 0
	writeContent, _ := json.Marshal(ret)
	this.Ctx.WriteString(string(writeContent))
}
