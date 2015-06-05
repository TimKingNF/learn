package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	
	beego.GlobalControllerRouter["learn/controllers:ViewController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewController"],
		beego.ControllerComments{
			"Index",
			`/index`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewController"],
		beego.ControllerComments{
			"Error",
			`/error`,
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

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"Table",
			`/table`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"Info",
			`/info`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"Setting",
			`/setting`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"UploadImg",
			`/uploadImg`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"Notice",
			`/notice`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"Login",
			`/login`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"EduLoading",
			`/eduLoading`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"EduManage",
			`/eduManage`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"],
		beego.ControllerComments{
			"File",
			`/file`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIAdminController"] = append(beego.GlobalControllerRouter["learn/controllers:APIAdminController"],
		beego.ControllerComments{
			"UpdatePassword",
			`/updatePassword`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"UpdatePassword",
			`/updatePassword`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"UpdateImg",
			`/updateImg`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"UploadImg",
			`/uploadImg`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"Login",
			`/login`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"EduLoading",
			`/eduLoading`,
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

}
