package controllers

import (
	"fmt"
	"learn/models"
	"strconv"
	"strings"
)

type APITeacherController struct {
	baseController
}

func (this *APITeacherController) Prepare() {
	//	recevie the key
	key := this.GetString("key")
	if len(key) <= 0 && key != models.Str2Sha1(this.Ctx.Input.Cookie("beegosessionID")) {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	//	recevie the signature
	signature := this.GetString("signature")
	if signature != this.GetSession("signature") {
		this.Redirect("/error", 302)
		this.StopRun()
	}
}

// @Title 修改密码
// @router /updatePassword [post]
func (this *APITeacherController) UpdatePassword() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "教师" {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	t_id, _ := strconv.ParseInt(this.GetSession("id").(string), 10, 64)
	if user, err := models.GetTeacherById(t_id); err == nil {
		newpwd := this.GetString("newpwd")
		chkpwd := this.GetString("chkpwd")
		if newpwd != chkpwd {
			this.Redirect("/error", 302)
			this.StopRun()
		}
		if oldpwd := this.GetString("oldpwd"); user.Password == oldpwd {
			//	update pwd
			user.Password = newpwd
			if err := models.UpdateTeacher(user); err != nil {
				this.Redirect("/error", 302)
				this.StopRun()
			}
			this.Redirect("/view/teacher/table?change=true", 302)
			this.StopRun()
		}
	}
	this.Redirect("/error", 302)
	this.StopRun()
}

// @Title 获取学生提交的作业
// @router /getStudentHomeworksByTeacherCourseHomework [post]
func (this *APITeacherController) GetStudentHomeworksByTeacherCourseHomework() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "教师" {
		this.SendFailedJSON(1, "用户类型不正确")
		return
	}
	if t_course_homework_id, err := this.GetInt64("t_course_homework"); err == nil {
		if t_course_homework, err := models.GetTeacherCourseHomeworkById(t_course_homework_id); err == nil {
			t_id, _ := strconv.ParseInt(this.GetSession("id").(string), 10, 64)
			if t_course_homework.TeacherCourse.Teacher.Id == t_id {
				//	get student homework list
				if s_homeworks, err := models.GetUploadedStudentHomeworkByTeacherCourseHomework(t_course_homework); err == nil {
					this.SendSuccessJSON(s_homeworks)
					return
				}
				this.SendFailedJSON(1, "获取学生作业失败")
				return
			}
			this.SendFailedJSON(1, "用户不匹配")
			return
		}
		this.SendFailedJSON(1, "获取教师作业失败")
		return
	}
	this.SendFailedJSON(1, "获取教师作业id失败")
	return
}

// @Ttile 给学生作业打分
// @router /setGradeByStudentHomework [post]
func (this *APITeacherController) SetGradeByStudentHomework() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "教师" {
		this.SendFailedJSON(1, "用户类型不正确")
		return
	}
	if s_homework_id, err := this.GetInt64("s_homework"); err == nil {
		if s_homework, err := models.GetStudentHomeworkById(s_homework_id); err == nil {
			t_id, _ := strconv.ParseInt(this.GetSession("id").(string), 10, 64)
			if s_homework.TeacherCourseHomework.TeacherCourse.Teacher.Id == t_id {
				s_homework.Grade = this.GetString("grade")
				if err := models.UpdateStudentHomework(s_homework); err == nil {
					this.SendSuccessJSON("操作成功")
					return
				}
				this.SendFailedJSON(1, "评分失败")
				return
			}
			this.SendFailedJSON(1, "用户不匹配")
			return
		}
		this.SendFailedJSON(1, "获取学生作业失败")
		return
	}
	this.SendFailedJSON(1, "获取学生作业id失败")
	return
}

// @Title 查询个人课表
// @router /teacherCourse [post]
func (this *APITeacherController) TeacherCourse() {
	xnd := this.GetString("xnd")
	xqd := this.GetString("xqd")
	week := this.GetString("week")
	if len(xnd) <= 0 || len(xqd) <= 0 {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	if len(week) > 0 {
		week_num, _ := strconv.Atoi(week)
		if week_num <= 0 || week_num > 20 {
			this.Redirect("/error", 302)
			this.StopRun()
		}
		this.Redirect("/view/teacher/table?xnd="+xnd+"&xqd="+xqd+"&week="+fmt.Sprintf("%d", week_num), 302)
		this.StopRun()
	}
	this.Redirect("/view/teacher/table?xnd="+xnd+"&xqd="+xqd, 302)
	this.StopRun()
}

// @Title 查询作业
// @router /courseHomework [post]
func (this *APITeacherController) CourseHomework() {
	if course_id, err := this.GetInt64("course"); err == nil {
		if course_id > 0 {
			this.Redirect("/view/teacher/courseHomework?course="+fmt.Sprintf("%d", course_id), 302)
			this.StopRun()
		}
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}
}

// @Title 查询作业
// @router /courseHistory [post]
func (this *APITeacherController) CourseHistory() {
	if course_id, err := this.GetInt64("course"); err == nil {
		if course_id > 0 {
			week := this.GetString("week")
			if len(week) > 0 {
				week_num, _ := strconv.Atoi(week)
				if week_num <= 0 || week_num > 20 {
					this.Redirect("/error", 302)
					this.StopRun()
				}
				this.Redirect("/view/teacher/courseHistory?course="+fmt.Sprintf("%d", course_id)+"&week="+fmt.Sprintf("%d", week_num), 302)
				this.StopRun()
			}
			this.Redirect("/view/teacher/courseHistory?course="+fmt.Sprintf("%d", course_id), 302)
			this.StopRun()
		}
	}
	this.Redirect("/error", 302)
	this.StopRun()
}

// @Title 学生签到
// @router /studentCheck [post]
func (this *APITeacherController) StudentCheck() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "教师" {
		this.SendFailedJSON(1, "操作失败")
		return
	}
	//	get teacher course
	if t_course_id, err := this.GetInt64("t_course"); err == nil {
		if t_course, err := models.GetTeacherCourseById(t_course_id); err == nil {
			t_id, _ := strconv.ParseInt(this.GetSession("id").(string), 10, 64)
			if t_course.Teacher.Id == t_id {
				week, _ := this.GetInt("week")
				result := this.GetString("result")
				if result != "迟到" && result != "已到" && result != "未到" && result != "请假" {
					this.SendFailedJSON(1, "点名结果不正确")
					return
				}
				var tmp_s_check = &models.StudentCheck{
					Result:        result,
					Student:       &models.Student{Id: this.GetString("student")},
					Week:          week,
					TeacherCourse: t_course,
				}
				if exist := models.ExistStudentCheck(tmp_s_check); !exist {
					if err := models.AddStudentCheck(tmp_s_check); err == nil {
						this.SendSuccessJSON("操作成功")
						return
					}
					this.SendFailedJSON(1, "添加学生签到记录失败")
					return
				} else {
					//	search student check
					if s_check, err := models.SearchStudentCheck(t_course.Id, tmp_s_check.Student.Id, tmp_s_check.Week); err == nil {
						tmp_s_check.Id = s_check.Id
						if err := models.UpdateStudentCheck(tmp_s_check); err == nil {
							this.SendSuccessJSON("操作成功")
							return
						}
						this.SendFailedJSON(1, "更新学生签到记录失败")
						return
					}
					this.SendFailedJSON(1, "获取学生签到记录失败")
					return
				}
			}
			this.SendFailedJSON(1, "用户不匹配")
			return
		}
		this.SendFailedJSON(1, "获取教师课程失败")
		return
	}
	this.SendFailedJSON(1, "获取教师课程id失败")
	return
}

// @Title 更新头像
// @router /updateImg [post]
func (this *APITeacherController) UpdateImg() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "教师" {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	t_id, _ := strconv.ParseInt(this.GetSession("id").(string), 10, 64)
	if user, err := models.GetTeacherById(t_id); err == nil {
		//	recevie the img
		user.Headimgurl = this.GetString("img")
		//	update teacher
		if err := models.UpdateTeacher(user); err != nil {
			this.Redirect("/error", 302)
			this.StopRun()
		}
		this.Redirect("/view/teacher/table?change=true", 302)
		this.StopRun()
	}
	this.Redirect("/error", 302)
	this.StopRun()
}

// @Title 上传头像
// @router /uploadImg [post]
func (this *APITeacherController) UploadImg() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "教师" {
		this.SendFailedJSON(1, "操作失败")
		return
	}
	//	recevie the img
	path := "UPLOADS/headimg/"
	models.CreateDir(path)
	filepath := path + models.GetRandString(2014) + ".jpeg"
	if err := this.SaveToFile("img", filepath); err != nil {
		this.SendFailedJSON(1, "操作失败")
		return
	}
	this.SendSuccessJSON("/" + filepath)
	return
}

// @Title 发布作业
// @router /publisHomework [post]
func (this *APITeacherController) PublisHomework() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "教师" {
		this.SendFailedJSON(1, "操作失败")
		return
	}
	//	judge teacher course' teacher is true
	t_id, _ := strconv.ParseInt(this.GetSession("id").(string), 10, 64)
	t_course_id, _ := this.GetInt64("t_course")
	if t_course, err := models.GetTeacherCourseById(t_course_id); err == nil {
		if t_course.Teacher.Id != t_id {
			this.SendFailedJSON(1, "教师不匹配")
			return
		}
		//	recevie the data
		publish_week, _ := this.GetInt("publish_week")
		var t_course_homework = &models.TeacherCourseHomework{
			Title:         this.GetString("title"),
			Remark:        this.GetString("remark"),
			AsOfTime:      this.GetString("as_of_time"),
			Attachment:    this.GetString("attachment"),
			PublishWeek:   publish_week,
			TeacherCourse: t_course,
		}
		//	judge the teacher coursehomework exist
		if exist := models.ExistTeacherCourseHomework(t_course_homework); exist {
			this.SendFailedJSON(1, "作业已存在")
			return
		}
		//	add the teacher course homework
		if t_course_homework_id, err := models.AddTeacherCourseHomework(t_course_homework); err == nil {
			//	get students
			if students, err := models.GetStudentsByTeacherCourse(t_course); err == nil {
				for k, _ := range students {
					//	add studentHomework and studentNotice
					var student_homwork = &models.StudentHomework{
						Student:               students[k],
						TeacherCourseHomework: &models.TeacherCourseHomework{Id: t_course_homework_id},
					}
					models.AddStudentHomework(student_homwork)
					var student_notice = &models.StudentNotice{
						Student:    students[k],
						SenderType: "teacher",
						Teacher:    t_course.Teacher,
						Content:    `课程《` + t_course.Course.Name + `》发布了新的作业，请同学们及时提交作业 <a href="/view/student/courseHomework?course=` + fmt.Sprintf("%d", t_course.Id) + `">课程链接</a>`,
					}
					models.AddStudentNotice(student_notice)
				}
			}
			this.SendSuccessJSON("操作成功")
			return
		} else {
			this.SendFailedJSON(1, "操作失败")
			return
		}
	} else {
		this.SendFailedJSON(1, "操作失败")
		return
	}
}

// @Title 更新作业
// @router /updateHomework [post]
func (this *APITeacherController) UpdateHomework() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "教师" {
		this.SendFailedJSON(1, "操作失败")
		return
	}
	//	judge teacher course' teacher is true
	t_id, _ := strconv.ParseInt(this.GetSession("id").(string), 10, 64)
	t_course_id, _ := this.GetInt64("t_course")
	t_course_homeworl_id, _ := this.GetInt64("t_course_homework")
	if t_course, err := models.GetTeacherCourseById(t_course_id); err == nil {
		if t_course.Teacher.Id != t_id {
			this.SendFailedJSON(1, "教师不匹配")
			return
		}
		//	recevie the data
		publish_week, _ := this.GetInt("publish_week")
		var t_course_homework = &models.TeacherCourseHomework{
			Id:            t_course_homeworl_id,
			Title:         this.GetString("title"),
			Remark:        this.GetString("remark"),
			AsOfTime:      this.GetString("as_of_time"),
			Attachment:    this.GetString("attachment"),
			PublishWeek:   publish_week,
			TeacherCourse: t_course,
		}
		//	update the teacher course homework
		if err := models.UpdateTeacherCourseHomework(t_course_homework); err == nil {
			this.SendSuccessJSON("操作成功")
			return
		} else {
			this.SendFailedJSON(1, "操作失败")
			return
		}
	} else {
		this.SendFailedJSON(1, "操作失败")
		return
	}
}

// @Title 上传作业附件
// @router /uploadAttachment [post]
func (this *APITeacherController) UploadAttachment() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "教师" {
		this.SendFailedJSON(1, "操作失败")
		return
	}
	//	recevie the attachment
	path := "UPLOADS/attachment/"
	models.CreateDir(path)
	filepath := path + models.GetRandString(2014) + ".zip"
	if err := this.SaveToFile("attachment", filepath); err != nil {
		this.SendFailedJSON(1, err)
		return
	}
	this.SendSuccessJSON("/" + filepath)
	return
}

// @Title 更新个人资料
// @router /updateProfile [post]
func (this *APITeacherController) UpdateProfile() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "教师" {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	t_id, _ := strconv.ParseInt(this.GetSession("id").(string), 10, 64)
	if user, err := models.GetTeacherById(t_id); err == nil {
		//	recevie the data
		sex, _ := this.GetInt("sex")
		user.Sex = sex
		user.Profile = &models.TeacherProfile{
			Teacher: &models.Teacher{Id: user.Id},
			Phone:   this.GetString("phone"),
			Email:   this.GetString("email"),
			Degree:  this.GetString("degree"),
			Title:   this.GetString("title"),
			Subject: this.GetString("subject"),
			Remark:  this.GetString("remark"),
		}
		//	update teacher
		if err := models.UpdateTeacherProfile(user); err != nil {
			this.Redirect("/error", 302)
			this.StopRun()
		}
		this.Redirect("/view/teacher/table?change=true", 302)
		this.StopRun()
	}
	this.Redirect("/error", 302)
	this.StopRun()
}

// @Title 删除教师作业
// @router /delTeacherCourseHomework [post]
func (this *APITeacherController) DelTeacherCourseHomework() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "教师" {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	if t_course_homework_id, err := this.GetInt64("t_course_homework"); err == nil {
		if t_course_homework_id > 0 {
			if t_course_homework, err := models.GetTeacherCourseHomeworkById(t_course_homework_id); err == nil {
				t_id, _ := strconv.ParseInt(this.GetSession("id").(string), 10, 64)
				if t_course_homework.TeacherCourse.Teacher.Id != t_id {
					this.SendFailedJSON(1, "用户不匹配")
					return
				}
				//	删除所有学生作业
				models.DelStudentHomeworkByTeacherCourseHomework(t_course_homework)
				if err := models.DelTeacherCourseHomework(t_course_homework); err != nil {
					this.SendFailedJSON(1, "删除失败")
					return
				}
				this.SendSuccessJSON("操作成功")
				return
			} else {
				this.SendFailedJSON(1, "获取失败")
				return
			}
		} else {
			this.SendFailedJSON(1, "id不正确")
			return
		}
	} else {
		this.SendFailedJSON(1, "id不正确")
		return
	}
}

// @Title 获取教师作业
// @router /getTeacherCourseHomework [post]
func (this *APITeacherController) GetTeacherCourseHomework() {
	if t_course_homework_id, err := this.GetInt64("t_course_homework"); err == nil {
		if t_course_homework_id > 0 {
			if t_course_homework, err := models.GetTeacherCourseHomeworkById(t_course_homework_id); err == nil {
				this.SendSuccessJSON(t_course_homework)
				return
			} else {
				this.SendFailedJSON(1, "获取失败")
				return
			}
		} else {
			this.SendFailedJSON(1, "id不正确")
			return
		}
	} else {
		this.SendFailedJSON(1, "id不正确")
		return
	}
}

// @Title 根据学年度 与 教师获取学期
// @router /getTermNumberByTeacher [post]
func (this *APITeacherController) GetTermNumberByTeacher() {
	xnd := this.GetString("xnd")
	if len(xnd) > 0 {
		arr := strings.Split(xnd, "-")
		if len(arr) == 2 {
			startYear, _ := strconv.Atoi(arr[0])
			endYear, _ := strconv.Atoi(arr[1])
			t_id, _ := strconv.ParseInt(this.GetSession("id").(string), 10, 64)
			if terms, err := models.GetTermNumberByTeacher(t_id, startYear, endYear); err == nil {
				this.SendSuccessJSON(terms)
				return
			} else {
				this.SendFailedJSON(1, "获取学期数据失败")
				return
			}
		}
	}
	this.SendFailedJSON(1, "学年度不正确")
	return
}

// @Title 根据学年度 学期 与 教师 获取教师课程表
// @router /getTeacherCourseByTerm [post]
func (this *APITeacherController) GetTeacherCourseByTerm() {
	xnd := this.GetString("xnd")
	if len(xnd) > 0 {
		arr := strings.Split(xnd, "-")
		if len(arr) == 2 {
			startYear, _ := strconv.Atoi(arr[0])
			endYear, _ := strconv.Atoi(arr[1])
			xqd, _ := this.GetInt("xqd")
			if term, err := models.SearchTerm(xqd, startYear, endYear); err == nil {
				t_id, _ := strconv.ParseInt(this.GetSession("id").(string), 10, 64)
				if t_courses, err := models.GetTeacherCourseByTerm(term, t_id); err == nil {
					this.SendSuccessJSON(t_courses)
					return
				} else {
					this.SendFailedJSON(1, "获取教师课程表失败")
					return
				}
			} else {
				this.SendFailedJSON(1, "获取学期数据失败")
				return
			}

		}
	}
	this.SendFailedJSON(1, "学年度不正确")
	return
}
