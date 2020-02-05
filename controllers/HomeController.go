package controllers

import (
	"bookms/models"
	"bookms/utils"
	"github.com/astaxie/beego/logs"
)

type HomeController struct {
	BaseController
	Muser models.User
}

func (c *HomeController) Prepare() {
	c.Muser.Id = -1
	if cookie, ok := c.GetSecureCookie(secretCookie, "user"); ok {
		if err := utils.Decode(cookie, &c.Muser); err == nil && c.Muser.Id > 0 {
			logs.Debug("user: ", c.Muser.Id, " "+c.Muser.Nickname)
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
	topCategories, err := category.GetTopCategories()
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	c.Data["TopCategories"] = topCategories
	var homeBooks []models.Book

	for _,tCate := range topCategories {
		cates, err := tCate.GetCategoriesByPid(tCate.Id)
		if err != nil {
			//c.JsonResult(500, err.Error())
		}
		var book models.Book
		for _, cate := range cates {
			books,_,err := book.GetBooksByCategory2(cate.Id, 1, 5)
			if err != nil {
				c.JsonResult(500, err.Error())
			}
			logs.Debug("Index: ", books)
			homeBooks = append(homeBooks, books...)
		}
	}
	c.Data["HomeBooks"] = homeBooks
	if c.Muser.Id > 0 {
		c.Data["IsLogin"] = 1
	} else {
		c.Data["IsLogin"] = 0
	}
}
