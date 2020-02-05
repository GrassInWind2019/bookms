package models

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"sync"
)

type Category struct {
	Id int `orm:"pk;auto" json:"id"`
	Pid int `json:"pid"`
	CategoryName string `orm:"column(category_name);size(100);unique" json:"category_name"`
	Description string `orm:"size(500)" json:"description"`
	Icon string `orm:"size(500)" json:"icon"`
	BookCount int `orm:"column(book_count)" json:"book_count"`
	Sort int `json:"sort"`
	Status int `orm:"size(4)" json:"status"`
}

var (
	working = false
	mutex sync.Mutex
)

func (m *Category) GetTopCategories() (categories []Category, err error) {
	//select id from Category where pid = 0
	o := GetOrm("r")
	_, err = o.QueryTable(TNCategory()).Filter("pid", 0).All(&categories)

	return
}

func (m *Category) GetCategoriesByPid(pid int) (categories []Category, err error) {
	if pid < 0 {
		err = errors.New("Invalid arguments")
		return
	}
	o := GetOrm("r")
	_, err = o.QueryTable(TNCategory()).Filter("pid", pid).All(&categories)

	return
}

func (m *Category) AddCategory() (err error){
	if m.Pid < 0 {
		err = errors.New("invalid arguments")
		return
	}
	o := GetOrm("w")
	category := Category{
		CategoryName:m.CategoryName,
	}
	if o.Read(&category); category.Id == 0 {
		_, err = o.Insert(m)
	} else {
		errStr := "分类名 " + m.CategoryName + " 已存在"
		err = errors.New(errStr)
	}
	return
}

func (m *Category) DeleteCategory() (err error) {
	if m.Id < 0 {
		err = errors.New("Invalid arguments")
		return
	}
	category := Category{
		Id:m.Id,
	}
	o := GetOrm("w")
	if err = o.Read(&category); category.BookCount > 0 {
		err = errors.New("删除失败，当前分类下的图书不为0，不允许删除")
		return
	}
	if _, err = o.Delete(&category, "id"); err != nil {
		return
	}
	_, err = o.QueryTable(TNCategory()).Filter("pid", m.Id).Delete()
	return
}

func (m *Category) UpdateCategory(field, val string) error {
	o := GetOrm("w")
	_ ,err := o.QueryTable(TNCategory()).Filter("id", m.Id).Update(orm.Params{field:val})
	return err
}

func UpdateCategoryCount() {
	//select count("book_id") as cnt,category_id from bookms_category group by category_id
	mutex.Lock()
	defer mutex.Unlock()
	if working == true {
		return
	}
	working = true
	defer func() {
		working = false
	}()
	type Count struct {
		Cnt int
		Cid int
	}
	var count []Count
	o := GetOrm("w")
	//select sum(b.book_count) as cnt, c.category_id as cid from bookms_book b left join bookms_book_category c on b.identify=c.identify group by c.category_id 不会使用索引
	//select sum(b.book_count) as cnt, c.category_id as cid from bookms_book_category c left join bookms_book b on b.identify=c.identify group by c.category_id 使用索引
	sql := "select sum(b.book_count) as cnt, c.category_id as cid from " + TNBookCategory() + " c left join " + TNBook() + " b on b.identify=c.identify group by c.category_id"
	o.Raw(sql).QueryRows(&count)
	logs.Debug("count result: ", count)
	if len(count) == 0 {
		return
	}
	var categories []Category
	o.QueryTable(TNCategory()).All(&categories, "id", "pid", "book_count")
	if len(categories)==0 {
		return
	}

	var err error
	o.Begin()
	defer func() {
		if err != nil {
			o.Rollback()
		} else {
			o.Commit()
		}
	}()
	o.QueryTable(TNCategory()).Update(orm.Params{"book_count":0})

	for _,v := range count {
		if v.Cnt > 0 {
			logs.Debug("id="+strconv.Itoa(v.Cid)+",book_count="+strconv.Itoa(v.Cnt))
			_, err = o.QueryTable(TNCategory()).Filter("id", v.Cid).Update(orm.Params{"book_count":v.Cnt})
			if err != nil {
				logs.Error(err.Error())
				return
			}
		}
	}
}