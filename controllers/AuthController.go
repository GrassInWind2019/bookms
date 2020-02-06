package controllers

import (
	"bookms/models"
	"bookms/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/session"
	_ "github.com/astaxie/beego/session/redis"
)

type AuthController struct {
	BaseController
	Muser models.User
}

var GlobalSessions *session.Manager

func init() {
	config := beego.AppConfig.String("ProviderConfig")
	sessionConfig := &session.ManagerConfig{
		CookieName:"bookms-sessionid",
		EnableSetCookie: true,
		Gclifetime:1800,
		Maxlifetime: 1800,
		Secure: false,
		CookieLifeTime: 1800,
		ProviderConfig: config,
	}
	GlobalSessions, _ = session.NewManager("redis",sessionConfig)
	go GlobalSessions.GC()
}

func (c *AuthController) Prepare() {
	sess,err := GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		logs.Error(err.Error())
		c.Abort("500")
	}
	defer sess.SessionRelease(c.Ctx.ResponseWriter)
	user := sess.Get("user")
	if nil != user {
		c.Muser = user.(models.User)
		if c.Muser.Id > 0 {
			logs.Debug("Get user from session: ", c.Muser)
			return
		}
	}

	if cookie, ok := c.GetSecureCookie(secretCookie, "user"); ok {
		if err := utils.Decode(cookie, &c.Muser); err == nil && c.Muser.Id > 0 {
			if err = sess.Set("user", c.Muser); err != nil {
				logs.Error("set session failed:", err.Error())
			}
			return
		}
	}
	c.Redirect(beego.URLFor("LoginController.Login"), 302)
	c.StopRun()
}
