package routers

import (
	"bookms/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{}, "*:Get")

    //user
    beego.Router("/login", &controllers.LoginController{}, "*:Login")
    beego.Router("/logout",&controllers.LoginController{}, "*:Logout")
    beego.Router("/register", &controllers.LoginController{}, "Get:Register;Post:RegisterDo")
}
