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
