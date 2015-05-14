package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	
	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"Table",
			`/table`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"],
		beego.ControllerComments{
			"File",
			`/file`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewAdminController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewAdminController"],
		beego.ControllerComments{
			"Index",
			`/index`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewAdminController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewAdminController"],
		beego.ControllerComments{
			"Setting",
			`/setting`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewAdminController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewAdminController"],
		beego.ControllerComments{
			"Student",
			`/student`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewAdminController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewAdminController"],
		beego.ControllerComments{
			"Backup",
			`/backup`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewAdminController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewAdminController"],
		beego.ControllerComments{
			"Teacher",
			`/teacher`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewAdminController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewAdminController"],
		beego.ControllerComments{
			"Course",
			`/course`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewController"],
		beego.ControllerComments{
			"Index",
			`/index`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewController"],
		beego.ControllerComments{
			"Login",
			`/index`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewController"],
		beego.ControllerComments{
			"SignOut",
			`/signOut`,
			[]string{"get"},
			nil})

}
