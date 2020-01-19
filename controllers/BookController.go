package controllers

import (
	"bookms/models"
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
		identify := c.GetString("identify")
		if identify == "" {
			c.JsonResult(400, "标识不能为空")
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
		c.JsonResult(200, "添加成功")
	}
}

func (c *BookController) UpdateBookByIdentify() {
	c.TplName="book/updatebook.html"
	logs.Debug("Enter UpdateBookByIdentify")
	//identify := c.GetString("identify")
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
			c.JsonResult(400, err.Error())
		}
		cids, err := new(models.BookCategory).GetBookCategories(identify)
		if err != nil {
			c.JsonResult(400, err.Error())
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
		//c.JsonResult(200, "OK")
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
			c.JsonResult(500, err.Error())
		}
		book_category := models.BookCategory{
			Identify:identify,
			CategoryId:category_id,
		}
		cIds := []int{}
		cIds = append(cIds, category_id)
		err = book_category.SetBookCategories(identify, cIds)
		if err != nil {
			c.JsonResult(500, err.Error())
		}
		c.JsonResult(200, "更新成功")
	}
}
//TODO: multi tables operations
func (c *BookController) DeleteBooksByIdentify() {
	logs.Debug("Enter DeleteBooksByIdentify")
	identify := c.Ctx.Input.Param(":identify")
	if identify=="" {
		c.Abort("404")
	}
	ids, err := new(models.BookRecord).GetBookIdsByIdentify(identify)
	if err != nil {
		c.JsonResult(400, "标识不存在")
	}
	if len(ids) > 0 {
		if err = new(models.BookRecord).DeleteBookRecordsByIds(ids); err != nil {
			c.JsonResult(500, err.Error())
		}
	} else {
		logs.Info("book records didn't have identify:"+identify)
		c.JsonResult(500, "book records didn't have identify:"+identify)
	}

	if err = new(models.Book).DeleteBook(identify); err != nil {
		c.JsonResult(500, err.Error())
	}
	if err = new(models.BookCategory).DeleteBookCategories(identify); err != nil {
		c.JsonResult(500, err.Error())
	}
	c.JsonResult(200, "删除成功")
}

//delete one book record
//TODO: multi tables operations
func (c *BookController) DeleteBookById() {
	book_id_str := c.Ctx.Input.Param("book_id")
	book_id, err := strconv.ParseInt(book_id_str, 10, 64)
	if err != nil || book_id <= 0 {
		c.Abort("404")
	}
	book_record := models.BookRecord{
		Id:int(book_id),
	}
	err = book_record.GetBookById()
	if err != nil {
		c.JsonResult(400, err.Error())
	}
	if "" == book_record.Identify {
		c.JsonResult(500, "identify is null")
	}
	var identifies []string
	identifies = append(identifies, book_record.Identify)
	books, err := new(models.Book).GetBooksByIdentifies(identifies)
	if err != nil {
		c.JsonResult(500, err.Error())
	}
	books[0].BookCount -= 1
	if err = books[0].UpdateBookByIdentify("book_count"); err != nil {
		c.JsonResult(500, err.Error())
	}
	//all book information need be deleted
	if  0 == books[0].BookCount {
		ids, err := new(models.BookRecord).GetBookIdsByIdentify(book_record.Identify)
		if err != nil {
			c.JsonResult(400, "标识不存在")
		}
		if err = new(models.BookRecord).DeleteBookRecordsByIds(ids); err != nil {
			c.JsonResult(500, err.Error())
		}
		if err = new(models.Book).DeleteBook(book_record.Identify); err != nil {
			c.JsonResult(500, err.Error())
		}
		if err = new(models.BookCategory).DeleteBookCategories(book_record.Identify); err != nil {
			c.JsonResult(500, err.Error())
		}
	}

	c.JsonResult(200, "删除成功")
}

func (c *BookController) LendBookById() {
	logs.Debug("Enter LendBookById")
	book_id, err := c.GetInt(":book_id")
	if err != nil {
		c.JsonResult(400, err.Error())
	}
	var book_record models.BookRecord
	book_record.Id = book_id
	if err = book_record.GetBookById(); err != nil {
		c.JsonResult(400, err.Error())
	}
	if book_record.LendStatus != 0 {
		c.JsonResult(400, "该书已被借出")
	}
	book_record.LendCount++
	book_record.LendStatus = 1
	book_record.LendTime=time.Now()
	book_record.UserId=c.Muser.Id
	if err = book_record.UpdateBookRecordById(book_id, "lend_count", "lend_status", "lend_time", "user_id"); err != nil {
		c.JsonResult(500, err.Error())
	}
	c.JsonResult(200, "借书成功")
}

func (c *BookController) ReturnBookById() {
	book_id, err := c.GetInt(":book_id")
	if err != nil {
		c.JsonResult(400, err.Error())
	}
	var book_record models.BookRecord
	book_record.Id = book_id
	if err = book_record.GetBookById(); err != nil {
		c.JsonResult(400, err.Error())
	}
	if book_record.LendStatus != 1 {
		c.JsonResult(400, "操作有误")
	}
	book_record.LendStatus = 0
	book_record.ReturnTime = time.Now()
	book_record.UserId = -1
	if err = book_record.UpdateBookRecordById(book_id, "lend_status", "return_time", "user_id"); err != nil {
		c.JsonResult(500, err.Error())
	}
	c.JsonResult(200, "还书成功")
}