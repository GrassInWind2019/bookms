package controllers

import (
	"bookms/models"
	"bookms/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strings"
	"time"
)

const (
	secretCookie = "bookms"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Login() {
	logs.Debug("Enter login!")
	c.TplName = "login/login.html"
	u := &models.User{}
	if cookie, ok := c.GetSecureCookie(secretCookie, "user"); ok {
		if err := utils.Decode(cookie, *u); err == nil {
			if err = c.login(u.Id); err == nil {
				c.Redirect(beego.URLFor("HomeController.Index"), 302)
				c.StopRun()
			}
		}
	}
	if c.Ctx.Input.IsPost() {
		logs.Debug("Post login")
		account := c.GetString("账号")
		password := c.GetString("密码")
		if account == "" || password == "" {
			c.JsonResult(401, "账号或密码为空,登录失败")
		}
		u, err := u.Login(account, password)
		if err != nil {
			c.JsonResult(403, "账号或密码错误,登录失败")
		}
		if err := c.login(u.Id); err != nil {
			logs.Error("Login: login ", err.Error())
			c.JsonResult(500, "Internal error")
		}
		c.Success("登录成功",beego.URLFor("HomeController.Index"), 3)
		//c.Redirect(beego.URLFor("HomeController.Index"), 301)
	}
}

func (c *LoginController) Logout() {
	logs.Info("logout")
	c.SetSecureCookie(secretCookie, "user", "", 0)
	sess,err := GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		logs.Error("Logout: ",err.Error())
		c.Abort("500")
	}
	defer sess.SessionRelease(c.Ctx.ResponseWriter)
	sess.Delete("user")
	c.Success("登出成功",beego.URLFor("HomeController.Index"), 3)
}

func (c *LoginController) Register() {
	c.TplName = "register/register.html"
}

func (c *LoginController) RegisterDo() {
	account := c.GetString("Account")
	nickName := strings.TrimSpace(c.GetString("Nickname"))
	password1 := c.GetString("Password1")
	password2 := c.GetString("Password2")
	email := c.GetString("Email")
	phone := c.GetString("Phone")
	biography := c.GetString("Biography")

	if password1 != password2 {
		c.JsonResult(400, "两次输入密码不一致，请重新输入！")
	}
	if strings.Count(nickName, "") > 20 {
		c.JsonResult(400, "昵称长度不能超过20")
	}
	if strings.Count(biography, "") > 150 {
		c.JsonResult(400, "个人简介不能超过150个字")
	}

	u := models.User{
		Account:account,
		Nickname:nickName,
		Password:password1,
		Email:email,
		Phone:phone,
		Biography:biography,
		Role:2,
		CreateTime:time.Now(),
		LastLoginTime:time.Now(),
		Status: 0,
	}

	if err := u.Add(); err != nil {
		logs.Error("RegisterDo: ",err.Error())
		c.JsonResult(500, err.Error())
	}
	c.login(u.Id)
	c.Success("注册成功", beego.URLFor("HomeController.Index"), 3)
}

func (c *LoginController) login(id int) (error) {
	u := &models.User{}
	u, err := u.Find(id)
	if err != nil {
		return err
	}
	u.LastLoginTime = time.Now()
	u.Update()
	v, err := utils.Encode(*u)
	if err != nil {
		return err
	}
	c.SetSecureCookie(secretCookie, "user", v, 600000)
	return err
}