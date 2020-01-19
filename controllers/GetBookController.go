package controllers

import "bookms/models"

type GetBookController struct {
	BaseController
}

func (c *GetBookController) GetBooksByIdentify() {
	c.TplName = "book/bookdetail.html"
	identify := c.Ctx.Input.Param(":identify")
	if "" == identify {
		c.Abort("404")
	}
	var identifies []string
	identifies = append(identifies, identify)
	books, err := new(models.Book).GetBooksByIdentifies(identifies)
	if err != nil {
		c.JsonResult(400, err.Error())
	}
	book_records, err := new(models.BookRecord).GetBookRecordsByIdentify(identify)
	if err != nil {
		c.JsonResult(400, err.Error())
	}
	c.Data["BookRecords"] = book_records
	c.Data["Books"] = books[0]
	//c.JsonResult(200, "OK")
}
