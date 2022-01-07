package controllers

import (
	"shorturl/models"
)

type ExpandController struct {
	BaseController
}

func (e *ExpandController) Get() {
	var result models.ShortResult
	shortUrl := e.Ctx.Input.Query("shorturl")
	result.UrlShort = shortUrl

	if models.CacheCond.IsExist(shortUrl) {
		longUrl := models.CacheCond.Get(shortUrl)
		result.UrlLong = longUrl.(string)
	} else {
		result.UrlLong = ""
	}
	e.Data["json"] = result
	e.ServeJSON()
}
