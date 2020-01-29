package models

import (
	"errors"
	"time"
)

type Score struct {
	Id int `orm:"pk;auto" json:"id"`
	Identify string `orm:"size(100)" json:"identify"`
	UserId int `orm:"column(user_id)" json:"user_id"`
	Score int `json:"score"`
	CreateTime time.Time `orm:"column(create_time)" json:"create_time"`
}

func (m *Score) AddScore() error  {
	if "" == m.Identify || m.UserId <= 0 || m.Score <= 0 || m.Score > 5 {
		return errors.New("Invalid arguments")
	}
	m.Score = m.Score*10
	o := GetOrm("uw")
	//select id from bookms_score where user_id=? and identify=?
	score := Score{
		Identify:m.Identify,
		UserId:m.UserId,
	}
	//TODO:check errors except "no row found"
	o.Read(&score, "user_id", "identify")
	if score.Id > 0 {
		return errors.New("The user had already scored the book!")
	}
	_, err := o.Insert(m)
	o = GetOrm("w")
	book := Book{
		Identify:m.Identify,
	}
	err = o.Read(&book, "identify")
	if err != nil {
		return err
	}
	if 0 == book.ScoreCount {
		book.AverageScore = m.Score
		book.ScoreCount = 1
	} else {
		book.AverageScore = (book.AverageScore*book.ScoreCount+m.Score)/(book.ScoreCount+1)
		book.ScoreCount = book.ScoreCount + 1
	}
	_, err = o.Update(&book, "average_score", "score_count")
	return err
}

func (m *Score) GetBookScore() (score Score, err error) {
	if "" == m.Identify || m.UserId <= 0 {
		err = errors.New("Invalid arguments")
		return
	}
	o := GetOrm("ur")
	err = o.Read(m, "user_id", "identify")
	if err == nil {
		score = *m
	}
	return
}
