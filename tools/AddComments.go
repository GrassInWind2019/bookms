package main

import (
	_ "bookms/init"
	"bookms/models"
	"bookms/utils"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
)

func AddMultiComments(score, user_id, comments_cnt int, identify string) {
	comment := models.Comments{
		Identify: identify,
		UserId:user_id,
		Content:"test publish comment, just for fun! ",
	}
	scoreObj := models.Score{
		Score:score,
		UserId:user_id,
		Identify:identify,
		CreateTime:time.Now(),
	}
	err := scoreObj.AddScore()
	if err != nil {
		logs.Error(err.Error())
	}
	var content string
	for i:=0;i<comments_cnt;i++ {
		comment.Id = 0 //AddComments.go:44]  Error 1062: Duplicate entry '1' for key 'PRIMARY'
		comment.CreateTime = time.Now()
		utils.DeepCopy(&content, &comment.Content)
		comment.Content = comment.Content+strconv.Itoa(i)
		err = comment.AddComment()
		if err != nil {
			logs.Error(err.Error())
			return
		}
		utils.DeepCopy(&comment.Content, &content)
	}
}

func main() {
	AddMultiComments(5, 5, 100, "2")
}
