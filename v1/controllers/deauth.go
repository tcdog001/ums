package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"radgo"
	mod "ums/v1/models"
)

type deauthInput struct {
	UserMac  	string 	`json:"usermac"`
	Reason   	int 	`json:"reason"`
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
		code.Write(this.Ctx, -2)
		
		return
	}
	beego.Debug("deauth input", input)
	
	//step 2: get user from db
	user := &mod.UserStatus{
		UserMac: 	input.UserMac,
	}
	
	if nil != user.Get() {
		code.Write(this.Ctx, -2)
		
		return
	}
	user.Reason = input.Reason
	
	//step 3: radius acct stop
	raduser := &mod.RadUser{
		User: user,
	}
	
	if err, aerr := radgo.ClientAcctStop(raduser); err != nil {
		beego.Debug("error:Failed when check with radius!")
		code.Write(this.Ctx, -3)
		
		return
	} else if aerr != nil {
		beego.Debug("error:Radius failed!")
		code.Write(this.Ctx, -3)
		
		return
	}
	beego.Debug("Redius stop success!")

	if nil != user.Delete() {
		code.Write(this.Ctx, -2)
		
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
	code.Write(this.Ctx, 0)
	
	return
}
