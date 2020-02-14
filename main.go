package main

import (
	"bookms/cache"
	_ "bookms/init"
	_ "bookms/models"
	_ "bookms/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	logs.SetLevel(logs.LevelDebug)
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
	defer cache.ClosePool()
	beego.Run()
}

