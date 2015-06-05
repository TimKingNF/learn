package routers

import (
	"github.com/astaxie/beego"
	"learn/controllers"
)

func init() {
	//	注册路由
	beego.Router("/", &controllers.ViewController{}, "Get:Index;Post:Login")
	beego.Router("/signOut", &controllers.ViewController{}, "Get:SignOut")
	beego.Router("/error", &controllers.ViewController{}, "Get:Error")

	api := beego.NewNamespace("/api",
		beego.NSNamespace("/student",
			beego.NSInclude(
				&controllers.APIStudentController{},
			),
		),
		beego.NSNamespace("/admin",
			beego.NSInclude(
				&controllers.APIAdminController{},
			),
		),
		beego.NSNamespace("/teacher",
			beego.NSInclude(
				&controllers.APITeacherController{},
			),
		),
	)

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
	beego.AddNamespace(api)
	beego.AddNamespace(view)
	AddRouterFilter()
}
