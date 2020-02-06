package controllers

import (
	"bookms/cache"
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
	var page int
	page,err = c.GetInt(":page", 1)
	if err != nil {
		page = 1
	}
	fav := models.Favorite{
		UserId:c.Muser.Id,
	}
	books, cnt, err := fav.ListFavoriteByUserId(page,100)
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	is_admin := user.IsAdmin()
	c.Data["IsAdmin"] = is_admin

	c.Data["UserInfo"] = *user
	if 0 >= cnt {
		c.Data["MyFavoriteCount"] = 0
	} else {
		c.Data["MyFavoriteCount"] = cnt
		c.Data["MyFavorite"] = books
	}
	c.TplName = "user/userCenter.html"
}

func (c *UserController) GetUserCenterInfo2() {
	user := &models.User{
		Id:c.Muser.Id,
	}
	user,err := user.Find(c.Muser.Id)
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	var page int
	page,err = c.GetInt(":page", 1)
	if err != nil {
		page = 1
	}
	fav := models.Favorite{
		UserId:c.Muser.Id,
	}
	books, cnt, err := fav.ListFavoriteByUserId2(page,100)
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	is_admin := user.IsAdmin()
	c.Data["IsAdmin"] = is_admin

	c.Data["UserInfo"] = *user
	if 0 >= cnt {
		c.Data["MyFavoriteCount"] = 0
	} else {
		c.Data["MyFavoriteCount"] = cnt
		c.Data["MyFavorite"] = books
	}
	c.TplName = "user/userCenter.html"
}

func (c *UserController) GetUserCenterInfo3() {
	user := &models.User{
		Id:c.Muser.Id,
	}
	user,err := user.Find(c.Muser.Id)
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	var page int
	page,err = c.GetInt(":page", 1)
	if err != nil {
		page = 1
	}
	fav := models.Favorite{
		UserId:c.Muser.Id,
	}
	books, cnt, err := fav.ListFavoriteByUserId3(page,100)
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	is_admin := user.IsAdmin()
	c.Data["IsAdmin"] = is_admin

	c.Data["UserInfo"] = *user
	if 0 >= cnt {
		c.Data["MyFavoriteCount"] = 0
	} else {
		c.Data["MyFavoriteCount"] = cnt
		c.Data["MyFavorite"] = books
	}
	c.TplName = "user/userCenter.html"
}

func (c *UserController) GetUserCenterFav() {
	user := &models.User{
		Id:c.Muser.Id,
	}
	user,err := user.Find(c.Muser.Id)
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	var page int
	page,err = c.GetInt(":page", 1)
	if err != nil {
		page = 1
	}
	fav := models.Favorite{
		UserId:c.Muser.Id,
	}
	books, cnt, err := fav.ListFavoriteByUserIdReturnUserFav(page,100)
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	is_admin := user.IsAdmin()
	c.Data["IsAdmin"] = is_admin

	c.Data["UserInfo"] = *user
	if 0 >= cnt {
		c.Data["MyFavoriteCount"] = 0
	} else {
		c.Data["MyFavoriteCount"] = cnt
		c.Data["MyFavorite"] = books
	}
	c.TplName = "user/userCenter.html"
}

func (c *UserController) GetUserCenterFav2() {
	user := &models.User{
		Id:c.Muser.Id,
	}
	user,err := user.Find(c.Muser.Id)
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	var page int
	page,err = c.GetInt(":page", 1)
	if err != nil {
		page = 1
	}
	fav := models.Favorite{
		UserId:c.Muser.Id,
	}
	var user_favs []*models.UserFavorite
	var cnt int64
	err1 := cache.GetInterface("user_favs-"+strconv.Itoa(c.Muser.Id), &user_favs)
	err2, res := cache.GetInt("user_favs_cnt-"+strconv.Itoa(c.Muser.Id))
	if err1 == nil && err2 == nil {
		logs.Debug("Get user_favs-",strconv.Itoa(c.Muser.Id)," from cache")
		cnt = int64(res)
	} else {
		logs.Debug("Get user favs ",err1.Error(),err2.Error())
		user_favs, cnt, err = fav.ListFavoriteByUserIdReturnUserFav2(page,100)
		if err != nil {
			c.JsonResult(500, err.Error())
		}
		cache.SetInterface("user_favs-"+strconv.Itoa(c.Muser.Id), user_favs, 600)
		cache.SetInt("user_favs_cnt-"+strconv.Itoa(c.Muser.Id), int(cnt), 600)
	}

	var is_admin int
	err, is_admin = cache.GetInt("is_admin-"+strconv.Itoa(c.Muser.Id))
	if err != nil {
		logs.Debug("Get is_admin",strconv.Itoa(c.Muser.Id),err.Error())
		is_admin = user.IsAdmin()
		c.Data["IsAdmin"] = is_admin
		cache.SetInt("is_admin-"+strconv.Itoa(c.Muser.Id), is_admin)
	} else {
		c.Data["IsAdmin"] = is_admin
		logs.Debug("Get is_admin",strconv.Itoa(c.Muser.Id)," from cache")
	}

	c.Data["UserInfo"] = *user
	if 0 >= cnt {
		c.Data["MyFavoriteCount"] = 0
	} else {
		c.Data["MyFavoriteCount"] = cnt
		c.Data["MyFavorite"] = user_favs
	}
	c.TplName = "user/userCenter.html"
}

func (c *UserController) GetUserCenterFav3() {
	user := &models.User{
		Id:c.Muser.Id,
	}
	user,err := user.Find(c.Muser.Id)
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	var page int
	page,err = c.GetInt(":page", 1)
	if err != nil {
		page = 1
	}
	fav := models.Favorite{
		UserId:c.Muser.Id,
	}
	books, cnt, err := fav.ListFavoriteByUserIdReturnUserFav3(page,100)
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	is_admin := user.IsAdmin()
	c.Data["IsAdmin"] = is_admin

	c.Data["UserInfo"] = *user
	if 0 >= cnt {
		c.Data["MyFavoriteCount"] = 0
	} else {
		c.Data["MyFavoriteCount"] = cnt
		c.Data["MyFavorite"] = books
	}
	c.TplName = "user/userCenter.html"
}
