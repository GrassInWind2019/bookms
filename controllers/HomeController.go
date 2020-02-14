package controllers

import (
	"bookms/cache"
	"bookms/models"
	"bookms/utils"
	"github.com/astaxie/beego/logs"
	"strconv"
)

type HomeController struct {
	BaseController
	Muser models.User
}

func (c *HomeController) Prepare() {
	sess,err := GlobalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		logs.Error(err.Error())
		return
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
			logs.Debug("user: ", c.Muser.Id, " "+c.Muser.Nickname)
			if err = sess.Set("user", c.Muser); err != nil {
				logs.Error("set session failed:", err.Error())
			}
			return
		}
		logs.Error("Get user from cookie failed.")
	}
	c.Muser.Id = -1
	logs.Debug("User anonymous.")
}

func (c *HomeController) Index() {
	c.TplName = "index.tpl"
	var category models.Category
	var topCategories []models.Category
	err := cache.GetInterface("top_categories", &topCategories)
	if err != nil {
		logs.Debug(err.Error())
		topCategories, err = category.GetTopCategories()
		if err != nil {
			logs.Error("Index: GetTopCategories ", err.Error())
			c.JsonResult(500, "服务器内部错误，获取分类信息失败")
		}
		if err,_ = cache.SetInterface("top_categories", topCategories, 3600); err != nil {
			logs.Debug("set top_categories", err.Error())
		}
	}else {
		logs.Debug("Get top category from cache")
	}

	c.Data["TopCategories"] = topCategories
	var homeBooks []models.Book

	err = cache.GetInterface("home_books", &homeBooks)
	if err != nil {
		logs.Debug(err.Error())
		for _,tCate := range topCategories {
			cates, err := tCate.GetCategoriesByPid(tCate.Id)
			if err != nil {
				logs.Error("Index: GetCategoriesByPid ", strconv.Itoa(tCate.Id)," ", err.Error())
			}
			var book models.Book
			for _, cate := range cates {
				books,_,err := book.GetBooksByCategory2(cate.Id, 1, 5)
				if err != nil {
					logs.Error("Index: GetBooksByCategory2 ", strconv.Itoa(cate.Id), " ", err.Error())
					c.JsonResult(500, "获取图书信息失败")
				}
				logs.Debug("Index: ", books)
				homeBooks = append(homeBooks, books...)
			}
		}
		if err,_ = cache.SetInterface("home_books", homeBooks); err != nil {
			logs.Debug("Set home_books", err.Error())
		}
	}else {
		logs.Debug("Get home books from cache")
	}

	c.Data["HomeBooks"] = homeBooks
	if c.Muser.Id > 0 {
		c.Data["IsLogin"] = 1
	} else {
		c.Data["IsLogin"] = 0
	}
}
