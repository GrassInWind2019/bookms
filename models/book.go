package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type Book struct {
	Id int `orm:"pk;auto" json:"id"`
	BookName string `orm:"column(book_name);size(200)" json:"book_name"`
	Identify   string `orm:"size(100);unique" json:"identify"`
	Description string `orm:"size(1000)" json:"description"`
	Catalog    string `orm:"size(6000)" json:"catalog"` //目录
	Cover string `orm:"size(1000)" json:"cover"`
	Status int `orm:"default(0)" json:"status"`  //状态:0 正常 ; 1 已下架
	Sort int `orm:"type(int);default(0)" json:"sort"`
	CreateTime time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	DocCount int `orm:"column(doc_count)" json:"doc_count"`
	CommentCount int `orm:"type(int)" json:"comment_count"`
	FavoriteCount int `orm:"column(favorite_count);type(int);default(0)" json:"favorite_count"`  //收藏次数
	AverageScore int `orm:"column(average_score);default(0)" json:"average_score"`
	ScoreCount int `orm:"column(score_count);default(0)" json:"score_count"`
	Author string `orm:"size(100)" json:"author"`
	BookCount int `orm:"column(book_count)" json:"book_count"`
}

func (m *Book) GetBooksByCategory(category, page, pagesize int, fields ...string) (books []Book, count int, err error) {
	if len(fields) == 0 {
		fields = append(fields, "book_name", "identify", "cover", "status", "create_time", "author")
	}
	fieldStr := "b."+strings.Join(fields, ",b.")
	//select * from bookms_book b left join bookms_book_category c on b.id = c.book_id where b.status=0 and c.category_id=1 order by sort limit 0,10
	sqlFmt := "select %v from " + TNBook() + " b left join " + TNBookCategory() +
	" c on b.identify=c.identify where c.category_id=" + strconv.Itoa(category)
	sql := fmt.Sprintf(sqlFmt, fieldStr)
	sql = sql + " limit %d,%d"
	sql = fmt.Sprintf(sql, (page-1)*pagesize,pagesize)
	sqlCount := fmt.Sprintf(sqlFmt, "count(*) cnt")
	logs.Debug("GetBooksByCategory: "+sql)
	logs.Debug("GetBooksByCategory: "+sqlCount)
	o := GetOrm("r")
	var params []orm.Params
	if _, err = o.Raw(sqlCount).Values(&params); err == nil {
		if len(params) > 0 {
			count, _ = strconv.Atoi(params[0]["cnt"].(string))
		}
	}
	_, err = o.Raw(sql).QueryRows(&books)
	return
}

func (m *Book) GetBooksByCategory2(category, page, page_size int, fields ...string) (books []Book, count int, err error) {
	if len(fields) == 0 {
		fields = append(fields, "id","book_name", "identify", "cover","create_time", "status", "author")
	}
	fieldStr := strings.Join(fields, ",")
	sql := "select identify from "+TNBookCategory()+" where category_id=" + strconv.Itoa(category) + " limit 0,1000"
	//logs.Debug(sql)
	var identifies []string
	o := GetOrm("r")
	if _,err = o.Raw(sql).QueryRows(&identifies); err != nil {
		return
	}
	identifies_str := "'"
	identifies_str += strings.Join(identifies, "','")
	identifies_str += "'"
	sql = "select count(*) cnt from "+TNBook()+" where identify in ("+ identifies_str+ ")"
	//logs.Debug(sql)
	var params []orm.Params
	if _,err = o.Raw(sql).Values(&params); err != nil {
		return
	}
	if len(params) > 0 {
		count, _ = strconv.Atoi(params[0]["cnt"].(string))
	} else {
		return
	}
	sql = "select %v from "+TNBook()+ " where identify in ("+identifies_str+") order by id limit %d,%d"
	sql = fmt.Sprintf(sql, fieldStr, (page-1)*page_size, page_size)
	//logs.Debug(sql)
	_, err = o.Raw(sql).QueryRows(&books)
	return
}

func (m *Book) GetBooksByIdentifies(identifies []string) (books []Book, err error) {
	o := GetOrm("r")
	_, err = o.QueryTable(TNBook()).Filter("identify__in", identifies).All(&books)
	return
}

func (m *Book) AddBook() error {
	o := GetOrm("w")

	if "" == m.Identify {
		return errors.New("Identify cannot be null")
	}
	if res := o.QueryTable(TNBook()).Filter("identify", m.Identify).Exist(); res == true {
		return errors.New("Identify "+m.Identify+" already exist!")
	}

	_, err := o.Insert(m)
	return err
}

func (m *Book) DeleteBook(identify string) error {
	if "" == identify {
		return errors.New("Invalid arguments")
	}
	o := GetOrm("w")
	_, err := o.QueryTable(TNBook()).Filter("identify", identify).Delete()
	go UpdateCategoryCount()
	return err
}

//func (m *Book) UpdateBookById(fields ...string) error {
//	if m.Id < 0 || len(fields) == 0 {
//		return errors.New("Invalid arguments")
//	}
//
//	o := GetOrm("w")
//	book := Book{
//		Id:m.Id,
//	}
//	if err := o.Read(&book);err != nil {
//		errStr := "Id为"+strconv.Itoa(m.Id)+"的书不存在"
//		return errors.New(errStr)
//	}
//	_, err := o.Update(m, fields...)
//	return err
//}

func (m *Book) UpdateBookByIdentify(fields ...string) error {
	if "" == m.Identify || 0 == len(fields) {
		return errors.New("Invalid arguments")
	}
	o := GetOrm("w")
	book := Book{
		Identify:m.Identify,
	}
	if err := o.QueryTable(TNBook()).Filter("identify", m.Identify).One(&book); err != nil {
		return errors.New(m.Identify+" not exist")
	}
	m.Id = book.Id
	_, err := o.Update(m, fields...)
	return err
}

func (m *Book) SearchBook(keyword string, page, page_size int) (books []Book, cnt int, err error) {
	if "" == keyword {
		err = errors.New("keyword cannot been null")
		return
	}
	o := GetOrm("r")
	sql := "select count(*) cnt from "+ TNBook()+" where book_name like '%"+keyword+"%' or description like '%"+keyword+"%'"
	logs.Debug(sql)
	cnt = 0
	var param []orm.Params
	_,err = o.Raw(sql).Values(&param)
	if err != nil {
		return
	}

	if len(param) > 0 {
		cnt, _ = strconv.Atoi(param[0]["cnt"].(string))
	} else {
		cnt = 0
		err = errors.New("cannot find any books")
		return
	}

	sql = "select * from "+ TNBook()+" where book_name like '%"+keyword+"%' or description like '%"+keyword+"%'"
	sql2 := fmt.Sprintf(" limit %v,%v", (page-1)*page_size, page_size)
	sql = sql+sql2
	logs.Debug(sql)
	_,err = o.Raw(sql).QueryRows(&books)
	return
}