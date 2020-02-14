package controllers

import (
	"bookms/models"
	"github.com/astaxie/beego/logs"
)

type SearchBookController struct {
	BaseController
}

func (c *SearchBookController) SearchBooks() {
	keyword := c.GetString(":keyword")
	books,cnt, err := new(models.Book).SearchBook(keyword,1,10)
	if err != nil {
		logs.Error("找不到keyword为",keyword,"的相关书籍", err.Error())
	}
	c.Data["SearchBooks"] = books
	c.Data["Count"] = cnt
	c.Data["KeyWord"] = keyword
	c.TplName="search/search.html"
}
