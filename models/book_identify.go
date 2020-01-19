package models

import "errors"

type BookIdentify struct {
	Id int `orm:"pk;auto" json:"id"`
	BookId int `orm:"column(book_id)" json:"book_id"`
	Identify string `orm:"size(100)" json:"identify"`
}

func (m *BookIdentify) AddBookIdentifyByIds(ids []int, identify string) error {
	if len(ids) == 0 || identify == "" {
		return errors.New("Invalid arguments")
	}
	o := GetOrm("w")
	if res := o.QueryTable(TNBookIdentify()).Filter("identify", identify).Exist(); res == true {
		return errors.New("标识已存在，请更换一个新的标识")
	}
	var bookIdentifies []BookIdentify
	for _, id := range ids {
		bookIdentify := BookIdentify{
			BookId:id,
			Identify:identify,
		}
		bookIdentifies = append(bookIdentifies, bookIdentify)
	}
	_, err := o.InsertMulti(len(bookIdentifies),&bookIdentifies)
	return err
}

func (m *BookIdentify) GetBookIdsByIdentify(identify string) (ids []int, err error) {
	if identify == "" {
		err = errors.New("Invalid arguments")
		return
	}
	o := GetOrm("r")
	var bookIdentifies []BookIdentify
	if _ ,err = o.QueryTable(TNBookIdentify()).Filter("identify", identify).All(&bookIdentifies); err != nil {
		return
	}
	for _, bookIdentify := range bookIdentifies {
		ids = append(ids, bookIdentify.BookId)
	}
	return
}

func (m *BookIdentify) DeleteBookIdentifyByIds(ids []int) error {
	if len(ids) == 0 {
		return errors.New("Id is null")
	}
	o := GetOrm("w")
	_, err := o.QueryTable(TNBookIdentify()).Filter("book_id__in", ids).Delete()
	return err
}