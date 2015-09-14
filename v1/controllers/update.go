package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"radgo"
	mod "ums/v1/models"
)

type updateInput struct {
	UserMac  string `json:"usermac"`
	FlowUp   uint64 `json:"flowup"`
	FlowDown uint64 `json:"flowdown"`
}

func (this *updateInput) UserStatus() *mod.UserStatus {
	return &mod.UserStatus{
		UserMac: this.UserMac,
	}
}

func (this *updateInput) UpdateUserStatus(user *mod.UserStatus) {
	user.FlowDown 	= this.FlowDown
	user.FlowUp		= this.FlowUp
}

type UpdateController struct {
	beego.Controller
}

func (this *UpdateController) Get() {
	this.TplNames = "home.html"
}

func (this *UpdateController) Post() {
	body := this.Ctx.Input.RequestBody
	beego.Debug("requestBody=", string(body))
	
	//step 1: get input
	code := &StatusCode{}
	input := &updateInput{}

	if err := json.Unmarshal(body, input); nil != err {
		code.Write(this.Ctx, ErrUmsInputError, err)
		
		return
	}
	beego.Debug("update input", input)

	//step 2: get and update user(local)
	user := input.UserStatus()
	if nil != user.Get() {
		code.Write(this.Ctx, ErrUmsUserStatusNotExist, nil)
		
		return
	}
	
	input.UpdateUserStatus(user)
	
	//step 3: radius acct update
	raduser := user.RadUser()
	
	if err, aerr := radgo.ClientAcctUpdate(raduser); nil != err {
		code.Write(this.Ctx, ErrUmsRadAcctUpdateError, err)
		
		return
	} else if nil != aerr {
		code.Write(this.Ctx, ErrUmsRadError, aerr)
		
		return
	}
	
	//step 4: update user(db)
	if err := user.Update(); nil!=err {
		beego.Debug("update", user, err)
		
		// NOT abort when update error
		// because not keepalive, wait timeout
	} else {
		//keepalive(just update ok)
		mod.AddAlive(user.UserName, user.UserMac)
	}

	//step 5: output
	code.Write(this.Ctx, 0, nil)
}
