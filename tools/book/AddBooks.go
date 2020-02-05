package main

import (
	bookms_init "bookms/init"
	"bookms/models"
	"fmt"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"strconv"
	"time"
)
//slow
func AddMultiBooks(add_count int) {
	for i:=0; i<add_count; i++ {
		book_name := "test "+strconv.Itoa(i)
		description := "test description "+strconv.Itoa(i)
		catalog := "test catalog "+strconv.Itoa(i)
		sort := rand.Intn(9)+2
		author := "test author "+strconv.Itoa(i)
		book_count := 3
		category_id := rand.Intn(6)+6
		store_position := "test store_position "+strconv.Itoa(i)
		//get an uuid for book by sonyflake
		sonyflakeObj := bookms_init.GetSonyFlakeObj()
		uuid, err := sonyflakeObj.NextID()
		if err != nil {
			logs.Error(500, "UUID create failed! ", err)
		}
		identify := fmt.Sprint(uuid)
		if "" == identify {
			logs.Error(500, "UUID create failed!")
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
			logs.Error(500, err)
		}
		var book_category models.BookCategory
		var cids []int
		cids = append(cids, category_id)
		if err = book_category.SetBookCategories(identify, cids); err != nil {
			logs.Error(500, err)
		}
		var positions []string
		for i:=0;i<book_count;i++ {
			position := store_position
			positions = append(positions, position)
		}
		if err = new(models.BookRecord).AddBookRecords(book_count, identify, positions); err != nil {
			logs.Error(500, err)
		}
	}
}

//fast with not set category
func AddMultiBooksFast(add_count int) {
	for i:=0; i<add_count; i++ {
		if i % 5000 == 0 {
			fmt.Println()
			logs.Info("第",i,"c次")
			fmt.Println()
		}
		book_name := "test "+strconv.Itoa(i)
		description := "test description "+strconv.Itoa(i)
		catalog := "test catalog "+strconv.Itoa(i)
		sort := rand.Intn(9)+2
		author := "test author "+strconv.Itoa(i)
		book_count := 3
		//category_id := rand.Intn(6)+6
		store_position := "test store_position "+strconv.Itoa(i)
		//get an uuid for book by sonyflake
		sonyflakeObj := bookms_init.GetSonyFlakeObj()
		uuid, err := sonyflakeObj.NextID()
		if err != nil {
			logs.Error(500, "UUID create failed! ", err)
		}
		identify := fmt.Sprint(uuid)
		if "" == identify {
			logs.Error(500, "UUID create failed!")
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
			logs.Error(500, err)
		}
		var positions []string
		for i:=0;i<book_count;i++ {
			position := store_position
			positions = append(positions, position)
		}
		if err = new(models.BookRecord).AddBookRecords(book_count, identify, positions); err != nil {
			logs.Error(500, err)
		}
	}
}

func main() {
	logs.SetLevel(logs.LevelInfo)
	//AddMultiBooks(500000)
	t1 := time.Now()
	AddMultiBooksFast(10000);
	cost := time.Since(t1)
	logs.Info("cost: ",cost)
}
