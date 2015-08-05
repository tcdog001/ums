package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"radgo"
	"time"
	"ums/v1/models"
)

type AuthRetData struct {
	Code        int64  `json:"code"`
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

	ret := AuthRetData{}
	beego.Info("request body=", string(this.Ctx.Input.RequestBody))

	var user models.Userstatus
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &user)
	if err != nil {
		ret.Code = -2
		setRetZero(&ret)
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	(&user).Init()

	//var policy radgo.Policy

	//check with redius
	radusr := models.RadUserstatus{
		User: &user,
	}
	policy, err := radgo.ClientAuth(&radusr)
	if err != nil {
		ret.Code = -3
		setRetZero(&ret)
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	_, err1 := radgo.ClientAcctStart(&radusr)
	if err1 != nil {
		ret.Code = -3
		setRetZero(&ret)
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	//注册user到数据库
	if !models.RegisterUserstatus(&user) {
		ret.Code = -2
		setRetZero(&ret)
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

	ret.UpFlowLimit = policy.UpFlowLimit
	ret.UpRateMax = policy.UpRateMax
	ret.UpRateAvg = policy.UpRateAvg
	ret.DownFlowLimit = policy.DownFlowLimit
	ret.DownRateMax = policy.DownRateMax
	ret.DownRateAvg = policy.DownRateAvg

	writeContent, _ := json.Marshal(ret)
	this.Ctx.WriteString(string(writeContent))

	return
}

func setRetZero(user *AuthRetData) bool {
	user.IdleTimeout = 0
	user.UpFlowLimit = 0
	user.UpRateMax = 0
	user.UpRateAvg = 0
	user.DownFlowLimit = 0
	user.DownRateMax = 0
	user.DownRateAvg = 0
	return true
}
