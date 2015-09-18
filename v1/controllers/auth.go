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

	me.IdleTimeout = policy.IdleTimeout
	me.OnlineTime = policy.OnlineTime
	me.UpFlowLimit = policy.UpFlowLimit
	me.UpRateMax = policy.UpRateMax
	me.UpRateAvg = policy.UpRateAvg
	me.DownFlowLimit = policy.DownFlowLimit
	me.DownRateMax = policy.DownRateMax
	me.DownRateAvg = policy.DownRateAvg

	me.Write(ctx, 0, nil)
}

func (me *AuthStatusCode) Write(ctx *context.Context, code EUmsError, err error) {

	me.Code = int(code)

	j, _ := json.Marshal(me)

	ctx.WriteString(string(j))

	beego.Info(string(j), code.ToString(), err)
}

type authInput struct {
	UserName string `json:"username"`
	UserIp   string `json:"userip"`
	UserMac  string `orm:"pk";json:"usermac"`
	DevMac   string `json:"devmac"`
	Ssid     string `json:"ssid"`
	AuthCode string `json:"authcode"`
}

func (this *authInput) Init() {
	this.UserName = mod.CutLastChar(this.UserName)
	this.AuthCode = mod.CutLastChar(this.AuthCode)
	this.Ssid = mod.CutLastChar(this.Ssid)
}

func (this *authInput) UserInfo() *mod.UserInfo {
	return &mod.UserInfo{
		UserName: this.UserName,
	}
}

func (this *authInput) UserStatus() *mod.UserStatus {
	user := &mod.UserStatus{
		UserName: this.UserName,
		UserIp:   this.UserIp,
		UserMac:  this.UserMac,
		DevMac:   this.DevMac,
		Ssid:     this.Ssid,
		AuthCode: this.AuthCode,
	}
	user.Init()

	return user
}

type UserAuthController struct {
	beego.Controller
}

func (this *UserAuthController) Get() {
	this.TplNames = "home.html"
}

func (this *UserAuthController) Post() {
	body := this.Ctx.Input.RequestBody
	beego.Info("request body=", string(body))

	//step 1: get input
	code := &AuthStatusCode{}
	input := &authInput{}

	if err := json.Unmarshal(body, input); err != nil {
		code.Write(this.Ctx, ErrUmsInputError, err)

		return
	}
	input.Init()
	beego.Debug("auth input", input)

	//step 2: check registered
	info := input.UserInfo()
	if !info.IsRegistered() {
		code.Write(this.Ctx, ErrUmsUserInfoNotRegistered, nil)

		return
	}

	//step 3: have authed ?
	user := input.UserStatus()
	if user.Exist() {
		code.Write(this.Ctx, ErrUmsUserHaveAuthed, nil)

		return
	}

	//step 4: radius auth and acct start
	raduser := user.RadUser()

	policy, err, aerr := radgo.ClientAuth(raduser)
	if nil != err {
		code.Write(this.Ctx, ErrUmsRadAuthError, err)

		return
	} else if nil != aerr {
		code.Write(this.Ctx, ErrUmsRadError, aerr)

		return
	}
	beego.Debug("auth received policy:[", policy, "]")

	err, aerr = radgo.ClientAcctStart(raduser)
	if nil != err {
		code.Write(this.Ctx, ErrUmsRadAcctStartError, err)

		return
	} else if nil != aerr {
		code.Write(this.Ctx, ErrUmsRadError, aerr)

		return
	}

	//step 5: register user status
	beego.Debug("user.UserName=", user.UserName)
	if err := user.Register(); nil != err {
		beego.Info("auth", user, err)

		//radius acct stop when register error
		//user.Reason = int(radgo.DeauthReasonNasError)
		//radgo.ClientAcctStop(raduser)

		code.Write(this.Ctx, ErrUmsUserStatusRegisterError, err)
		return
	}

	//step 6: keepalive(when register ok/fail)
	mod.AddAlive(user.UserName, user.UserMac)

	//step 7: output
	code.WritePolicy(this.Ctx, policy)

	return
}
