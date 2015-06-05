package controllers

import (
	"learn/models"
	"strconv"
	"strings"
)

type ViewTeatherController struct {
	ViewController
}

func (this *ViewTeatherController) Prepare() {
	//	get session
	user_type := this.GetSession("type").(string)
	if user_type == "教师" {
		id := this.GetSession("id").(string)
		if len(id) > 0 {
			if t_id, err := strconv.ParseInt(id, 10, 64); err == nil {
				teacher, err := models.GetTeacherById(t_id)
				if err == nil {
					this.Data["teacher"] = teacher
					//	设置操作签名，获取签名参数
					appid, sessid := this.SetSignature()
					this.Data["appid"] = appid
					this.Data["sessid"] = sessid
					this.Data["key"] = models.Str2Sha1(this.Ctx.Input.Cookie("beegosessionID"))
					return
				}
			}
		}
	}
	this.Redirect("/error", 302)
	this.StopRun()
}

// @Title 个人信息视图
// @router /info [get]
func (this *ViewTeatherController) Info() {
	//	get teacher profile
	t_profile, _ := models.GetTeacherProfile(this.Data["teacher"].(*models.Teacher).Id)
	this.Data["teacher"].(*models.Teacher).Profile = t_profile

	this.Layout = "teacher/base.html"
	this.TplNames = "teacher/info.html"
}

// @Tilte 查看班级名单
// @router /classInfo [get]
func (this *ViewTeatherController) ClassInfo() {
	if class_id, err := this.GetInt64("class"); err == nil {
		if class, err := models.GetClassById(class_id); err == nil {
			class.Students, _ = models.GetStudentsByClass(class)
			this.Data["class"] = class
		} else {
			this.Redirect("/error", 302)
			this.StopRun()
		}
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}

	this.Layout = "teacher/base.html"
	this.TplNames = "teacher/classInfo.html"
}

// @Title 查看课程点名记录
// @router /courseHistory [get]
func (this *ViewTeatherController) CourseHistory() {
	//	get week
	if week, _ := this.GetInt("week"); week > 0 {
		this.Data["now_week"] = week
	}

	if t_course_id, err := this.GetInt64("course"); err == nil {
		if t_course, err := models.GetTeacherCourseById(t_course_id); err == nil {
			//	get teacher course class and student list
			t_course.Classes, _ = models.GetClassesByTeacherCourse(t_course)
			students, _ := models.GetStudentsByTeacherCourse(t_course)
			week := 1
			if this.Data["now_week"] != nil {
				week = this.Data["now_week"].(int)
			}
			t_course.StudentChecks, _ = models.GetStudentChecksByTeacherCourseAndWeek(t_course, week)
			t_course.Orgs["students"] = students
			t_course.Orgs["s_checks"], _ = models.GetStudentChecksByTeacherCourse(t_course)
			this.Data["t_course"] = t_course
			if t_courses, err := models.GetTeacherCourseByTerm(t_course.Term, this.Data["teacher"].(*models.Teacher).Id); err == nil {
				this.Data["t_courses"] = t_courses
				this.Data["course_week"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"}
			} else {
				this.Redirect("/error", 302)
				this.StopRun()
			}
		} else {
			this.Redirect("/error", 302)
			this.StopRun()
		}
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	//	get all term
	term_list, err := models.GetTermListByTeacherCourse(this.Data["teacher"].(*models.Teacher).Id)
	if err != nil {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	this.Data["term_list"] = models.RankingTerm(term_list)

	this.Layout = "teacher/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "teacher/head/courseHistory_head.html"
	this.LayoutSections["Scripts"] = "teacher/scripts/courseHistory_scripts.html"
	this.TplNames = "teacher/courseHistory.html"
}

// @Title 点名
// @router /check [get]
func (this *ViewTeatherController) Check() {
	if t_course_id, err := this.GetInt64("course"); err == nil {
		if t_course, err := models.GetTeacherCourseById(t_course_id); err == nil {
			//	get teacher course class and student list
			t_course.Classes, _ = models.GetClassesByTeacherCourse(t_course)
			students, _ := models.GetStudentsByTeacherCourse(t_course)
			t_course.Orgs["students"] = students
			this.Data["t_course"] = t_course
		} else {
			this.Redirect("/error", 302)
			this.StopRun()
		}
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	//	get week
	if week, _ := this.GetInt("week"); week > 0 {
		this.Data["now_week"] = week
	}

	this.Layout = "teacher/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "teacher/scripts/check_scripts.html"
	this.TplNames = "teacher/check.html"
}

// @Title 修改个人资料视图
// @router /profile [get]
func (this *ViewTeatherController) Profile() {
	if exist := models.TeacherProfileExist(this.Data["teacher"].(*models.Teacher).Id); !exist {
		if err := models.AddTeacherProfile(&models.TeacherProfile{Teacher: this.Data["teacher"].(*models.Teacher)}); err != nil {
			this.Redirect("/error", 302)
			this.StopRun()
		}
	}
	//	get teacher profile
	if t_profile, err := models.GetTeacherProfile(this.Data["teacher"].(*models.Teacher).Id); err == nil {
		this.Data["teacher"].(*models.Teacher).Profile = t_profile
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	this.Layout = "teacher/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "scripts/signature_scripts.html"
	this.TplNames = "teacher/profile.html"
}

// @Title 设置
// @router /setting [get]
func (this *ViewTeatherController) Setting() {
	this.Layout = "teacher/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "scripts/signature_scripts.html"
	this.TplNames = "teacher/setting.html"
}

// @Tilte 查看课程信息
// @router /courseInfo [get]
func (this *ViewTeatherController) CourseInfo() {
	//	get teacher course
	if t_course_id, err := this.GetInt64("course"); err == nil {
		if t_course, err := models.GetTeacherCourseById(t_course_id); err == nil {
			//	get teacher course class and student list
			t_course.Classes, _ = models.GetClassesByTeacherCourse(t_course)
			students, _ := models.GetStudentsByTeacherCourse(t_course)
			t_course.Orgs["students"] = students
			t_course.Orgs["s_checks"], _ = models.GetStudentChecksByTeacherCourse(t_course)
			if homeworks, err := models.GetTeacherCourseHomeworkByTeacherCourse(t_course); err == nil {
				for k, _ := range homeworks {
					homeworks[k].StudentHomeworks, _ = models.GetUploadedStudentHomeworkByTeacherCourseHomework(homeworks[k])
				}
				t_course.Homeworks = homeworks
				this.Data["t_course"] = t_course
			}
		} else {
			this.Redirect("/error", 302)
			this.StopRun()
		}
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}

	this.Layout = "teacher/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "teacher/head/courseInfo_head.html"
	this.LayoutSections["Scripts"] = "teacher/scripts/courseInfo_scripts.html"
	this.TplNames = "teacher/courseInfo.html"
}

// @Title 课程作业
// @router /courseHomework [get]
func (this *ViewTeatherController) CourseHomework() {
	//	get teacher course
	if t_course_id, err := this.GetInt64("course"); err == nil {
		if t_course, err := models.GetTeacherCourseById(t_course_id); err == nil {
			//	get teacher course homework list
			if homeworks, err := models.GetTeacherCourseHomeworkByTeacherCourse(t_course); err == nil {
				//	get student homework  by teacher course homework
				for k, _ := range homeworks {
					homeworks[k].StudentHomeworks, _ = models.GetUploadedStudentHomeworkByTeacherCourseHomework(homeworks[k])
				}
				t_course.Homeworks = homeworks
				this.Data["t_course"] = t_course
				//	get the term all teachercourse
				if t_courses, err := models.GetTeacherCourseByTerm(t_course.Term, this.Data["teacher"].(*models.Teacher).Id); err == nil {
					this.Data["t_courses"] = t_courses
				} else {
					this.Redirect("/error", 302)
					this.StopRun()
				}
			} else {
				this.Redirect("/error", 302)
				this.StopRun()
			}
		} else {
			this.Redirect("/error", 302)
			this.StopRun()
		}
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}

	//	get all term
	term_list, err := models.GetTermListByTeacherCourse(this.Data["teacher"].(*models.Teacher).Id)
	if err != nil {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	this.Data["term_list"] = models.RankingTerm(term_list)

	//	get week
	if week, _ := this.GetInt("week"); week > 0 {
		this.Data["now_week"] = week
	}

	this.Layout = "teacher/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "teacher/head/courseHomework_head.html"
	this.LayoutSections["Scripts"] = "teacher/scripts/courseHomework_scripts.html"
	this.TplNames = "teacher/courseHomework.html"
}

// @Tilte 消息
// @router /notice [get]
func (this *ViewTeatherController) Notice() {
	this.Layout = "teacher/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "teacher/scripts/notice_scripts.html"
	this.LayoutSections["Head_html"] = "teacher/head/notice_head.html"
	this.TplNames = "teacher/notice.html"
}

// @Title 修改个人头像
// @router /uploadImg [get]
func (this *ViewTeatherController) UploadImg() {
	this.Layout = "teacher/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "teacher/scripts/uploadimg_scripts.html"
	this.TplNames = "teacher/uploadImg.html"
}

// @Title 课程表视图
// @router /table [get]
func (this *ViewTeatherController) Table() {
	change, _ := this.GetBool("change")
	this.Data["change"] = change
	//	get the term list
	term_list, err := models.GetTermListByTeacherCourse(this.Data["teacher"].(*models.Teacher).Id)
	if err != nil {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	this.Data["term_list"] = models.RankingTerm(term_list)
	var term *models.Term
	xnd := this.GetString("xnd")
	if len(xnd) > 0 {
		xqd, _ := this.GetInt("xqd")
		arr := strings.Split(xnd, "-")
		if len(arr) != 2 {
			this.Redirect("/error", 302)
			this.StopRun()
		}
		startYear, _ := strconv.Atoi(arr[0])
		endYear, _ := strconv.Atoi(arr[1])
		term, err = models.SearchTerm(xqd, startYear, endYear)
		if err != nil {
			this.Redirect("/error", 302)
			this.StopRun()
		}
	} else {
		//	get the now term
		term, err = models.GetTermByTimeNow()
		if err != nil {
			this.Redirect("/error", 302)
			this.StopRun()
		}
	}
	this.Data["term"] = term
	// get the teacher course
	t_courses, err := models.GetTeacherCourseByTerm(term, this.Data["teacher"].(*models.Teacher).Id)
	if err != nil {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	this.Data["course_color"] = []string{"tb-red", "tb-dark-blue", "tb-grey", "tb-violet", "tb-origin", "tb-blue-green", "tb-blue", "tb-yellow",
		"tb-green", "tb-light-blue", "tb-brown", "tb-pink", "tb-light-green", "tb-dark-blue", "tb-rose-red"}
	this.Data["course_time"] = []string{"08:00-08:40", "08:50-09:30", "09:45-10:25", "10:35-11:15", "11:20-12:00",
		"12:50-13:30", "13:40-14:20", "14:30-15:10", "15:15-15:55", "16:10-16:50", "16:55-17:35", "18:45-19:25",
		"19:30-20:10", "20:15-20:55", "21:05-21:45"}
	this.Data["course_week"] = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"}
	this.Data["t_courses"] = models.RangeToMapTeacherCourse(t_courses)
	this.Data["now_week"] = this.GetString("week")

	this.Layout = "teacher/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "teacher/head/table_head.html"
	this.LayoutSections["Scripts"] = "teacher/scripts/table_scripts.html"
	this.TplNames = "teacher/table.html"
}

// @Title 下载附件
// @router /download [get]
func (this *ViewTeatherController) Download() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "教师" {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	filepath := this.GetString("filepath")
	if exist := models.PathIsExist(filepath); exist {
		this.Ctx.Output.Download(string(filepath[1:]))
		this.StopRun()
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}
}
