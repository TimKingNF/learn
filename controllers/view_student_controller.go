package controllers

import (
	"encoding/json"
	Edu "learn/EDU"
	"learn/models"
	"strconv"
	"strings"
)

type ViewStudentController struct {
	ViewController
}

func (this *ViewStudentController) Prepare() {
	//	get session
	user_type := this.GetSession("type").(string)
	if user_type == "学生" {
		id := this.GetSession("id").(string)
		if len(id) > 0 {
			student, err := models.GetStudentById(id)
			if err == nil {
				this.Data["student"] = student
				//	获取未阅读的通知
				this.Data["noticeNum"] = models.CountNotReadStudentNotice(student.Id)
				//	设置操作签名，获取签名参数
				appid, sessid := this.SetSignature()
				this.Data["appid"] = appid
				this.Data["sessid"] = sessid
				this.Data["key"] = models.Str2Sha1(this.Ctx.Input.Cookie("beegosessionID"))
				return
			}
		}
	}
	this.Redirect("/error", 302)
	this.StopRun()
}

// @Title 课程表视图
// @router /table [get]
func (this *ViewStudentController) Table() {
	change, _ := this.GetBool("change")
	this.Data["change"] = change
	//	get the term list
	term_list, err := models.GetTermListByStudentCourse(this.Data["student"].(*models.Student).Id)
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
	//	get the student course
	s_courses, err := models.GetStudentCourseByTerm(term, this.Data["student"].(*models.Student).Id)
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
	this.Data["s_courses"] = models.RangeToMapStudentCourse(s_courses)
	this.Data["now_week"] = this.GetString("week")

	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "student/head/table_head.html"
	this.LayoutSections["Scripts"] = "student/scripts/table_scripts.html"
	this.TplNames = "student/table.html"
}

// @Title 个人信息视图
// @router /info [get]
func (this *ViewStudentController) Info() {
	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "student/scripts/info_scripts.html"
	this.LayoutSections["Head_html"] = "student/head/info_head.html"
	this.TplNames = "student/info.html"
}

// @Title 设置
// @router /setting [get]
func (this *ViewStudentController) Setting() {
	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "scripts/signature_scripts.html"
	this.TplNames = "student/setting.html"
}

// @Title 修改个人头像
// @router /uploadImg [get]
func (this *ViewStudentController) UploadImg() {
	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "student/scripts/uploadimg_scripts.html"
	this.TplNames = "student/uploadImg.html"
}

// @Tilte 消息
// @router /notice [get]
func (this *ViewStudentController) Notice() {
	//	get student notice
	if s_notices, err := models.GetStudentNoticeByStudent(this.Data["student"].(*models.Student).Id); err == nil {
		tmp_s_notices, _ := models.RankingNotice(s_notices)
		this.Data["s_notices"] = tmp_s_notices
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}

	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "student/scripts/notice_scripts.html"
	this.LayoutSections["Head_html"] = "student/head/notice_head.html"
	this.TplNames = "student/notice.html"
}

// @Tilte 查看教师资料
// @router /teacherInfo [get]
func (this *ViewStudentController) TeacherInfo() {
	if teacher_id, err := this.GetInt64("teacher"); err == nil {
		if teacher, err := models.GetTeacherById(teacher_id); err == nil {
			teacher.Profile, _ = models.GetTeacherProfile(teacher_id)
			this.Data["teacher"] = teacher
		} else {
			this.Redirect("/error", 302)
			this.StopRun()
		}
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}

	this.Layout = "student/base.html"
	this.TplNames = "student/teacherInfo.html"
}

// @Tilte 查看班级名单
// @router /classInfo [get]
func (this *ViewStudentController) ClassInfo() {
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

	this.Layout = "student/base.html"
	this.TplNames = "student/classInfo.html"
}

// @Tilte 查看课程信息
// @router /courseInfo [get]
func (this *ViewStudentController) CourseInfo() {
	//	get teacher course
	if t_course_id, err := this.GetInt64("course"); err == nil {
		if t_course, err := models.GetTeacherCourseById(t_course_id); err == nil {
			//	get teacher course class and student list
			t_course.Classes, _ = models.GetClassesByTeacherCourse(t_course)
			students, _ := models.GetStudentsByTeacherCourse(t_course)
			t_course.Orgs["students"] = students
			t_course.Orgs["s_checks"], _ = models.GetStudentChecksByTeacherCourseAndStudent(t_course, this.Data["student"].(*models.Student).Id)
			if s_homeworks, err := models.GetStudentHomeworkByTeacherCourseAndStudent(t_course, this.Data["student"].(*models.Student).Id); err == nil {
				t_course.Orgs["s_homeworks"] = s_homeworks
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

	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "student/head/courseInfo_head.html"
	this.LayoutSections["Scripts"] = "student/scripts/courseInfo_scripts.html"
	this.TplNames = "student/courseInfo.html"
}

// @Title 查看课程点名记录
// @router /courseHistory [get]
func (this *ViewStudentController) CourseHistory() {
	//	get teacher course
	if t_course_id, err := this.GetInt64("course"); err == nil {
		if t_course, err := models.GetTeacherCourseById(t_course_id); err == nil {
			//	get student check list
			if s_checks, err := models.GetStudentChecksByTeacherCourseAndStudent(t_course, this.Data["student"].(*models.Student).Id); err == nil {
				t_course.Orgs["s_checks"] = s_checks
				this.Data["t_course"] = t_course
				//	get the term all teachercourse
				if s_courses, err := models.GetStudentCourseByTerm(t_course.Term, this.Data["student"].(*models.Student).Id); err == nil {
					this.Data["s_courses"] = s_courses
				} else {
					this.Redirect("/error", 302)
					this.StopRun()
				}
				//	get all term
				term_list, err := models.GetTermListByStudentHomework(this.Data["student"].(*models.Student).Id)
				if err != nil {
					this.Redirect("/error", 302)
					this.StopRun()
				}
				this.Data["term_list"] = models.RankingTerm(term_list)
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

	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "student/head/courseHistory_head.html"
	this.LayoutSections["Scripts"] = "student/scripts/courseHistory_scripts.html"
	this.TplNames = "student/courseHistory.html"
}

// @Title 查看课程作业
// @router /courseHomework [get]
func (this *ViewStudentController) CourseHomework() {
	//	get teacher course
	if t_course_id, err := this.GetInt64("course"); err == nil {
		if t_course, err := models.GetTeacherCourseById(t_course_id); err == nil {
			//	get student course homework list
			if homeworks, err := models.GetStudentHomeworkByTeacherCourseAndStudent(t_course, this.Data["student"].(*models.Student).Id); err == nil {
				t_course.Orgs["s_homeworks"] = homeworks
				this.Data["t_course"] = t_course
				//	get the term all teachercourse
				if s_courses, err := models.GetStudentCourseByTerm(t_course.Term, this.Data["student"].(*models.Student).Id); err == nil {
					this.Data["s_courses"] = s_courses
				} else {
					this.Redirect("/error", 302)
					this.StopRun()
				}
				//	get all term
				term_list, err := models.GetTermListByStudentHomework(this.Data["student"].(*models.Student).Id)
				if err != nil {
					this.Redirect("/error", 302)
					this.StopRun()
				}
				this.Data["term_list"] = models.RankingTerm(term_list)
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

	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "student/head/courseHomework_head.html"
	this.LayoutSections["Scripts"] = "student/scripts/courseHomework_scripts.html"
	this.TplNames = "student/courseHomework.html"
}

// @Title 签到统计
// @router /studentCheck [get]
func (this *ViewStudentController) StudentCheck() {
	//	get the term list
	term_list, err := models.GetTermListByStudentCourse(this.Data["student"].(*models.Student).Id)
	if err != nil {
		this.Redirect("/error", 302)
		this.StopRun()
	}
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
	//	get the student course
	s_course_id, _ := this.GetInt64("course")
	if s_course_id > 0 {
		if s_course, err := models.GetStudentCourseById(s_course_id); err == nil {
			this.Data["s_course"] = s_course
		} else {
			this.Redirect("/error", 302)
			this.StopRun()
		}
	}
	//	get the student courses
	if s_courses, err := models.GetStudentCourseByTerm(term, this.Data["student"].(*models.Student).Id); err == nil {
		//	get the student homework
		for k, v := range s_courses {
			if s_checks, err := models.GetStudentChecksByTeacherCourseAndStudent(v.TeacherCourse, this.Data["student"].(*models.Student).Id); err == nil {
				s_courses[k].Orgs = make(map[string]interface{})
				s_courses[k].Orgs["s_checks"] = s_checks
				for k, v := range term_list {
					if term.Id == v.Id {
						term_list[k].Orgs = make(map[string]interface{})
						term_list[k].Orgs["s_course"] = s_courses
					}
				}
			} else {
				continue
			}
		}
		this.Data["s_courses"] = s_courses
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	//	get all student checks
	for k, v := range term_list {
		if term.Id != v.Id {
			term_list[k].Orgs = make(map[string]interface{})

			if s_courses, err := models.GetStudentCourseByTerm(term_list[k], this.Data["student"].(*models.Student).Id); err == nil {
				//	get the homework
				for k, v := range s_courses {
					if s_checks, err := models.GetStudentChecksByTeacherCourseAndStudent(v.TeacherCourse, this.Data["student"].(*models.Student).Id); err == nil {
						s_courses[k].Orgs = make(map[string]interface{})
						s_courses[k].Orgs["s_checks"] = s_checks
					} else {
						continue
					}
				}
				term_list[k].Orgs["s_course"] = s_courses
			}
		}
	}
	this.Data["term_list"] = models.RankingTerm(term_list)

	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "student/head/check_head.html"
	this.LayoutSections["Scripts"] = "student/scripts/check_scripts.html"
	this.TplNames = "student/check.html"
}

// @Title 作业统计
// @router /studentHomework [get]
func (this *ViewStudentController) StudentHomework() {
	//	get the term list
	term_list, err := models.GetTermListByStudentCourse(this.Data["student"].(*models.Student).Id)
	if err != nil {
		this.Redirect("/error", 302)
		this.StopRun()
	}
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
	//	get the student course
	s_course_id, _ := this.GetInt64("course")
	if s_course_id > 0 {
		if s_course, err := models.GetStudentCourseById(s_course_id); err == nil {
			this.Data["s_course"] = s_course
		} else {
			this.Redirect("/error", 302)
			this.StopRun()
		}
	}
	//	get the student courses
	if s_courses, err := models.GetStudentCourseByTerm(term, this.Data["student"].(*models.Student).Id); err == nil {
		//	get the student homework
		for k, v := range s_courses {
			if s_homeworks, err := models.GetStudentHomeworkByStudentCourse(v, this.Data["student"].(*models.Student).Id); err == nil {
				s_courses[k].Orgs = make(map[string]interface{})
				s_courses[k].Orgs["s_homeworks"] = s_homeworks
				for k, v := range term_list {
					if term.Id == v.Id {
						term_list[k].Orgs = make(map[string]interface{})
						term_list[k].Orgs["s_course"] = s_courses
					}
				}
			} else {
				continue
			}
		}
		this.Data["s_courses"] = s_courses
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	//	get all student homework
	for k, v := range term_list {
		if term.Id != v.Id {
			term_list[k].Orgs = make(map[string]interface{})

			if s_courses, err := models.GetStudentCourseByTerm(term_list[k], this.Data["student"].(*models.Student).Id); err == nil {
				//	get the homework
				for k, v := range s_courses {
					if s_homeworks, err := models.GetStudentHomeworkByStudentCourse(v, this.Data["student"].(*models.Student).Id); err == nil {
						s_courses[k].Orgs = make(map[string]interface{})
						s_courses[k].Orgs["s_homeworks"] = s_homeworks
					} else {
						continue
					}
				}
				term_list[k].Orgs["s_course"] = s_courses
			}
		}
	}
	this.Data["term_list"] = models.RankingTerm(term_list)

	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "student/head/homework_head.html"
	this.LayoutSections["Scripts"] = "student/scripts/homework_scripts.html"
	this.TplNames = "student/homework.html"
}

// @Title 下载附件
// @router /download [get]
func (this *ViewStudentController) Download() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "学生" {
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

// @Tilte 我的成绩
// @router /studentScore [get]
func (this *ViewStudentController) StudentScore() {
	//	get the term list
	term_list, err := models.GetTermListByStudentCourse(this.Data["student"].(*models.Student).Id)
	if err != nil {
		this.Redirect("/error", 302)
		this.StopRun()
	}

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
	//	get the student course
	if s_courses, err := models.GetStudentCourseByTerm(term, this.Data["student"].(*models.Student).Id); err == nil {
		s_courses = models.FilterRepeat(s_courses).([]*models.StudentCourse)
		//	get the score
		for k, v := range s_courses {
			if score, err := models.SearchStudentGrade(v.Id); err == nil {
				s_courses[k].Score = score
			} else {
				continue
			}
		}
		this.Data["s_courses"] = s_courses
		for k, v := range term_list {
			if term.Id == v.Id {
				term_list[k].Orgs = make(map[string]interface{})
				for k1, _ := range s_courses {
					s_courses[k1].TeacherCourse.Term = nil
					s_courses[k1].Score.StudentCourse = nil
				}
				term_list[k].Orgs["s_course"] = s_courses
			}
		}
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	//	get all student course
	for k, v := range term_list {
		if term.Id != v.Id {
			term_list[k].Orgs = make(map[string]interface{})

			if s_courses, err := models.GetStudentCourseByTerm(term_list[k], this.Data["student"].(*models.Student).Id); err == nil {
				s_courses = models.FilterRepeat(s_courses).([]*models.StudentCourse)
				//	get the score
				for k, v := range s_courses {
					if score, err := models.SearchStudentGrade(v.Id); err == nil {
						score.StudentCourse = nil
						s_courses[k].TeacherCourse.Term = nil
						s_courses[k].Score = score
					} else {
						continue
					}
				}
				term_list[k].Orgs["s_course"] = s_courses
			}
		}
	}
	this.Data["term_list"] = models.RankingTerm(term_list)

	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "student/head/score_head.html"
	this.LayoutSections["Scripts"] = "student/scripts/score_scripts.html"
	this.TplNames = "student/score.html"
}

// @Tilte 我的统计数据
// @router /echart [get]
func (this *ViewStudentController) Echart() {
	//	get the term list
	term_list, err := models.GetTermListByStudentCourse(this.Data["student"].(*models.Student).Id)
	if err != nil {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	//	get all student course
	for k, _ := range term_list {
		term_list[k].Orgs = make(map[string]interface{})

		if s_courses, err := models.GetStudentCourseByTerm(term_list[k], this.Data["student"].(*models.Student).Id); err == nil {
			s_courses = models.FilterRepeat(s_courses).([]*models.StudentCourse)
			//	get the score
			for k, v := range s_courses {
				if score, err := models.SearchStudentGrade(v.Id); err == nil {
					score.StudentCourse = nil
					s_courses[k].TeacherCourse.Term = nil
					s_courses[k].Score = score
				} else {
					continue
				}
				if s_homeworks, err := models.GetStudentHomeworkByStudentCourse(v, this.Data["student"].(*models.Student).Id); err == nil {
					s_courses[k].Orgs = make(map[string]interface{})
					s_courses[k].Orgs["s_homeworks"] = s_homeworks
				} else {
					continue
				}
				if s_checks, err := models.GetStudentChecksByTeacherCourseAndStudent(v.TeacherCourse, this.Data["student"].(*models.Student).Id); err == nil {
					s_courses[k].Orgs["s_checks"] = s_checks
				} else {
					continue
				}
			}
			term_list[k].Orgs["s_course"] = s_courses
		}
	}
	this.Data["term_list"] = models.RankingTerm(term_list)

	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "student/head/info_head.html"
	this.LayoutSections["Scripts"] = "student/scripts/echart_scripts.html"
	this.TplNames = "student/echart.html"
}

// @Title 登录教务系统
// @router /login [get]
func (this *ViewStudentController) Login() {
	//	recevie data
	user := this.Data["student"]
	if user != nil && user.(*models.Student).IsEdu {
		this.Redirect("/view/student/eduManage", 302)
		this.StopRun()
	}
	//	取出 refer
	refer := this.Ctx.Request.Referer()
	if len(refer) > 0 && refer == "http://"+models.Host+":"+models.Port+"/view/student/login" {
		this.Data["fail"] = true
	}
	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "student/head/info_head.html"
	this.LayoutSections["Scripts"] = "scripts/signature_scripts.html"
	this.TplNames = "student/login.html"
}

// @Title 教务系统导入
// @router /eduLoading [get]
func (this *ViewStudentController) EduLoading() {
	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "student/scripts/eduLoading_scripts.html"
	this.TplNames = "student/eduLoading.html"
}

// @Title 教务系统管理
// @router /eduManage [get]
func (this *ViewStudentController) EduManage() {
	//	recevie data
	u_id := this.Data["student"].(*models.Student).Id
	if userinfo, err := Edu.GetUserinfoLog(u_id); err == nil {
		var m map[string]string
		if err := json.Unmarshal([]byte(userinfo.Content), &m); err == nil {
			this.Data["userinfo"] = m
		}
	}
	if department_log, err := Edu.GetUserDepartmentLog(u_id); err == nil {
		this.Data["department_log"] = department_log
	}
	if major_log, err := Edu.GetUserMajorLog(u_id); err == nil {
		this.Data["major_log"] = major_log
	}
	if class_log, err := Edu.GetUserClassLog(u_id); err == nil {
		this.Data["class_log"] = class_log
	}
	if failed_log, err := Edu.GetUserFailedLog(u_id); err == nil {
		this.Data["failed_log"] = failed_log
	}
	if term_log, err := Edu.GetUserTermLog(u_id); err == nil {
		this.Data["term_log"] = term_log
	}
	if course_log, err := Edu.GetUserCourseLog(u_id); err == nil {
		this.Data["course_log"] = course_log
	}
	if teacher_log, err := Edu.GetUserTeacherLog(u_id); err == nil {
		this.Data["teacher_log"] = teacher_log
	}
	if t_course_log, err := Edu.GetUserTeacherCourseLog(u_id); err == nil {
		this.Data["t_course_log"] = t_course_log
	}
	if s_course_log, err := Edu.GetUserStudentCourseLog(u_id); err == nil {
		this.Data["s_course_log"] = s_course_log
	}
	if t_course_class_log, err := Edu.GetUserTeacherCourseClassLog(u_id); err == nil {
		this.Data["t_course_class_log"] = t_course_class_log
	}
	if s_course_grade_log, err := Edu.GetUserStudentCourseGradeLog(u_id); err == nil {
		this.Data["s_course_grade_log"] = s_course_grade_log
	}

	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Head_html"] = "student/head/eduManage_head.html"
	this.LayoutSections["Scripts"] = "student/scripts/eduManage_scripts.html"
	this.TplNames = "student/eduManage.html"
}
