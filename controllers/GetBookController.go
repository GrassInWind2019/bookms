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

//type BookRecordStatus int
//
//const(
//	CannotLend BookRecordStatus = iota
//	CanLend
//	CanReturn
//)

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
	isScored := 0
	if c.Muser.Id > 0 {
		fav := models.Favorite{
			Identify:identify,
			UserId:c.Muser.Id,
		}
		isFav,err := fav.IsFavorite()
		if err != nil {
			//c.Data["IsFavorite"] = false
			c.Data["IsFavorite"] = 0
		} else {
			//c.Data["IsFavorite"] = isFav
			if isFav {
				c.Data["IsFavorite"] = 1
			} else {
				c.Data["IsFavorite"] = 0
			}
		}
		scoreObj := models.Score{
			UserId:c.Muser.Id,
			Identify:identify,
		}
		_, err = scoreObj.GetBookScore()
		if err != nil {
			isScored = 1
		}
	} else {
		c.Data["IsFavorite"] = 0
	}
	book_comments, err := new(models.Comments).GetBookCommentsAndScores(identify, 1, 5)
	if err != nil {
		//c.JsonResult(500, err.Error())
		logs.Debug("GetBookCommentsAndScores: ", err.Error())
	}
	c.Data["BookComments"] = book_comments
	c.Data["UserId"] = c.Muser.Id
	c.Data["IsScored"] = isScored
	c.Data["ScoreNum"] = book_comments[0].Score / 10
	logs.Debug("UserId: ", c.Muser.Id)
	logs.Debug(book_records)
}
