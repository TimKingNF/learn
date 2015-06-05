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
			"TeacherInfo",
			`/teacherInfo`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"ClassInfo",
			`/classInfo`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"CourseInfo",
			`/courseInfo`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"CourseHistory",
			`/courseHistory`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"CourseHomework",
			`/courseHomework`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"StudentCheck",
			`/studentCheck`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"StudentHomework",
			`/studentHomework`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"Download",
			`/download`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"StudentScore",
			`/studentScore`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewStudentController"],
		beego.ControllerComments{
			"Echart",
			`/echart`,
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
			"Info",
			`/info`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"],
		beego.ControllerComments{
			"ClassInfo",
			`/classInfo`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"],
		beego.ControllerComments{
			"CourseHistory",
			`/courseHistory`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"],
		beego.ControllerComments{
			"Check",
			`/check`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"],
		beego.ControllerComments{
			"Profile",
			`/profile`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"],
		beego.ControllerComments{
			"Setting",
			`/setting`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"],
		beego.ControllerComments{
			"CourseInfo",
			`/courseInfo`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"],
		beego.ControllerComments{
			"CourseHomework",
			`/courseHomework`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"],
		beego.ControllerComments{
			"Notice",
			`/notice`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"],
		beego.ControllerComments{
			"UploadImg",
			`/uploadImg`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"],
		beego.ControllerComments{
			"Table",
			`/table`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"] = append(beego.GlobalControllerRouter["learn/controllers:ViewTeatherController"],
		beego.ControllerComments{
			"Download",
			`/download`,
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
			"UploadAttachment",
			`/uploadAttachment`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"ReadNoticeByStudent",
			`/readNoticeByStudent`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"DeleteStudentNotice",
			`/deleteStudentNotice`,
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
			"StudentCourse",
			`/studentCourse`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"StudentHomework",
			`/studentHomework`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"StudentCheck",
			`/studentCheck`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"CourseHistory",
			`/courseHistory`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"StudentScore",
			`/studentScore`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"CourseHomework",
			`/courseHomework`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"GetTermNumberByStudent",
			`/getTermNumberByStudent`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APIStudentController"] = append(beego.GlobalControllerRouter["learn/controllers:APIStudentController"],
		beego.ControllerComments{
			"GetStudentCourseByTerm",
			`/getStudentCourseByTerm`,
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

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"UpdatePassword",
			`/updatePassword`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"GetStudentHomeworksByTeacherCourseHomework",
			`/getStudentHomeworksByTeacherCourseHomework`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"SetGradeByStudentHomework",
			`/setGradeByStudentHomework`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"TeacherCourse",
			`/teacherCourse`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"CourseHomework",
			`/courseHomework`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"CourseHistory",
			`/courseHistory`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"StudentCheck",
			`/studentCheck`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"UpdateImg",
			`/updateImg`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"UploadImg",
			`/uploadImg`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"PublisHomework",
			`/publisHomework`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"UpdateHomework",
			`/updateHomework`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"UploadAttachment",
			`/uploadAttachment`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"UpdateProfile",
			`/updateProfile`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"DelTeacherCourseHomework",
			`/delTeacherCourseHomework`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"GetTeacherCourseHomework",
			`/getTeacherCourseHomework`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"GetTermNumberByTeacher",
			`/getTermNumberByTeacher`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["learn/controllers:APITeacherController"] = append(beego.GlobalControllerRouter["learn/controllers:APITeacherController"],
		beego.ControllerComments{
			"GetTeacherCourseByTerm",
			`/getTeacherCourseByTerm`,
			[]string{"post"},
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

}
