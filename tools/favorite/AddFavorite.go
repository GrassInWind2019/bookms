package main

import (
	"bookms/models"
	"database/sql"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"
)

var (
	db *sql.DB
	cookie = "login=P_-DAwEBDkNvb2tpZVJlbWVtYmVyAf-EAAEDAQhNZW1iZXJJZAEEAAEHQWNjb3VudAEMAAEEVGltZQH_hgAAABD_hQUBAQRUaW1lAf-GAAAAHf-EAQIBBWFkbWluAQ8BAAAADtWk7EQjNfwUAeAA|1578300740590740500|61c14cc80ff50249b47557a63c97e67cce81ca45; beegosessionID=4af65f6c2f612aacce92a261afc6ec0f; user=_7D_gQMBAQRVc2VyAf-CAAENAQJJZAEEAAEHQWNjb3VudAEMAAEITmlja25hbWUBDAABCFBhc3N3b3JkAQwAAQVQaG9uZQEMAAEFRW1haWwBDAABBFJvbGUBBAABCFJvbGVOYW1lAQwAAQZBdmF0YXIBDAABBlN0YXR1cwEEAAEKQ3JlYXRlVGltZQH_hAABDUxhc3RMb2dpblRpbWUB_4QAAQlCaW9ncmFwaHkBDAAAABD_gwUBAQRUaW1lAf-EAAAA__v_ggECAQtHcmFzc0luV2luZAELR3Jhc3NJbldpbmQBeDM2MzIzNjY2MzY2NjM2NjIzNjY0MzczMzMwMzMzMzYzNjIzNjMxMzQzMjM2NjQ2MzMyNjMzMjM4NjE2NjY1MzQzMDMwMzEzNzM3Mzg2NDM4MzUzMTYxNjUyMDJjYjk2MmFjNTkwNzViOTY0YjA3MTUyZDIzNGI3MAELMTMwMTIzNDU2NzgBE0dyYXNzSW5XaW5kQHNpbmEuY24BBAEM5pmu6YCa55So5oi3Aw8BAAAADtXDIksAAAAAAeABDwEAAAAO1cRKAhwA1kQB4AEM5pmu6YCa55So5oi3AA==|1580356354477819000|8b8826fdc26c4292bc34a1b75280f2a4e93d7efb"

	)

func init() {
	var err error
	db, err = sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/bookms")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
}

func AddFavorites() {
	client := &http.Client{}
	sql := "select identify from "+models.TNBook()
	rows, err := db.Query(sql)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	var identify string
	i := 0
	t := time.Now()
	for rows.Next() {
		i++
		if i % 2000 == 0 {
			logs.Info("cost: ", time.Since(t))
		}
		if err = rows.Scan(&identify); err != nil {
			panic(err.Error())
		}
		url := "http://localhost:8080/favoritedo/"+identify
		req,_ := http.NewRequest("GET",url,nil)
		req.Header.Add("cookie",cookie)
		resp,_ := client.Do(req)
		//body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Printf(string(body))
		resp.Body.Close()
	}
}

func main() {
	defer db.Close()
	AddFavorites()
}
