package controllers

import "bookms/models"

type SearchBookController struct {
	BaseController
}

func (c *SearchBookController) SearchBooks() {
	keyword := c.GetString(":keyword")
	books,cnt, err := new(models.Book).SearchBook(keyword,1,10)
	if err != nil {
		c.JsonResult(404, "找不到相关书籍",err.Error())
	}
	c.Data["SearchBooks"] = books
	c.Data["Count"] = cnt
	c.TplName="search/search.html"
}
