package main

import (
	_ "bookms/init"
	_ "bookms/models"
	_ "bookms/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"runtime"
	"time"
)

func main() {
	logs.SetLevel(logs.LevelDebug)
	//models.UpdateCategoryCount()
	//go RunGC()
	beego.Run()
}

func RunGC() {
	logs.SetLevel(logs.LevelInfo)
	t := time.NewTicker(60*time.Second)
	for {
		select {
			case <-t.C:
				runtime.GC()
				logs.Debug("GC runned!")
		}
	}
}

