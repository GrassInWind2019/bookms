package main

import (
	"bookms/cache"
	_ "bookms/init"
	_ "bookms/models"
	_ "bookms/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	_ "net/http/pprof"
)

func main() {
	logs.SetLevel(logs.LevelWarn)
	runtime.GOMAXPROCS(2)
	cpu_file, err := os.OpenFile("cpu.profile", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	pprof.StartCPUProfile(cpu_file)
	defer pprof.StopCPUProfile()
	mem_file, err := os.OpenFile("mem.profile", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	pprof.WriteHeapProfile(mem_file)
	//models.UpdateCategoryCount()
	//go RunGC()
	defer cache.ClosePool()
	beego.Run()
}

func RunGC() {
	logs.SetLevel(logs.LevelInfo)
	t := time.NewTicker(5*time.Second)
	for {
		select {
			case <-t.C:
				runtime.GC()
				logs.Debug("GC runned!")
		}
	}
}

