package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"radgo"
	mod "ums/v1/models"
)

type deauthInput struct {
	UserMac  	string 	`json:"usermac"`
	Reason   	string 	`json:"reason"`
}

func (this *deauthInput) UserStatus() *mod.UserStatus {
	return &mod.UserStatus{
		UserMac: this.UserMac,
	}
}

type DeauthController struct {
	beego.Controller
}

func (this *DeauthController) Get() {
	this.TplNames = "home.html"
}

func (this *DeauthController) Post() {
	body := this.Ctx.Input.RequestBody
	beego.Info("request body=", string(body))
	
	//step 1: get input
	code := &StatusCode{}
	input := &deauthInput{}
	
	if err := json.Unmarshal(body, input); nil!=err {
		code.Write(this.Ctx, ErrUmsInputError, err)
		
		return
	}
	beego.Debug("deauth input", input)
	
	//step 2: get user from db
	user := input.UserStatus()	
	if nil != user.Get() {
		code.Write(this.Ctx, ErrUmsUserStatusNotExist, nil)
		
		return
	}
	
	user.Reason = int(radgo.GetDeauthReason(input.Reason))
	
	//step 3: radius acct stop
	raduser := user.RadUser()
	if err, aerr := radgo.ClientAcctStop(raduser); err != nil {
		code.Write(this.Ctx, ErrUmsRadAcctStopError, err)
		
		return
	} else if aerr != nil {
		code.Write(this.Ctx, ErrUmsRadError, aerr)
		
		return
	}
	beego.Debug("Redius stop success!")

	if nil != user.Delete() {
		code.Write(this.Ctx, ErrUmsUserStatusDeleteError, nil)
		
		return
	}
	
	//step 4: delete user(db)
	if err := user.Delete(); nil!=err {
		beego.Debug("delete user", user, )
		
		// NOT abort, must do below
	}
	
	//step 5: stop keepalive
	mod.DelAlive(user.UserMac)
	
	//step 6: log user record
	mod.LogUserRecord(user)
	
	//step 7: output
	code.Write(this.Ctx, 0, nil)
	
	return
}
