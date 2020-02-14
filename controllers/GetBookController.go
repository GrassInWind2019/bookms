package controllers

import (
	"bookms/cache"
	"bookms/models"
	"bookms/utils"
	"github.com/astaxie/beego/logs"
)

type GetBookController struct {
	BaseController
	Muser models.User
}

func (c *GetBookController) Prepare() {
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

func (c *GetBookController) GetBooksByIdentify() {
	c.TplName = "book/bookdetail.html"
	identify := c.Ctx.Input.Param(":identify")
	if "" == identify {
		c.Abort("404")
	}
	var books []models.Book
	err := cache.GetInterface("book-"+identify, &books)
	if err != nil {
		logs.Debug(identify, err.Error())
		var identifies []string
		identifies = append(identifies, identify)
		books, err = new(models.Book).GetBooksByIdentifies(identifies)
		if err != nil {
			logs.Error("GetBooksByIdentify: GetBooksByIdentifies ", err.Error())
			c.JsonResult(404, "找不到标识为"+identify+"的书")
		}
		err, _ = cache.SetInterface("book-"+identify, books, 3600)
		if err != nil {
			logs.Error("set ",identify,err.Error())
		}
	} else {
		logs.Debug("Get books ",identify, " from cache")
	}

	var book_records []models.BookRecord
	err = cache.GetInterface("book_record-"+identify, &book_records)
	if err != nil {
		logs.Debug("Get book record ", identify, err.Error())
		book_records, err = new(models.BookRecord).GetBookRecordsByIdentify(identify)
		if err != nil {
			logs.Error("GetBooksByIdentify: GetBookRecordsByIdentify ", err.Error())
			c.JsonResult(500, "获取图书标识为"+identify+"的借阅记录失败")
		}
		err,_ = cache.SetInterface("book_record-"+identify, book_records, 600)
		if err != nil {
			logs.Error("Set book record ", identify, err.Error())
		}
	}else {
		logs.Debug("Get book record", identify, " from cache")
	}

	c.Data["BookRecords"] = book_records
	c.Data["Books"] = books[0]
	isScored := 0
	if c.Muser.Id > 0 {
		var isFav int
		err, isFav = cache.GetInt("is_fav-"+identify)
		if err != nil {
			logs.Debug("Get is_fav-"+identify, err.Error())
			fav := models.Favorite{
				Identify:identify,
				UserId:c.Muser.Id,
			}
			isFav,err = fav.IsFavorite()
			if err != nil {
				logs.Error(err.Error())
				c.Data["IsFavorite"] = 0
				cache.SetInt("is_fav-"+identify, isFav,600)
			} else {
				cache.SetInt("is_fav-"+identify, isFav,600)
				c.Data["IsFavorite"] = isFav
			}
		} else {
			logs.Debug("Get is_fav-",identify," from cache")
			c.Data["IsFavorite"] = isFav
		}
		err, isScored = cache.GetInt("is_scored-"+identify)
		if err != nil {
			logs.Debug("Get is_scored-", identify, err.Error())
			scoreObj := models.Score{
				UserId:c.Muser.Id,
				Identify:identify,
			}
			_, err = scoreObj.GetBookScore()
			if err == nil {
				isScored = 1
				err, reply := cache.SetInt("is_scored-"+identify, isScored, 600)
				if err != nil {
					logs.Error(err.Error(), reply)
				}
			}
		} else {
			logs.Debug("Get is_scored-", identify, " from cache")
		}
	} else {
		c.Data["IsFavorite"] = 0
	}
	var book_comments []models.BookComment
	err = cache.GetInterface("book_comments-"+identify, &book_comments)
	if err != nil {
		logs.Debug("Get book_comments-",identify, err.Error())
		book_comments, err = new(models.Comments).GetBookCommentsAndScores(identify, 1, 5)
		if err != nil {
			//c.JsonResult(500, err.Error())
			logs.Error("GetBookCommentsAndScores: ", err.Error())
		}else {
			cache.SetInterface("book_comments-"+identify,book_comments, 600)
		}
	} else {
		logs.Debug("Get book_comments-",identify," from cache")
	}

	c.Data["BookComments"] = book_comments
	c.Data["UserId"] = c.Muser.Id
	c.Data["IsScored"] = isScored
	if len(book_comments) > 0 {
		c.Data["ScoreNum"] = book_comments[0].Score / 10
	}
	logs.Debug("UserId: ", c.Muser.Id)
	logs.Debug(book_records)
}
