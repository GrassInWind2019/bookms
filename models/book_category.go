package models

import (
	"errors"
)

type BookCategory struct {
	Id int `orm:"pk;auto" json:"id"`
	//BookId int `orm:"column(book_id);size(15)" json:"book_id"`
	Identify string `orm:"size(100)" json:"identify"`
	CategoryId int `orm:"column(category_id);size(8)" json:"category_id"`
}

func (m *BookCategory) SetBookCategories(identify string, cids []int) error {
	if "" == identify || 0 == len(cids) {
		return errors.New("Invalid arguments")
	}
	o := GetOrm("w")
	var bookCates []BookCategory
	for _, id := range cids {
		bookCategory := BookCategory{
			Identify:identify,
			CategoryId:id,
		}
		bookCates = append(bookCates, bookCategory)
	}
	var err error
	o.Begin()
	defer func() {
		if err != nil {
			o.Rollback()
		} else {
			o.Commit()
			go UpdateCategoryCount()
		}
	}()
	if _, err = o.QueryTable(TNBookCategory()).Filter("identify", identify).Delete(); err != nil {
		return err
	}
	_, err = o.InsertMulti(len(bookCates), &bookCates)
	return err
}

func (m *BookCategory) GetBookCategories(identify string) (cids []int, err error) {
	if "" == identify {
		err = errors.New("Invalid arguments")
		return
	}
	var bookCates []BookCategory
	o := GetOrm("r")
	if _, err = o.QueryTable(TNBookCategory()).Filter("identify", identify).All(&bookCates); err != nil {
		return
	}
	for _, bookCategory := range bookCates {
		cids = append(cids, bookCategory.CategoryId)
	}
	return
}

func (m *BookCategory) DeleteBookCategories(identify string) error {
	if "" == identify {
		return errors.New("Invalid arguments")
	}
	o := GetOrm("w")
	_, err := o.QueryTable(TNBookCategory()).Filter("identify", identify).Delete()
	go UpdateCategoryCount()
	return err
}

//func (m *BookCategory) SetBookCategories(bookId int, cIds []int) error {
//	if len(cIds) == 0 || bookId < 0 {
//		return errors.New("Invalid arguments")
//	}
//	o := GetOrm("w")
//	var categories []Category
//	o.QueryTable(TNCategory()).Filter("id__in", cIds).All(&categories, "id", "pid")
//	cidMap := make(map[int]bool)
//	for _, category := range categories {
//		cidMap[category.Id] = true
//		cidMap[category.Pid] = true
//	}
//	var bookCategories []BookCategory
//	for _, cId := range cIds {
//		if cidMap[cId] == false {
//			errStr := strconv.Itoa(cId)+"分类不存在"
//			return errors.New(errStr)
//		}
//		bookCategory := BookCategory{
//			//BookId: bookId,
//			CategoryId:cId,
//		}
//		bookCategories = append(bookCategories, bookCategory)
//	}
//	o.QueryTable(TNBookCategory()).Filter("book_id", bookId).Delete()
//	_, err := o.InsertMulti(len(bookCategories), &bookCategories)
//	go UpdateCategoryCount()
//	return err
//}
//
//func (m *BookCategory) GetBookCategoryById(book_id int) (category_ids []int, err error) {
//	if book_id < 0 {
//		err = errors.New("book id is invalid")
//		return
//	}
//	o := GetOrm("r")
//	//m.BookId = book_id
//	if err = o.Read(m, "book_id"); err == nil {
//		category_ids = append(category_ids, m.CategoryId)
//	}
//	return
//}