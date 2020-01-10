package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func init() {
	//orm.RegisterModel(new(User))
	orm.RegisterModelWithPrefix("bookms_user_", new(User))
}

func TNUser() string {
	return "bookms_user_user"
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