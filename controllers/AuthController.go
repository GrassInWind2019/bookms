package controllers

import (
	"bookms/models"
	"bookms/utils"
	"github.com/astaxie/beego"
)

type AuthController struct {
	BaseController
	Muser models.User
}

func (c *AuthController) Prepare() {
	if cookie, ok := c.GetSecureCookie(secretCookie, "user"); ok {
		if err := utils.Decode(cookie, &c.Muser); err == nil && c.Muser.Id > 0 {
			return
		}
	}
	c.Redirect(beego.URLFor("LoginController.Login"), 302)
	c.StopRun()
}
