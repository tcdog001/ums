package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type StatusCode struct {
	Code int64 `json:"code"`
}

func (me *StatusCode) Write(ctx *context.Context, code int64) {
	me.Code = code
	
	wc, _ := json.Marshal(me)
	
	ctx.WriteString(string(wc))
	
	beego.Info(string(wc))
}