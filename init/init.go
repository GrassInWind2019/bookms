package init

import (
	"encoding/gob"
	"bookms/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)
func init() {
	gob.Register(models.User{})
	orm.RegisterDriver("mysql", orm.DRMySQL)
	dbinit("w","r")
	dbinit("uw","ur")
}

func dbinit(aliases ...string) {
	isDev := (beego.AppConfig.String("runmode") == "dev")
	if len(aliases) > 0 {
		for _, alias := range aliases {
			registDatabase(alias)
			if "w" == alias {
				orm.RunSyncdb("default", false, isDev)
			}
		}
	} else {
		registDatabase("w")
		orm.RunSyncdb("default", false, isDev)
	}
	orm.Debug = isDev
}

func registDatabase(alias string) {
	if len(alias) == 0 {
		return
	}
	//连接名称
	dbAlias := alias
	if "w" == alias || "default" == alias {
		dbAlias = "default"
		alias = "w"
	}
	dbName := beego.AppConfig.String("db_"+alias+"_database")
	dbUser := beego.AppConfig.String("db_"+alias+"_username")
	dbPwd := beego.AppConfig.String("db_"+alias+"_password")
	dbHost := beego.AppConfig.String("db_"+alias+"_host")
	dbPort := beego.AppConfig.String("db_"+alias+"_port")
	logs.Debug("Register DB " + dbName+ " for "+dbAlias)
	err := orm.RegisterDataBase(dbAlias, "mysql",dbUser+":"+dbPwd+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8",30)
	if err != nil {
		logs.Error(err.Error())
	}
}