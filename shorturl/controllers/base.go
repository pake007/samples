package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

func (b *BaseController) Prepare() {
	logs.SetLogger(logs.AdapterConsole)
}
