package routers

import (
	"bookms/controllers"
	"github.com/astaxie/beego"
)

func init() {
    //beego.Router("/", &controllers.MainController{}, "*:Get")
    //home
    beego.Router("/", &controllers.HomeController{}, "*:Index")

    //user
    beego.Router("/login", &controllers.LoginController{}, "*:Login")
    beego.Router("/logout",&controllers.LoginController{}, "*:Logout")
    beego.Router("/register", &controllers.LoginController{}, "get:Register;post:RegisterDo")
    beego.Router("/favoritedo/:identify", &controllers.UserController{}, "*:FavoriteDo")
    beego.Router("/bookcomment/:identify", &controllers.UserController{}, "*:AddCommentAndScore")
    beego.Router("/usercenter/:page", &controllers.UserController{}, "*:GetUserCenterFav2")
    //beego.Router("/usercenter/:page", &controllers.UserController{}, "*:GetUserCenterInfo")
    beego.Router("/usercenter2/:page", &controllers.UserController{}, "*:GetUserCenterInfo2")
    beego.Router("/usercenter3/:page", &controllers.UserController{}, "*:GetUserCenterInfo3")
    beego.Router("/usercenterfav/:page", &controllers.UserController{}, "*:GetUserCenterFav")
    beego.Router("/usercenterfav2/:page", &controllers.UserController{}, "*:GetUserCenterFav2")
    beego.Router("/usercenterfav3/:page", &controllers.UserController{}, "*:GetUserCenterFav3")

    //category
    beego.Router("/addcategory", &controllers.CategoryController{}, "*:AddCategory")

    //book
    beego.Router("/addbook", &controllers.BookController{}, "*:AddBook")
    beego.Router("/updatebook/:identify", &controllers.BookController{}, "*:UpdateBookByIdentify")
    beego.Router("/deletebook/:identify", &controllers.BookController{}, "*:DeleteBooksByIdentify")
    beego.Router("/deletebook/:book_id", &controllers.BookController{}, "*:DeleteBookById")
    beego.Router("/lendbook/:book_id", &controllers.BookController{}, "*:LendBookById")
    beego.Router("/returnbook/:book_id", &controllers.BookController{}, "*:ReturnBookById")
    beego.Router("/bookdetail/:identify", &controllers.GetBookController{}, "get:GetBooksByIdentify")
}
