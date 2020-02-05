package models

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"time"
)

type BookRecord struct {
	Id int `orm:"pk;auto" json:"id"`
	//BookId int `orm:"column(book_id)" json:"book_id"`
	Identify string `orm:"size(100)" json:"identify"`
	LendStatus int `orm:"column(lend_status);default(0)" json:"lend_status"`  //状态:0 在馆 ; 1 已借出;  5-->已下架
	UserId int `orm:"column(user_id);size(20)" json:"user_id"`  //借阅人
	LendTime time.Time `orm:"column(lend_time);type(datetime);not null" json:"lend_time"`  //借出时间
	ReturnTime time.Time `orm:"column(return_time);type(datetime);not null" json:"return_time"` //归还时间
	LendCount int `orm:"column(lend_count);type(int);default(0)" json:"lend_count"` //借出次数
	StorePosition string `orm:"column(store_position);size(200);not null" json:"store_position"` //存放位置
}

func (m *BookRecord) AddBookRecords(book_count int, identify string, positions []string) error {
	if "" == identify || book_count == 0 || len(positions) == 0 || book_count != len(positions) {
		return errors.New("Invalid arguments")
	}
	o := GetOrm("w")
	//select count(*) from bookms_book_record where identify=
	//select identify from bookms_book_record where identify=
	//if res := o.QueryTable(TNBookRecord()).Filter("identify", identify).Exist(); res == true {
	//	return errors.New("标识已存在，请更换一个新的标识")
	//}
	sql := "select count(*) from " + TNBookRecord() + " where identify=" + identify
	//sql := "select identify from " + TNBookRecord() + " where identify=" + identify
	logs.Debug(sql)
	type Count struct {
		Cnt int
	}
	var count Count
	if err := o.Raw(sql).QueryRow(&count); err != nil {
		return err
	}
	if count.Cnt > 0 {
		return errors.New("标识已存在，请更换一个新的标识")
	}
	var bookRecords []BookRecord
	for i := 0; i < book_count; i++ {
		record := BookRecord{
			Identify:identify,
			LendStatus:0,
			UserId:-1,
			LendTime:time.Now(),
			ReturnTime:time.Now(),
			LendCount:0,
			StorePosition:positions[i],
		}
		bookRecords = append(bookRecords, record)
	}
	_, err := o.InsertMulti(book_count, &bookRecords)

	return err
}

func (m *BookRecord) GetBookIdsByIdentify(identify string) (ids []int, err error) {
	if identify == "" {
		err = errors.New("Invalid arguments")
		return
	}
	o := GetOrm("r")
	var bookRecords []BookRecord
	//if _ ,err = o.QueryTable(TNBookRecord()).Filter("identify", identify).All(&bookRecords); err != nil {
	//	return
	//}
	sql := "select * from "+TNBookRecord()+" where identify="+identify
	logs.Debug(sql)
	if _, err = o.Raw(sql).QueryRows(&bookRecords); err != nil {
		return
	}
	for _, bookRecord := range bookRecords {
		ids = append(ids, bookRecord.Id)
	}
	return
}

func (m *BookRecord) GetBookRecordsByIdentify(identify string) (bookRecords []BookRecord, err error) {
	if identify == "" {
		err = errors.New("Invalid arguments")
		return
	}
	o := GetOrm("r")
	//if _ ,err = o.QueryTable(TNBookRecord()).Filter("identify", identify).All(&bookRecords); err != nil {
	//	return
	//}
	sql := "select * from "+TNBookRecord()+" where identify='"+identify+"'"
	logs.Debug(sql)
	if _, err = o.Raw(sql).QueryRows(&bookRecords); err != nil {
		return
	}
	return
}

func (m *BookRecord) GetBookById() error {
	if m.Id <= 0 {
		return errors.New("Invalid book id")
	}
	o := GetOrm("r")
	err := o.Read(m, "id")
	return err
}

func (m *BookRecord) DeleteBookRecordsByIds(ids []int) error {
	if len(ids) == 0 {
		return errors.New("Id is null")
	}
	o := GetOrm("w")
	_, err := o.QueryTable(TNBookRecord()).Filter("book_id__in", ids).Delete()
	return err
}

func (m *BookRecord) UpdateBookRecordById(book_id int, fields ...string) error {
	if book_id <= 0 || 0 == len(fields) {
		return errors.New("Invalid arguments")
	}
	o := GetOrm("w")
	bookRecord := BookRecord{
		Id:book_id,
	}
	if err := o.Read(&bookRecord); err != nil {
		return errors.New(err.Error())
	}
	_, err := o.Update(m, fields...)
	return err
}