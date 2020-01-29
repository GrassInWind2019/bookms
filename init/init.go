package init

import (
	"encoding/gob"
	"bookms/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/sony/sonyflake"
	"strconv"
	"time"
)

var (
	sonyflakeObj *sonyflake.Sonyflake
)
func init() {
	initUUID()
	gob.Register(models.User{})
	orm.RegisterDriver("mysql", orm.DRMySQL)
	dbinit("w","r")
	dbinit("uw","ur")
}

func initUUID() {
	machine_idstr := beego.AppConfig.String("machine_id")
	if "" == machine_idstr {
		panic("UUID machine id cannot be null")
	}
	machine_id, err := strconv.ParseInt(machine_idstr, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	st := sonyflake.Settings{
		StartTime:time.Now(),
		MachineID: func() (uint16, error) {
			if machine_id > 2^16 {
				return 0, errors.New("Machine id overflow!")
			}
			return uint16(machine_id),nil
		},
		CheckMachineID: func(u uint16) bool {
			return  true
		},
	}
	sonyflakeObj = sonyflake.NewSonyflake(st)
	if sonyflakeObj == nil {
		panic("Create sonyflake object failed!")
	}
	logs.Debug("sonyflake object create success", sonyflakeObj)
}

func GetSonyFlakeObj() *sonyflake.Sonyflake {
	return sonyflakeObj
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