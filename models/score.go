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

func (m *Score) AddScore() (average_score int,err error)  {
	if "" == m.Identify || m.UserId <= 0 || m.Score <= 0 || m.Score > 5 {
		return 0, errors.New("Invalid arguments")
	}
	m.Score = m.Score*10
	o1 := GetOrm("uw")
	o1.Begin()
	defer func() {
		if err != nil {
			o1.Rollback()
		} else {
			o1.Commit()
		}
	}()
	//select id from bookms_score where user_id=? and identify=?
	score := Score{
		Identify:m.Identify,
		UserId:m.UserId,
	}
	//TODO:check errors except "no row found"
	o1.Read(&score, "user_id", "identify")
	if score.Id > 0 {
		return 0, errors.New("The user had already scored the book!")
	}
	_, err = o1.Insert(m)
	o2 := GetOrm("w")
	o2.Begin()
	defer func() {
		if err != nil {
			o2.Rollback()
		} else {
			o2.Commit()
		}
	}()
	book := Book{
		Identify:m.Identify,
	}
	err = o2.Read(&book, "identify")
	if err != nil {
		return 0, err
	}
	if 0 == book.ScoreCount {
		book.AverageScore = m.Score
		book.ScoreCount = 1
	} else {
		book.AverageScore = (book.AverageScore*book.ScoreCount+m.Score)/(book.ScoreCount+1)
		book.ScoreCount = book.ScoreCount + 1
	}
	_, err = o2.Update(&book, "average_score", "score_count")
	return book.AverageScore,err
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
