package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"radgo"
	"ums/v1/models"
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

	code := &AuthCode{}
	user := &models.UserStatus{}
	
	if err := json.Unmarshal(body, user); err != nil {
		code.Write(this.Ctx, -2)
		
		return
	}
	user.Init()
	//var policy radgo.Policy

	// liujf:
	//	check user state from db
	// 		is registered: go on
	//		not registered: error, abort it
	
	//check with redius
	radusr := &models.RadUserstatus{
		User: user,
	}
	
	var policy *radgo.Policy
	
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
	}else if aerr != nil {
		beego.Info("ClientAcctStart:Radius failed!")
		code.Write(this.Ctx, -3)
		
		return
	}
	
	//注册user到数据库
	if !user.RegisterUserstatus() {
		code.Write(this.Ctx, -2)
		
		return
	}
	
	//插入listener
	addListener(user.UserMac)

	//返回给设备处理结果
	code.UpFlowLimit = policy.UpFlowLimit
	code.UpRateMax = policy.UpRateMax
	code.UpRateAvg = policy.UpRateAvg
	code.DownFlowLimit = policy.DownFlowLimit
	code.DownRateMax = policy.DownRateMax
	code.DownRateAvg = policy.DownRateAvg
	
	code.Write(this.Ctx, 0)

	return
}

