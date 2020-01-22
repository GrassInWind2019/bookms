package models

import (
	"errors"
	"strconv"
)

type Favorite struct {
	Id int `orm:"pk;auto" json:"id"`
	UserId int `orm:"column(user_id)" json:"user_id"`
	Identify string `orm:"size(100)" json:"identify"`
}

func (m *Favorite) FavoriteDo() error {
	if m.UserId <= 0 || "" == m.Identify || m.Id > 0 {
		return errors.New("Invalid arguments")
	}
	o := GetOrm("uw")
	o.Read(m, "user_id", "identify")
	if m.Id > 0 { //取消收藏
		_, err := o.QueryTable(TNFavorite()).Filter("user_id", m.UserId).Filter("identify", m.Identify).Delete()
		if err != nil {
			return err
		}
		o:= GetOrm("w")
		//update bookms_book set favorite_count=favorite_count-1 where identify=?
		sql := "update " + TNBook() + " set favorite_count=favorite_count-1 where identify=" + m.Identify
		_, err = o.Raw(sql).Exec()
		return err
	} else { //添加收藏
		_, err := o.Insert(m)
		if err != nil {
			return err
		}
		o:= GetOrm("w")
		//update bookms_book set favorite_count=favorite_count+1 where identify=?
		sql := "update " + TNBook() + " set favorite_count=favorite_count+1 where identify=" + m.Identify
		_, err = o.Raw(sql).Exec()
		return err
	}

	return nil
}

func (m *Favorite) IsFavorite() (res bool, err error) {
	if "" == m.Identify || m.UserId <= 0 || m.Id > 0 {
		return false, errors.New("Invalid arguments")
	}
	o := GetOrm("ur")
	if err = o.Read(m, "user_id", "identify"); err != nil {
		return false, err
	}
	if m.Id > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (m *Favorite) ListFavoriteByUserId() (books []Book, err error) {
	if m.UserId <= 0 {
		err = errors.New("Invalid argument")
		return
	}
	o := GetOrm("ur")
	sql := "select identify from " + TNFavorite() + " where user_id=" + strconv.Itoa(m.UserId)
	var identifies []string
	_,err = o.Raw(sql).QueryRows(&identifies)
	if err != nil {
		return
	}
	o = GetOrm("r")
	if _, err = o.QueryTable(TNBook()).Filter("identify__in", identifies).All(&books); err != nil {
		return
	}
	return
}