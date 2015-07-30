package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
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
	//查询statmap

	var ret RegisterData

	beego.Info("request body=", string(this.Ctx.Input.RequestBody))

	var account models.Userinfo

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &account)
	if err != nil {
		ret.Code = -4
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	} else {
		//注册account到数据库
		if !models.RegisterUserinfo(&account) {
			ret.Code = -5
			writeContent, _ := json.Marshal(ret)
			this.Ctx.WriteString(string(writeContent))
			return
		}
	}
	beego.Info("insert table useraccount success!")

	//request webserver??

	//查询添加statmap??

	//返回给设备处理结果
	ret.Code = 0
	writeContent, _ := json.Marshal(ret)
	this.Ctx.WriteString(string(writeContent))
	return
}
