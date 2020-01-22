package controllers

import (
	"bookms/models"
	"bookms/utils"
	"github.com/astaxie/beego/logs"
)

type GetBookController struct {
	BaseController
	Muser models.User
}

func (c *GetBookController) Prepare() {
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
	if c.Muser.Id > 0 {
		fav := models.Favorite{
			Identify:identify,
			UserId:c.Muser.Id,
		}
		isFav,err := fav.IsFavorite()
		if err != nil {
			c.Data["IsFavorite"] = false
		} else {
			c.Data["IsFavorite"] = isFav
		}
	} else {
		c.Data["IsFavorite"] = false
	}
	book_comments, err := new(models.Comments).GetBookCommentsAndScores(identify, 1, 10)
	if err != nil {
		//c.JsonResult(500, err.Error())
		logs.Debug("GetBookCommentsAndScores: ", err.Error())
	}
	c.Data["BookComments"] = book_comments
}
