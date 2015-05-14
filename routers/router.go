package routers

import (
	"github.com/astaxie/beego"
	"learn/controllers"
)

func init() {
	//	注册路由
	beego.Router("/", &controllers.ViewController{}, "Get:Index;Post:Login")
	beego.Router("/signOut", &controllers.ViewController{}, "Get:SignOut")

	view := beego.NewNamespace("/view",
		beego.NSNamespace("/student",
			beego.NSInclude(
				&controllers.ViewStudentController{},
			),
		),
		beego.NSNamespace("/admin",
			beego.NSInclude(
				&controllers.ViewAdminController{},
			),
		),
		beego.NSNamespace("/teacher",
			beego.NSInclude(
				&controllers.ViewTeatherController{},
			),
		),
	)
	// beego.AddNamespace(ns)
	beego.AddNamespace(view)
	AddRouterFilter()
}
