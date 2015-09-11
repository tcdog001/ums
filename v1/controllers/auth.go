package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"radgo"
	mod "ums/v1/models"
)

type AuthCode struct {
	StatusCode
	
	IdleTimeout uint32 `josn:"idletimeout"`
	OnlineTime  uint32 `json:"onlinetime"`

	UpFlowLimit uint64 `json:"upflowlimit"`
	UpRateMax   uint32 `json:"upratemax"`
	UpRateAvg   uint32 `json:"uprateavg"`

	DownFlowLimit uint64 `json:"downflowlimit"`
	DownRateMax   uint32 `json:"downratemax"`
	DownRateAvg   uint32 `json:"downrateavg"`
}

func (me *AuthCode) WritePolicy(ctx *context.Context, policy *radgo.Policy) {
	me.UpFlowLimit 		= policy.UpFlowLimit
	me.UpRateMax 		= policy.UpRateMax
	me.UpRateAvg 		= policy.UpRateAvg
	me.DownFlowLimit 	= policy.DownFlowLimit
	me.DownRateMax 		= policy.DownRateMax
	me.DownRateAvg 		= policy.DownRateAvg
	
	me.Write(ctx, 0)
}

type UserAuthController struct {
	beego.Controller
}

func (this *UserAuthController) Get() {
	this.TplNames = "home.html"
}

func (this *UserAuthController) Post() {
	//解析json
	//查询redius(验证码+phoneno)
	//insert userinfotable
	//modify statmap
	body := this.Ctx.Input.RequestBody
	beego.Info("request body=", string(body))

	//step 1: get input
	code := &AuthCode{}
	user := &mod.UserStatus{}
	if err := json.Unmarshal(body, user); err != nil {
		code.Write(this.Ctx, -2)
		
		return
	}
	user.Init()
	
	//step 2: check registered
	info := &mod.UserInfo{
		UserName: user.UserName,
	}
	
	if !info.IsRegistered() {
		code.Write(this.Ctx, -3)
		
		return
	}
	
	//step 3: radius auth and acct start
	var policy *radgo.Policy
	
	radusr := &mod.RadUser{
		User: user,
	}
	
	if p, err, aerr := radgo.ClientAuth(radusr); err != nil {
		beego.Info("ClientAuth:username/password failed!")
		code.Write(this.Ctx, -3)
		
		return
	} else if aerr != nil {
		beego.Info("ClientAuth:Radius failed!")
		code.Write(this.Ctx, -1)
		
		return
	} else {
		policy = p
	}
	
	if err, aerr := radgo.ClientAcctStart(radusr); err != nil {
		beego.Info("ClientAcctStart:Failed when check with radius!")
		code.Write(this.Ctx, -3)
		
		return
	} else if aerr != nil {
		beego.Info("ClientAcctStart:Radius failed!")
		code.Write(this.Ctx, -3)
		
		return
	}
	
	//step 4: register user status
	if nil == user.Register() {
		code.Write(this.Ctx, -2)
		
		return
	}
	
	//step 5: keepalive
	mod.AddAlive(user.UserName, user.UserMac)

	//step 6: output
	code.WritePolicy(this.Ctx, policy)

	return
}

