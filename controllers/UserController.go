package controllers

import (
	"bookms/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
)

type UserController struct {
	AuthController
}

func (c *UserController) FavoriteDo() {
	identify := c.GetString(":identify")
	fav := models.Favorite{
		UserId:c.Muser.Id,
		Identify:identify,
	}
	if err := fav.FavoriteDo(); err != nil {
		c.JsonResult(500, err.Error())
	}
	//c.JsonResult(200, "收藏成功")
	logs.Debug(beego.URLFor("GetBookController.GetBooksByIdentify"))
	//c.Redirect(beego.URLFor("GetBookController.GetBooksByIdentify"), 302)
	url := "/bookdetail/"+identify
	logs.Debug("AddCommentAndScore: redirect to "+url)
	c.Redirect(url, 302)
}

func (c *UserController) AddCommentAndScore() {
	var score int64
	var err error
	scoreStr := c.GetString("starLevel")
	if scoreStr != "" {
		score, err = strconv.ParseInt(scoreStr, 10, 64)
		if err != nil {
			c.JsonResult(400, err.Error())
		}
		if score > 5 || score <= 0 {
			c.JsonResult(400, "请在0~5之间打分")
		}
	}

	comment_content := c.GetString("userComment")
	if "" == comment_content {
		c.JsonResult(400, "评论内容不能为空")
	}
	identify := c.Ctx.Input.Param(":identify")
	if "" == identify {
		c.JsonResult(400, "identify不能为空")
	}
	if scoreStr != "" {
		scoreObj := models.Score{
			Identify:identify,
			UserId:c.Muser.Id,
			Score:int(score),
			CreateTime:time.Now(),
		}
		if err := scoreObj.AddScore(); err != nil {
			c.JsonResult(500, err.Error())
		}
	}

	comment := models.Comments{
		Identify:identify,
		UserId:c.Muser.Id,
		Content:comment_content,
		CreateTime:time.Now(),
	}
	if err := comment.AddComment(); err != nil {
		c.JsonResult(500, err.Error())
	}
	//c.JsonResult(200, "发表成功")
	logs.Debug(beego.URLFor("GetBookController.GetBooksByIdentify"))
	//c.Redirect(beego.URLFor("GetBookController.GetBooksByIdentify"), 302)
	url := "/bookdetail/"+identify
	logs.Debug("AddCommentAndScore: redirect to "+url)
	c.Redirect(url, 302)
}

func (c *UserController) GetUserCenterInfo() {
	user := &models.User{
		Id:c.Muser.Id,
	}
	user,err := user.Find(c.Muser.Id)
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	fav := models.Favorite{
		UserId:c.Muser.Id,
	}
	books, cnt, err := fav.ListFavoriteByUserId()
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	is_admin := user.IsAdmin()
	if is_admin {
		c.Data["IsAdmin"] = 1
	} else {
		c.Data["IsAdmin"] = 0
	}

	c.Data["UserInfo"] = *user
	if 0 >= cnt {
		c.Data["MyFavoriteCount"] = 0
	} else {
		c.Data["MyFavoriteCount"] = cnt
		c.Data["MyFavorite"] = books
	}
	c.TplName = "user/userCenter.html"
}