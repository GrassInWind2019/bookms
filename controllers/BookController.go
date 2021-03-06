package controllers

import (
	"bookms/cache"
	bookms_init "bookms/init"
	"bookms/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
)

type BookController struct {
	AuthController
}

func (c *BookController) AddBook() {
	c.TplName = "book/addbook.html"
	if c.Ctx.Input.IsPost(){
		book_name := c.GetString("book_name")
		if book_name == "" {
			c.JsonResult(400, "书名不能为空")
		}
		description := c.GetString("description")
		catalog := c.GetString("catalog")
		sort, err := c.GetInt("sort")
		if err != nil {
			c.JsonResult(400, "排序应为数字")
		}
		author := c.GetString("author")
		if author == "" {
			c.JsonResult(400, "作者不能为空")
		}
		book_count, err := c.GetInt("book_count")
		if err != nil || book_count <= 0 {
			c.JsonResult(400, "图书数量应为数字")
		}
		category_id, err := c.GetInt("category_id")
		if err != nil {
			c.JsonResult(400, "分类id应为数字")
		}
		store_position := c.GetString("store_position")
		if "" == store_position {
			c.JsonResult(400, "存放位置不能为空")
		}
		//get an uuid for book by sonyflake
		sonyflakeObj := bookms_init.GetSonyFlakeObj()
		uuid, err := sonyflakeObj.NextID()
		if err != nil {
			c.JsonResult(500, "UUID create failed! "+err.Error())
		}
		identify := fmt.Sprint(uuid)
		if "" == identify {
			c.JsonResult(500, "UUID create failed!")
		}

		book := models.Book{
			BookName:book_name,
			Identify:identify,
			Description:description,
			Catalog:catalog,
			Sort:sort,
			Author:author,
			BookCount:book_count,
		}
		err = book.AddBook()
		if err != nil {
			c.JsonResult(500, err.Error())
		}
		var book_category models.BookCategory
		var cids []int
		cids = append(cids, category_id)
		if err = book_category.SetBookCategories(identify, cids); err != nil {
			c.JsonResult(500, err.Error())
		}
		var positions []string
		for i:=0;i<book_count;i++ {
			position := store_position
			positions = append(positions, position)
		}
		if err = new(models.BookRecord).AddBookRecords(book_count, identify, positions); err != nil {
			c.JsonResult(500, err.Error())
		}
		//c.JsonResult(200, "添加成功")
		c.Success("添加成功", beego.URLFor("UserController.GetUserCenterInfo"),1)
	}
}

func (c *BookController) UpdateBookByIdentify() {
	c.TplName="book/updatebook.html"
	logs.Debug("Enter UpdateBookByIdentify")
	identify := c.Ctx.Input.Param(":identify")
	if identify == "" {
		c.Abort("404")
	}
	if c.Ctx.Input.IsGet() {
		var books []models.Book
		var identifies []string
		identifies = append(identifies, identify)
		books, err := new(models.Book).GetBooksByIdentifies(identifies)
		if err != nil {
			logs.Error("UpdateBookByIdentify: GetBooksByIdentifies ", err.Error())
			c.JsonResult(404, "找不到标识为"+identify+"的书")
		}
		cids, err := new(models.BookCategory).GetBookCategories(identify)
		if err != nil {
			logs.Error("UpdateBookByIdentify: GetBookCategories", err.Error())
			c.JsonResult(500, err.Error())
		}
		if len(cids) <= 0 || cids[0] <= 0 {
			c.JsonResult(500, "get category id failed")
		}
		c.Data["BookName"] = books[0].BookName
		c.Data["Description"] = books[0].Description
		c.Data["Catalog"] = books[0].Catalog
		c.Data["Sort"] = books[0].Sort
		c.Data["Author"] = books[0].Author
		c.Data["BookCount"] = books[0].BookCount
		c.Data["CategoryId"] = cids
	}
	if c.Ctx.Input.IsPost() {
		book_name := c.GetString("book_name")
		description := c.GetString("description")
		catalog := c.GetString("catalog")
		sort,err := c.GetInt("sort")
		if err != nil {
			c.JsonResult(400, err.Error())
		}
		author := c.GetString("author")
		book_count, err := c.GetInt("book_count")
		if err != nil {
			c.JsonResult(400, err.Error())
		}
		category_id, err := c.GetInt("category_id")
		if err != nil {
			c.JsonResult(400, err.Error())
		}
		book := &models.Book{
			BookName:book_name,
			Identify:identify,
			Description:description,
			Catalog:catalog,
			Sort:sort,
			Author:author,
			BookCount:book_count,
		}

		err = book.UpdateBookByIdentify("book_name", "description","catalog","sort","author", "book_count")
		if err != nil {
			logs.Error("UpdateBookByIdentify: UpdateBookByIdentify", err.Error())
			c.JsonResult(500, "服务器内部错误，更新失败")
		}
		err, reply:=cache.SetExpire("book-"+identify, 0)
		if err != nil {
			logs.Error("Set book-"+identify+" expire:",err.Error(), reply)
		}
		logs.Debug("Set cache book-"+identify+" expired")
		book_category := models.BookCategory{
			Identify:identify,
			CategoryId:category_id,
		}
		cIds := []int{}
		cIds = append(cIds, category_id)
		err = book_category.SetBookCategories(identify, cIds)
		if err != nil {
			logs.Error("UpdateBookByIdentify: SetBookCategories "+identify+ " ", err.Error())
			c.JsonResult(500, "服务器内部错误，设置图书分类失败")
		}
		c.JsonResult(200, "更新成功")
	}
}

func (c *BookController) DeleteBooksByIdentify() {
	logs.Debug("Enter DeleteBooksByIdentify")
	identify := c.Ctx.Input.Param(":identify")
	if identify=="" {
		c.Abort("404")
	}
	ids, err := new(models.BookRecord).GetBookIdsByIdentify(identify)
	if err != nil {
		c.JsonResult(404, "标识为"+identify+"的书不存在")
	}
	if len(ids) > 0 {
		if err = new(models.BookRecord).DeleteBookRecordsByIds(ids); err != nil {
			logs.Error("DeleteBooksByIdentify: DeleteBookRecordsByIds ", err.Error())
			c.JsonResult(500, "服务器内部错误，删除失败")
		}
		err, reply := cache.SetExpire("book_record-"+identify, 0)
		if err != nil {
			logs.Error("Set book_record-"+identify+" expire:",err.Error(), reply)
		}
	} else {
		c.JsonResult(404, "标识为"+identify+"的书不存在")
	}

	if err = new(models.Book).DeleteBook(identify); err != nil {
		logs.Error("DeleteBooksByIdentify: DeleteBook ", err.Error())
		c.JsonResult(500, "服务器内部错误，删除失败")
	}
	err, reply:=cache.SetExpire("book-"+identify, 0)
	if err != nil {
		logs.Error("Set book-"+identify+" expire:",err.Error(), reply)
	}
	if err = new(models.BookCategory).DeleteBookCategories(identify); err != nil {
		logs.Error("DeleteBooksByIdentify: DeleteBookCategories ", err.Error())
		c.JsonResult(500, "服务器内部错误，删除图书分类失败")
	}
	c.JsonResult(200, "删除成功")
}

func (c *BookController) LendBookById() {
	logs.Debug("Enter LendBookById")
	book_id, err := c.GetInt(":book_id")
	if err != nil {
		c.JsonResult(404, err.Error())
	}
	var book_record models.BookRecord
	book_record.Id = book_id
	if err = book_record.GetBookById(); err != nil {
		logs.Error("LendBookById: GetBookById ", err.Error())
		c.JsonResult(404, "id为"+strconv.Itoa(book_id)+"图书不存在")
	}
	if book_record.LendStatus != 0 {
		c.JsonResult(400, "该书已被借出")
	}
	book_record.LendCount++
	book_record.LendStatus = 1
	book_record.LendTime=time.Now()
	book_record.UserId=c.Muser.Id
	if err = book_record.UpdateBookRecordById(book_id, "lend_count", "lend_status", "lend_time", "user_id"); err != nil {
		logs.Error("LendBookById: UpdateBookRecordById ", err.Error())
		c.JsonResult(500, "借书失败")
	}
	err, reply := cache.SetExpire("book_record-"+book_record.Identify, 0)
	if err != nil {
		logs.Error("Set book_record-"+book_record.Identify+" expire:",err.Error(), reply)
	}
	url := "/bookdetail/"+book_record.Identify
	c.Success("借书成功", url, 1)
}

func (c *BookController) ReturnBookById() {
	book_id, err := c.GetInt(":book_id")
	if err != nil {
		c.JsonResult(404, err.Error())
	}
	var book_record models.BookRecord
	book_record.Id = book_id
	if err = book_record.GetBookById(); err != nil {
		logs.Error("ReturnBookById: GetBookById ", err.Error())
		c.JsonResult(400, "id为"+strconv.Itoa(book_id)+"图书不存在")
	}
	if book_record.LendStatus != 1 {
		c.JsonResult(400, "操作有误")
	}
	book_record.LendStatus = 0
	book_record.ReturnTime = time.Now()
	book_record.UserId = -1
	if err = book_record.UpdateBookRecordById(book_id, "lend_status", "return_time", "user_id"); err != nil {
		logs.Error("ReturnBookById: UpdateBookRecordById ", err.Error())
		c.JsonResult(500, "还书失败")
	}
	err, reply := cache.SetExpire("book_record-"+book_record.Identify, 0)
	if err != nil {
		logs.Error("Set book_record-"+book_record.Identify+" expire:",err.Error(), reply)
	}
	//c.JsonResult(200, "还书成功")
	url := "/bookdetail/"+book_record.Identify
	c.Success("还书成功", url, 1)
}