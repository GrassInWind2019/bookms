package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
)

type Favorite struct {
	Id int `orm:"pk;auto" json:"id"`
	UserId int `orm:"column(user_id)" json:"user_id"`
	Identify string `orm:"size(100)" json:"identify"`
}

type UserFavorite struct {
	Id int `json:"id"`
	UserId int `json:"user_id"`
	Identify string `json:"identify"`
	BookName string `json:"book_name"`
	Cover string `json:"cover"`
	Author string `json:"author"`
	CategoryId int `json:"category_id"`
	CategoryName string `json:"category_name"`
	LendStatus string `json:"lend_status"`   //0 -->可借 1-->正在借阅 2-->不可借 5-->已下架
}

func (m *Favorite) FavoriteDo() (err error) {
	if m.UserId <= 0 || "" == m.Identify || m.Id > 0 {
		return errors.New("Invalid arguments")
	}
	o1 := GetOrm("uw")
	o1.Begin()
	defer func() {
		if err != nil {
			o1.Rollback()
		} else {
			o1.Commit()
		}
	}()
	o1.Read(m, "user_id", "identify")
	if m.Id > 0 { //取消收藏
		_, err = o1.QueryTable(TNFavorite()).Filter("user_id", m.UserId).Filter("identify", m.Identify).Delete()
		if err != nil {
			return err
		}
		o2:= GetOrm("w")
		o2.Begin()
		defer func() {
			if err != nil {
				o2.Rollback()
			} else {
				o2.Commit()
			}
		}()
		//update bookms_book set favorite_count=favorite_count-1 where identify=?
		sql := "update " + TNBook() + " set favorite_count=favorite_count-1 where identify='" + m.Identify+"'"
		_, err = o2.Raw(sql).Exec()
		return err
	} else { //添加收藏
		_, err = o1.Insert(m)
		if err != nil {
			return err
		}
		o2:= GetOrm("w")
		o2.Begin()
		defer func() {
			if err != nil {
				o2.Rollback()
			} else {
				o2.Commit()
			}
		}()
		//update bookms_book set favorite_count=favorite_count+1 where identify=?
		sql := "update " + TNBook() + " set favorite_count=favorite_count+1 where identify='" + m.Identify+"'"
		_, err = o2.Raw(sql).Exec()
		return err
	}

	return nil
}

func (m *Favorite) IsFavorite() (res int, err error) {
	if "" == m.Identify || m.UserId <= 0 || m.Id > 0 {
		return 0, errors.New("Invalid arguments")
	}
	o := GetOrm("ur")
	if err = o.Read(m, "user_id", "identify"); err != nil {
		return 0, err
	}
	if m.Id > 0 {
		return 1, nil
	} else {
		return 0, nil
	}
}

func (m *Favorite) ListFavoriteByUserId2(page,page_size int) (books []Book, cnt int64, err error) {
	if m.UserId <= 0 {
		err = errors.New("Invalid argument")
		return
	}
	o := GetOrm("ur")
	sql := "select identify from " + TNFavorite() + " where user_id=" + strconv.Itoa(m.UserId) + " limit %v offset %v"
	sql = fmt.Sprintf(sql, page_size, (page-1)*page_size)
	logs.Debug(sql)
	var identifies []string
	cnt,err = o.Raw(sql).QueryRows(&identifies)
	if err != nil || 0 == cnt {
		return
	}
	if cnt < 0 {
		err = errors.New("internal error")
		return
	}
	o = GetOrm("r")
	if _, err = o.QueryTable(TNBook()).Filter("identify__in", identifies).All(&books); err != nil {
		return
	}
	//identifies_str := strings.Join(identifies, ",")
	//sql = "select * from "+TNBook()+" where identify in ("+identifies_str+") limit %v offset %v"
	//sql = fmt.Sprintf(sql, page_size, (page-1)*page_size)
	//_, err = o.Raw(sql).QueryRows(&books)
	return
}

func (m *Favorite) ListFavoriteByUserId(page,page_size int) (books []Book, cnt int64, err error) {
	if m.UserId <= 0 {
		err = errors.New("Invalid argument")
		return
	}
	//TODO: add lend info
	sql := "select b.* from bookms_book b inner join bookms_user_favorite using(identify) where user_id="+strconv.Itoa(m.UserId)+" limit %v offset %v"
	sql = fmt.Sprintf(sql, page_size, (page-1)*page_size)
	logs.Debug(sql)
	o := GetOrm("r")
	cnt,err = o.Raw(sql).QueryRows(&books)
	return
}

func (m *Favorite) ListFavoriteByUserId3(page, page_size int) (books []Book, cnt int64, err error) {
	if m.UserId <= 0 {
		err = errors.New("Invalid argument")
		return
	}
	//TODO: add lend info
	sql := "select * from bookms_book inner join (select identify from bookms_user_favorite where user_id="+strconv.Itoa(m.UserId)+ " limit %v offset %v) as book_fav using(identify)"
	//sql := "select b.* from bookms_user_favorite f left join bookms_book b on f.identify=b.identify where user_id="+strconv.Itoa(m.UserId)+" limit %v offset %v"
	sql = fmt.Sprintf(sql, page_size, (page-1)*page_size)
	logs.Debug(sql)
	o := GetOrm("r")
	cnt,err = o.Raw(sql).QueryRows(&books)
	return
}

func (m *Favorite) ListFavoriteByUserIdReturnUserFav(page,page_size int) (book_favs []UserFavorite, cnt int64, err error) {
	if m.UserId <= 0 {
		err = errors.New("Invalid argument")
		return
	}
	//select f.id,f.user_id,f.identify,b.book_name,b.cover,b.author,c.id,c.category_name,
	//	(case
	//when r.lend_status=0 then '可借'
	//	when r.lend_status=5 then '已下架'
	//	when r.lend_status=1 and r.user_id=1 then '正在借阅'
	//	when r.lend_status=1 and r.user_id<>1 then '不可借'
	//	end) as lend_status
	//	from bookms_book_record r
	//	left join bookms_book b using (identify)
	//	left join bookms_book_category bc using(identify)
	//	left join bookms_user_favorite f using(identify)
	//	left join bookms_category c on bc.category_id=c.id where f.user_id=1 limit 100 offset 20000;
	sql := "select f.id,f.user_id,f.identify,b.book_name,b.cover,b.author,bc.category_id,c.category_name," +
		"(case " +
		"when r.lend_status=0 then '可借' " +
		"when r.lend_status=5 then '已下架' " +
		"when r.lend_status=1 and r.user_id="+strconv.Itoa(m.UserId)+" then '正在借阅' " +
		"when r.lend_status=1 and r.user_id<>"+strconv.Itoa(m.UserId)+" then '不可借' " +
		"end) as lend_status " +
		"from bookms_book_record r " +
		"left join "+TNBook()+" b using(identify) " +
		"left join "+TNBookCategory()+" bc using(identify) " +
		//"left join "+TNFavorite()+" f using(identify) " +
		//using(user_id) result is not target page
		"inner join "+TNFavorite()+" f using(identify) " +
		"left join "+TNCategory()+" c on bc.category_id=c.id " +
		"where f.user_id="+strconv.Itoa(m.UserId)+" limit %v offset %v"
	sql = fmt.Sprintf(sql, page_size, (page-1)*page_size)
	logs.Debug(sql)
	o := GetOrm("r")
	cnt,err = o.Raw(sql).QueryRows(&book_favs)
	return
}

func (m *Favorite) ListFavoriteByUserIdReturnUserFav2(page, page_size int) (book_favs []*UserFavorite, cnt int64, err error) {
	if m.UserId <= 0 {
		err = errors.New("Invalid argument")
		return
	}
	//sql := "select f.id,f.user_id,f.identify,b.book_name,b.cover,b.author,bc.category_id,c.category_name," +
	//	"(case " +
	//	"when r.lend_status=0 then '可借' " +
	//	"when r.lend_status=5 then '已下架' " +
	//	"when r.lend_status=1 and r.user_id="+strconv.Itoa(m.UserId)+" then '正在借阅' " +
	//	"when r.lend_status=1 and r.user_id<>"+strconv.Itoa(m.UserId)+" then '不可借' " +
	//	"end) as lend_status " +
	//	"from"+ TNBookRecord()+ "r "+
	//	"left join "+TNBook()+" b using(identify) " +
	//	"left join "+TNBookCategory()+" bc using(identify) " +
	//	"inner join (select fav.id,fav.user_id,fav.identify from "+TNFavorite()+" fav where fav.user_id=1 limit %v offset %v) f using(user_id) " +
	//	"left join "+TNCategory()+" c on bc.category_id=c.id " +
	//	"where f.user_id="+strconv.Itoa(m.UserId)+" limit %v offset %v"
	sql := "select id,user_id,identify from "+TNFavorite()+" where user_id="+strconv.Itoa(m.UserId)+" limit %v offset %v"
	sql = fmt.Sprintf(sql, page_size, (page-1)*page_size)
	logs.Debug(sql)
	o := GetOrm("ur")
	if cnt,err = o.Raw(sql).QueryRows(&book_favs); err != nil {
		return
	}
	var identifies []string
	for _, bf := range book_favs {
		if "" != bf.Identify {
			identifies = append(identifies, bf.Identify)
		}
	}
	if 0 == len(identifies) {
		return
	}
	identifies_str := "'"
	identifies_str += strings.Join(identifies, "','")
	identifies_str += "'"
	sql = `select ` +
		`(case ` +
		`when r.lend_status=0 then '可借' ` +
		`when r.lend_status=5 then '已下架' ` +
		`when r.lend_status=1 and r.user_id=`+strconv.Itoa(m.UserId)+` then '正在借阅' ` +
		`when r.lend_status=1 and r.user_id<>`+strconv.Itoa(m.UserId)+` then '不可借' ` +
		`end) as lend_status,identify ` +
		`from `+ TNBookRecord()+ ` r where r.identify in(` + identifies_str + `)`
	logs.Debug(sql)
	type Status struct {
		Lend_status string
		Identify string
	}
	var statusObj []Status
	o = GetOrm("r")
	if cnt,err = o.Raw(sql).QueryRows(&statusObj); err != nil {
		return
	}
	statusMap := map[string]string{}
	for _,o := range statusObj {
		statusMap[o.Identify]=o.Lend_status
	}

	sql = `select book_name,cover,author,identify from `+TNBook()+` where identify in(` +identifies_str+ `)`
	logs.Debug(sql)
	type TempBook struct {
		BookName string
		Cover string
		Author string
		Identify string
	}
	var bookObj []TempBook
	if cnt,err = o.Raw(sql).QueryRows(&bookObj); err != nil {
		return
	}
	bookMap := make(map[string]TempBook)
	for _,b := range bookObj {
		bookMap[b.Identify]=b
	}

	sql = `select category_id,identify from `+ TNBookCategory()+ ` where identify in(`+identifies_str+`)`
	logs.Debug(sql)
	type TempBookCategory struct{
		CategoryId int
		Identify string
	}
	var bcategoryObj []TempBookCategory
	if cnt,err = o.Raw(sql).QueryRows(&bcategoryObj); err != nil {
		return
	}
	bcategoryMap := make(map[string]int)
	var cids []string
	for _, c := range bcategoryObj {
		bcategoryMap[c.Identify] = c.CategoryId
		cids = append(cids, strconv.Itoa(c.CategoryId))
	}

	cids_str := strings.Join(cids, ",")
	sql = `select category_name from `+TNCategory()+` where id in(`+cids_str+`)`
	logs.Debug(sql)
	type TempCategory struct {
		CategoryId int
		CategoryName string
	}
	var categoryObj []TempCategory
	if cnt,err = o.Raw(sql).QueryRows(&categoryObj); err != nil {
		return
	}
	categoryMap := make(map[int]TempCategory)
	for _,c := range categoryObj {
		categoryMap[c.CategoryId] = c
	}
	for _, fav := range book_favs {
		fav.LendStatus = statusMap[fav.Identify]
		fav.BookName = bookMap[fav.Identify].BookName
		fav.Author = bookMap[fav.Identify].Author
		fav.Cover = bookMap[fav.Identify].Cover
		fav.CategoryId = bcategoryMap[fav.Identify]
		fav.CategoryName = categoryMap[fav.CategoryId].CategoryName
	}

	return
}

func (m *Favorite) ListFavoriteByUserIdReturnUserFav3(page, page_size int) (book_favs []UserFavorite, cnt int64, err error) {
	if m.UserId <= 0 {
		err = errors.New("Invalid argument")
		return
	}
	//select f.id,f.user_id,f.identify,b.book_name,b.cover,b.author,c.id,c.category_name,
	//	(case
	//when r.lend_status=0 then '可借'
	//	when r.lend_status=5 then '已下架'
	//	when r.lend_status=1 and r.user_id=1 then '正在借阅'
	//	when r.lend_status=1 and r.user_id<>1 then '不可借'
	//	end) as lend_status
	//	from bookms_book_record r
	//	left join bookms_book b using (identify)
	//	left join bookms_book_category bc using(identify)
	//	inner join (select fav.id,fav.user_id,fav.identify from bookms_user_favorite fav where fav.user_id=1 limit 100 offset 20000)f using(user_id)
	//	left join bookms_category c on bc.category_id=c.id
	// where f.user_id=1 limit 100 offset 20000;
	sql := "select f.id,f.user_id,f.identify,b.book_name,b.cover,b.author,bc.category_id,c.category_name," +
		"(case " +
		"when r.lend_status=0 then '可借' " +
		"when r.lend_status=5 then '已下架' " +
		"when r.lend_status=1 and r.user_id="+strconv.Itoa(m.UserId)+" then '正在借阅' " +
		"when r.lend_status=1 and r.user_id<>"+strconv.Itoa(m.UserId)+" then '不可借' " +
		"end) as lend_status " +
		"from bookms_book_record r " +
		"left join "+TNBook()+" b using(identify) " +
		"left join "+TNBookCategory()+" bc using(identify) " +
		"inner join (select fav.id,fav.user_id,fav.identify from "+TNFavorite()+" fav where fav.user_id="+strconv.Itoa(m.UserId)+" limit %v offset %v) f using(user_id) " +
		"left join "+TNCategory()+" c on bc.category_id=c.id " +
		"where f.user_id="+strconv.Itoa(m.UserId)+" limit %v offset %v"
	sql = fmt.Sprintf(sql, page_size, (page-1)*page_size, page_size, (page-1)*page_size)
	logs.Debug(sql)
	o := GetOrm("r")
	cnt,err = o.Raw(sql).QueryRows(&book_favs)
	return
}

func (m *Favorite) ListFavoriteByUserIdReturnUserFavbak(page, page_size int) (book_favs []UserFavorite, cnt int64, err error) {
	if m.UserId <= 0 {
		err = errors.New("Invalid argument")
		return
	}
	sql := "select id,user_id,identify from "+TNFavorite()+" where user_id="+strconv.Itoa(m.UserId)+" limit %v offset %v"
	sql = fmt.Sprintf(sql, page_size, (page-1)*page_size)
	logs.Debug(sql)
	o := GetOrm("ur")
	if cnt,err = o.Raw(sql).QueryRows(&book_favs); err != nil {
		logs.Debug("ListFavoriteByUserIdReturnUserFavGC:",err.Error())
		return
	}
	var identifies []string
	for _, bf := range book_favs {
		if "" != bf.Identify {
			identifies = append(identifies, bf.Identify)
		}
	}
	if 0 == len(identifies) {
		return
	}
	identifies_str := "('"
	identifies_str += strings.Join(identifies, "','")
	identifies_str += "')"
	sql = `select (case when r.lend_status=0 then '可借' when r.lend_status=5 then '已下架' when r.lend_status=1 and r.user_id=`+
		strconv.Itoa(m.UserId)+
		` then '正在借阅' when r.lend_status=1 and r.user_id<>`+
		strconv.Itoa(m.UserId)+
		` then '不可借' end) as lend_status,identify from `+
		TNBookRecord()+ ` r where r.identify in` + identifies_str
	logs.Debug(sql)
	type Status struct {
		Lend_status string
		Identify string
	}
	var statusObj []Status
	o = GetOrm("r")
	if cnt,err = o.Raw(sql).QueryRows(&statusObj); err != nil {
		return
	}
	statusMap := make(map[string]string, cnt)
	for _,o := range statusObj {
		statusMap[o.Identify]=o.Lend_status
	}

	sql = `select book_name,cover,author,identify from `+TNBook()+` where identify in` + identifies_str
	logs.Debug(sql)
	type TempBook struct {
		BookName string
		Cover string
		Author string
		Identify string
	}
	var bookObj []TempBook
	if cnt,err = o.Raw(sql).QueryRows(&bookObj); err != nil {
		return
	}
	bookMap := make(map[string]TempBook, cnt)
	for _,b := range bookObj {
		bookMap[b.Identify]=b
	}

	sql = `select category_id,identify from `+ TNBookCategory()+ ` where identify in` + identifies_str
	logs.Debug(sql)
	type TempBookCategory struct{
		CategoryId int
		Identify string
	}
	var bcategoryObj []TempBookCategory
	if cnt,err = o.Raw(sql).QueryRows(&bcategoryObj); err != nil {
		return
	}
	bcategoryMap := make(map[string]int, cnt)
	var cids []string
	for _, c := range bcategoryObj {
		bcategoryMap[c.Identify] = c.CategoryId
		cids = append(cids, strconv.Itoa(c.CategoryId))
	}

	cids_str := "('"
	cids_str += strings.Join(cids, "','")
	cids_str += "')"
	sql = `select category_name from `+TNCategory()+` where id in`+cids_str
	logs.Debug(sql)
	type TempCategory struct {
		CategoryId int
		CategoryName string
	}
	//categoryObj := make([]TempCategory, 0)
	var categoryObj []TempCategory
	if _,err = o.Raw(sql).QueryRows(&categoryObj); err != nil {
		return
	}
	categoryMap := make(map[int]TempCategory, 20)
	for _,c := range categoryObj {
		categoryMap[c.CategoryId] = c
	}
	for i:=0;i<len(book_favs);i++{
		book_favs[i].LendStatus = statusMap[book_favs[i].Identify]
		book_favs[i].BookName = bookMap[book_favs[i].Identify].BookName
		book_favs[i].Author = bookMap[book_favs[i].Identify].Author
		book_favs[i].Cover = bookMap[book_favs[i].Identify].Cover
		book_favs[i].CategoryId = bcategoryMap[book_favs[i].Identify]
		book_favs[i].CategoryName = categoryMap[book_favs[i].CategoryId].CategoryName
	}

	return
}