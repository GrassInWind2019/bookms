package main

import (
	_ "bookms/routers"
	_ "bookms/init"
	_ "bookms/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	logs.SetLevel(logs.LevelDebug)
	beego.Run()
}

