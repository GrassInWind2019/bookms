package controllers

import (
	"bookms/models"
	"bookms/utils"
	"github.com/astaxie/beego"
)

type AuthController struct {
	BaseController
}

func (c *AuthController) Prepare() {
	u := &models.User{}
	if cookie, ok := c.GetSecureCookie(secretCookie, "user"); ok {
		if err := utils.Decode(cookie, *u); err == nil {
			return
		}
	}
	c.Redirect(beego.URLFor("LoginController.Login"), 302)
	c.StopRun()
}
