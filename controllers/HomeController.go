package controllers

import "bookms/models"

type HomeController struct {
	BaseController
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
			books,_,err := book.GetBooksByCategory(cate.Id, 1, 5)
			if err != nil {
				c.JsonResult(500, err.Error())
			}
			homeBooks = append(homeBooks, books...)
		}
	}
	c.Data["HomeBooks"] = homeBooks
}
