package models

import (
	"errors"
	"fmt"
	"time"
)

type Comments struct {
	Id int `orm:"pk;auto" json:"id"`
	Identify string `orm:"size(100)" json:"identify"`
	UserId int `orm:"column(user_id)" json:"user_id"`
	Content string `orm:"size(3000)" json:"content"`
	CreateTime time.Time `orm:"column(create_time)" json:"create_time"`
}

func (m *Comments) AddComment() error {
	if "" == m.Identify || m.UserId <= 0 || "" == m.Content {
		return errors.New("Invalid arguments")
	}
	m.CreateTime = time.Now()
	o := GetOrm("uw")
	_, err := o.Insert(m)
	return err
}

func (m *Comments) DeleteComment() error {
	if m.Id < 0 {
		return errors.New("Invalid argument")
	}
	o := GetOrm("uw")
	_, err := o.Delete(m, "id")
	return err
}

type BookComment struct {
	UserId int `json:"user_id"`
	Score int `json:"score"`
	Avatar string `json:"avatar"`
	NickName string `json:"nick_name"`
	Role int `json:"role"`
	RoleName string `json:"role_name"`
	Content string `json:"content"`
	CreateTime time.Time `json:"time_create"`
}

func (m *Comments) GetBookCommentsAndScores(identify string, page, page_size int) (book_comments []BookComment, err error) {
	if "" == identify {
		err = errors.New("Invalid identify")
		return
	}

	sql := "select c.content,c.create_time,s.score,u.id as user_id,u.avatar,u.nickname,u.role from bookms_user_comments c" +
		" left join bookms_user_score s on c.identify=s.identify left join bookms_user_user u on " +
		"u.id=c.user_id where c.identify='"+identify+"' and s.identify='"+identify+"' order by c.id desc limit %v offset %v"
	sql = fmt.Sprintf(sql, page_size, (page - 1)*page_size)
	o := GetOrm("ur")
	_, err = o.Raw(sql).QueryRows(&book_comments)
	return
}