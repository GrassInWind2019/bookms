package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func init() {
	//orm.RegisterModel(new(User))
	orm.RegisterModelWithPrefix("bookms_user_", new(User))
	orm.RegisterModelWithPrefix("bookms_", new(Book))
	orm.RegisterModelWithPrefix("bookms_", new(Category))
	orm.RegisterModelWithPrefix("bookms_", new(BookCategory))
	orm.RegisterModelWithPrefix("bookms_", new(BookRecord))
}

func TNUser() string {
	return "bookms_user_user"
}

func TNBook() string {
	return "bookms_book"
}

func TNBookCategory() string {
	return "bookms_book_category"
}

func TNCategory() string {
	return "bookms_category"
}

func TNBookIdentify() string {
	return "bookms_book_identify"
}

func TNBookRecord() string {
	return "bookms_book_record"
}

func GetOrm(alias string) orm.Ormer {
	o := orm.NewOrm()
	if len(alias) > 0 {
		logs.Debug("Using alias: " + alias)
		if "w" == alias {
			o.Using("default")
		} else {
			o.Using(alias)
		}
	}
	return o
}