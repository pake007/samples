package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"shorturl/models"
)

type ShortController struct {
	BaseController
}

// Get Use Get rather than Post so that we can simulate easier in the browser
func (s *ShortController) Get() {
	var result models.ShortResult
	longUrl := s.Ctx.Input.Query("longurl")
	result.UrlLong = longUrl
	urlMd5 := models.GetMD5(longUrl)

	if models.CacheCond.IsExist(urlMd5) {
		shortUrl := models.CacheCond.Get(urlMd5)
		result.UrlShort = shortUrl.(string)
	} else {
		result.UrlShort = models.Generate()
		err := models.CacheCond.Put(urlMd5, result.UrlShort, 0)
		if err != nil {
			logs.Info(err)
		}
		err = models.CacheCond.Put(result.UrlShort, longUrl, 0)
		if err != nil {
			logs.Info(err)
		}
	}
	s.Data["json"] = result
	s.ServeJSON()
}
