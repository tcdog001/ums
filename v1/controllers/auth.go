package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"radgo"
	mod "ums/v1/models"
)

type AuthStatusCode struct {
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

func (me *AuthStatusCode) WritePolicy(ctx *context.Context, policy *radgo.Policy) {
	me.UpFlowLimit 		= policy.UpFlowLimit
	me.UpRateMax 		= policy.UpRateMax
	me.UpRateAvg 		= policy.UpRateAvg
	me.DownFlowLimit 	= policy.DownFlowLimit
	me.DownRateMax 		= policy.DownRateMax
	me.DownRateAvg 		= policy.DownRateAvg
	
	me.Write(ctx, 0)
}

type authInput struct {
	UserName     string    `json:"username"`
	UserIp       string    `json:"userip"`
	UserMac      string    `orm:"pk";json:"usermac"`
	DevMac       string    `json:"devmac"`
	Ssid         string    `json:"ssid"`
	AuthCode     string    `json:"authcode"`
}

func (this *authInput) Init() {
	this.UserName = mod.CutLastChar(this.UserName)
}

func (this *authInput) UserInfo() *mod.UserInfo {
	return &mod.UserInfo{
		UserName: this.UserName,
	}
}

func (this *authInput) UserStatus() *mod.UserStatus {
	return &mod.UserStatus{
		UserName:	this.UserName,
		UserIp:		this.UserIp,
		UserMac:	this.UserMac,
		DevMac:		this.DevMac,
		Ssid:		this.Ssid,
		AuthCode:	this.AuthCode,
	}
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
	code := &AuthStatusCode{}
	input := &authInput{}
	
	if err := json.Unmarshal(body, input); err != nil {
		code.Write(this.Ctx, ErrUmsInputError)
		
		return
	}
	input.Init()
	beego.Debug("auth input", input)
	
	//step 2: check registered
	info := input.UserInfo()	
	if !info.IsRegistered() {
		code.Write(this.Ctx, ErrUmsUserInfoNotRegistered)
		
		return
	}
	
	//step 3: radius auth and acct start
	user := input.UserStatus()
	raduser := user.RadUser()
	
	policy, err, aerr := radgo.ClientAuth(raduser)
	if nil != err {
		beego.Info("ClientAuth:username/password failed!")
		code.Write(this.Ctx, ErrUmsRadAuthError)
		
		return
	} else if nil != aerr {
		beego.Info("ClientAuth:Radius failed!")
		code.Write(this.Ctx, ErrUmsRadError)
		
		return
	}
	
	err, aerr = radgo.ClientAcctStart(raduser)
	if nil != err {
		beego.Info("ClientAcctStart:Failed when check with radius!")
		code.Write(this.Ctx, ErrUmsRadAcctStartError)
		
		return
	} else if nil != aerr {
		beego.Info("ClientAcctStart:Radius failed!")
		code.Write(this.Ctx, ErrUmsRadError)
		
		return
	}
	
	//step 4: register user status
	if err := user.Register(); nil!=err {
		beego.Info("auth", user, err)
		
		//radius acct stop when register error
		user.Reason = int(radgo.DeauthReasonNasError)
		radgo.ClientAcctStop(raduser)
		
		code.Write(this.Ctx, ErrUmsUserStatusRegisterError)
		return
	}
	
	//step 5: keepalive(when register ok/fail)
	mod.AddAlive(user.UserName, user.UserMac)

	//step 6: output
	code.WritePolicy(this.Ctx, policy)

	return
}

