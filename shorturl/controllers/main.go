package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (m *MainController) Get() {
	m.Ctx.Output.Body([]byte("shorturl"))
}
