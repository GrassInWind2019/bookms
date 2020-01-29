package models

import (
	"errors"
	"bookms/utils"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

type User struct {
	Id int `orm:"pk;auto" json:"id"`
	Account string `orm:"size(20);unique" json:"account"`
	Nickname string `orm:"size(40)" json:"nickname"`
	Password string `json:"-"`
	Phone    string `orm:"size(15);default('')" json:"phone"`
	Email    string `orm:"size(100);unique" json:"email"`
	Role     int `orm:"default(2)" json:"role"`
	RoleName string `orm:"-" json:"role_name"`
	Avatar   string `orm:default('') json:"avatar"`
	Status   int `orm:default(0) json:"status"`
	CreateTime time.Time `orm:"type(datetime);auto_now_add" json:"create_time"`
	LastLoginTime time.Time `orm:"type(datetime)" json:"last_login_time"`
	Biography string `orm:size(500) json:"biography"`  //人物簡介
}

func getRoleName(role int) string {
	if role == 0 {
		return "超级管理员"
	} else if role == 1 {
		return "管理员"
	} else {
		return "普通用户"
	}
}

func (m *User) IsAdmin() bool {
	if 0 == m.Role {
		return true
	}
	return false
}

func (m *User) Add() error {
	if strings.Count(m.Nickname, "") > 20 {
		return errors.New("昵称长度不能超过20个字")
	}
	cond := orm.NewCondition().Or("account", m.Account).Or("email", m.Email).Or("phone", m.Phone)
	var u User
	o := GetOrm("uw")
	//o := orm.NewOrm()
	//o.Using("uw")
	//if o.QueryTable(TNUser()).SetCond(cond).One(&u, "account", "email", "phone"); u.Id > 0 {
	if o.QueryTable(TNUser()).SetCond(cond).One(&u, "account", "email", "phone"); u.Id > 0 {
		if u.Account == m.Account {
			return errors.New("用户已存在")
		}
		if u.Email == m.Email {
			return errors.New("邮箱已被注册，请换一个邮箱。")
		}
		if u.Phone == m.Phone {
			return errors.New("手机号已被注册，请换一个手机号。")
		}
	}
	encryptPassword, err := utils.EncryptPassword(m.Password)
	if err != nil {
		return err
	}
	m.Password = encryptPassword
	_, err = o.Insert(m)
	if err != nil {
		return err
	}
	m.RoleName = getRoleName(m.Role)
	return nil
}

func (m *User) Update(cols ...string) error {
	if _, err := GetOrm("uw").Update(m, cols...); err != nil {
		return err
	}
	return nil
}

func (m *User)Login(account, password string) (*User, error)  {
	m.Account = account
	if err := GetOrm("ur").Read(m, "account"); err != nil {
		return m, err
	}
	res, err := utils.VerifyPassword(m.Password, password)
	if err != nil {
		return m, err
	}
	if res == false {
		return m, errors.New("账号或密码错误")
	}
	return m, nil
}

func (m *User) Find(id int) (*User, error) {
	m.Id = id
	if err := GetOrm("ur").Read(m); err != nil {
		return m, err
	}
	m.RoleName = getRoleName(m.Role)
	return m, nil
}