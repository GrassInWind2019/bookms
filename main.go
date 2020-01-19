package main

import (
	_ "bookms/init"
	_ "bookms/models"
	_ "bookms/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	logs.SetLevel(logs.LevelDebug)
	//models.UpdateCategoryCount()
	beego.Run()
}

