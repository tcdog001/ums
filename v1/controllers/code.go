package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type StatusCode struct {
	Code int `json:"code"`
}

func (me *StatusCode) Write(ctx *context.Context, code EUmsError, err error) {
	me.Code = int(code)

	j, _ := json.Marshal(me)

	ctx.WriteString(string(j))

	beego.Info(string(j), code.ToString(), err)
}
