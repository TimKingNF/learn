package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	Edu "learn/EDU"
	"learn/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type APIStudentController struct {
	baseController
}

func (this *APIStudentController) Prepare() {
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
func (this *APIStudentController) UpdatePassword() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "学生" {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	id := this.GetSession("id").(string)
	if user, err := models.GetStudentById(id); err == nil {
		newpwd := this.GetString("newpwd")
		chkpwd := this.GetString("chkpwd")
		if newpwd != chkpwd {
			this.Redirect("/error", 302)
			this.StopRun()
		}
		if oldpwd := this.GetString("oldpwd"); user.Password == oldpwd {
			//	update pwd
			user.Password = newpwd
			if err := models.UpdateStudent(user); err != nil {
				this.Redirect("/error", 302)
				this.StopRun()
			}
			this.Redirect("/view/student/table?change=true", 302)
			this.StopRun()
		}
	}
	this.Redirect("/error", 302)
	this.StopRun()
}

// @Title 更新头像
// @router /updateImg [post]
func (this *APIStudentController) UpdateImg() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "学生" {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	id := this.GetSession("id").(string)
	if user, err := models.GetStudentById(id); err == nil {
		//	recevie the img
		user.Headimgurl = this.GetString("img")
		//	update student
		if err := models.UpdateStudent(user); err != nil {
			this.Redirect("/error", 302)
			this.StopRun()
		}
		this.Redirect("/view/student/table?change=true", 302)
		this.StopRun()
	}
	this.Redirect("/error", 302)
	this.StopRun()
}

// @Title 学生上传作业
// @router /uploadAttachment [post]
func (this *APIStudentController) UploadAttachment() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "学生" {
		this.SendFailedJSON(1, "操作失败")
		return
	}
	if s_homework_id, err := this.GetInt64("student_homework"); err == nil {
		if s_homework, err := models.GetStudentHomeworkById(s_homework_id); err == nil {
			if s_homework.Student.Id != this.GetSession("id").(string) {
				this.SendFailedJSON(1, "用户不匹配")
				return
			}
			if time.Now().After(models.GetTime(s_homework.TeacherCourseHomework.AsOfTime)) {
				this.SendFailedJSON(1, "已经截止提交")
				return
			}
			//	get attachment
			path := "UPLOADS/homework/"
			models.CreateDir(path)
			filepath := path + models.GetRandString(2014) + ".zip"
			if err := this.SaveToFile("attachment", filepath); err != nil {
				this.SendFailedJSON(1, err)
				return
			}
			s_homework.Attachment = "/" + filepath
			s_homework.IsUpload = 1
			if err := models.UpdateStudentHomework(s_homework); err == nil {
				this.SendSuccessJSON("操作成功")
				return
			} else {
				this.SendFailedJSON(1, "更新学生作业失败")
				return
			}
		} else {
			this.SendFailedJSON(1, "获取学生作业失败")
			return
		}

	} else {
		this.SendFailedJSON(1, "获取学生作业id失败")
		return
	}
	this.SendFailedJSON(1, "操作失败")
	return
}

// @Title 阅读消息
// @router /readNoticeByStudent [post]
func (this *APIStudentController) ReadNoticeByStudent() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "学生" {
		this.SendFailedJSON(1, "操作失败")
		return
	}
	s_id := this.GetSession("id").(string)
	if models.CountNotReadStudentNotice(s_id) > 0 {
		send_id, _ := this.GetInt64("sender")
		if err := models.ReadNoticeByStudent(s_id, this.GetString("sender_type"), send_id); err != nil {
			this.SendFailedJSON(1, "操作失败")
			return
		}
	}
	this.SendSuccessJSON("操作成功")
	return
}

// @Title 删除消息
// @router /deleteStudentNotice [post]
func (this *APIStudentController) DeleteStudentNotice() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "学生" {
		this.SendFailedJSON(1, "操作失败")
		return
	}
	if notice_id, err := this.GetInt64("notice"); err == nil {
		//	search notice
		if s_notice, err := models.GetStudentNoticeById(notice_id); err == nil {
			if s_notice.Student.Id != this.GetSession("id").(string) {
				this.SendFailedJSON(1, "用户不匹配")
				return
			}
			//	delete student notice
			if err := models.DelStudentNotice(s_notice); err == nil {
				this.SendSuccessJSON("操作成功")
				return
			}
		}
	}
	this.SendFailedJSON(1, "获取id失败")
	return
}

// @Title 上传头像
// @router /uploadImg [post]
func (this *APIStudentController) UploadImg() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "学生" {
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

// @Title 查询个人课表
// @router /studentCourse [post]
func (this *APIStudentController) StudentCourse() {
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
		this.Redirect("/view/student/table?xnd="+xnd+"&xqd="+xqd+"&week="+fmt.Sprintf("%d", week_num), 302)
		this.StopRun()
	}
	this.Redirect("/view/student/table?xnd="+xnd+"&xqd="+xqd, 302)
	this.StopRun()
}

// @Title 查询个人作业统计
// @router /studentHomework [post]
func (this *APIStudentController) StudentHomework() {
	xnd := this.GetString("xnd")
	xqd := this.GetString("xqd")
	course := this.GetString("course")
	if len(xnd) <= 0 || len(xqd) <= 0 {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	if len(course) > 0 {
		course_id, _ := strconv.Atoi(course)
		if course_id <= 0 {
			this.Redirect("/error", 302)
			this.StopRun()
		}
		this.Redirect("/view/student/studentHomework?xnd="+xnd+"&xqd="+xqd+"&course="+fmt.Sprintf("%d", course_id), 302)
		this.StopRun()
	}
	this.Redirect("/view/student/studentHomework?xnd="+xnd+"&xqd="+xqd, 302)
	this.StopRun()
}

// @Title 查询个人签到统计
// @router /studentCheck [post]
func (this *APIStudentController) StudentCheck() {
	xnd := this.GetString("xnd")
	xqd := this.GetString("xqd")
	course := this.GetString("course")
	if len(xnd) <= 0 || len(xqd) <= 0 {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	if len(course) > 0 {
		course_id, _ := strconv.Atoi(course)
		if course_id <= 0 {
			this.Redirect("/error", 302)
			this.StopRun()
		}
		this.Redirect("/view/student/studentCheck?xnd="+xnd+"&xqd="+xqd+"&course="+fmt.Sprintf("%d", course_id), 302)
		this.StopRun()
	}
	this.Redirect("/view/student/studentCheck?xnd="+xnd+"&xqd="+xqd, 302)
	this.StopRun()
}

// @Title 查询个人课程签到记录
// @router /courseHistory [post]
func (this *APIStudentController) CourseHistory() {
	if course_id, err := this.GetInt64("course"); err == nil {
		if course_id > 0 {
			this.Redirect("/view/student/courseHistory?course="+fmt.Sprintf("%d", course_id), 302)
			this.StopRun()
		}
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}
}

// @Title 查询个人成绩
// @router /studentScore [post]
func (this *APIStudentController) StudentScore() {
	xnd := this.GetString("xnd")
	xqd := this.GetString("xqd")
	if len(xnd) <= 0 || len(xqd) <= 0 {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	this.Redirect("/view/student/studentScore?xnd="+xnd+"&xqd="+xqd, 302)
	this.StopRun()
}

// @Title 查询作业
// @router /courseHomework [post]
func (this *APIStudentController) CourseHomework() {
	if course_id, err := this.GetInt64("course"); err == nil {
		if course_id > 0 {
			this.Redirect("/view/student/courseHomework?course="+fmt.Sprintf("%d", course_id), 302)
			this.StopRun()
		}
	} else {
		this.Redirect("/error", 302)
		this.StopRun()
	}
}

// @Title 根据学年度 与 学年度获取学期
// @router /getTermNumberByStudent [post]
func (this *APIStudentController) GetTermNumberByStudent() {
	xnd := this.GetString("xnd")
	if len(xnd) > 0 {
		arr := strings.Split(xnd, "-")
		if len(arr) == 2 {
			startYear, _ := strconv.Atoi(arr[0])
			endYear, _ := strconv.Atoi(arr[1])
			s_id := this.GetSession("id").(string)
			if terms, err := models.GetTermNumberByStudent(s_id, startYear, endYear); err == nil {
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
// @router /getStudentCourseByTerm [post]
func (this *APIStudentController) GetStudentCourseByTerm() {
	xnd := this.GetString("xnd")
	if len(xnd) > 0 {
		arr := strings.Split(xnd, "-")
		if len(arr) == 2 {
			startYear, _ := strconv.Atoi(arr[0])
			endYear, _ := strconv.Atoi(arr[1])
			xqd, _ := this.GetInt("xqd")
			if term, err := models.SearchTerm(xqd, startYear, endYear); err == nil {
				s_id := this.GetSession("id").(string)
				if s_courses, err := models.GetStudentCourseByTerm(term, s_id); err == nil {
					this.SendSuccessJSON(s_courses)
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

// @Title 登录教务系统
// @router /login [post]
func (this *APIStudentController) Login() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "学生" {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	id := this.GetString("id")
	pwd := this.GetString("pwd")
	if len(id) <= 0 || len(pwd) <= 0 {
		this.Redirect("/view/student/login", 302)
		this.StopRun()
	}
	//	login edu
	if _, ok, _, _ := Edu.Sign_in(id, pwd, "学生"); !ok {
		this.Redirect("/view/student/login", 302)
		this.StopRun()
	}
	//	update student
	if user, err := models.GetStudentById(id); err == nil {
		user.EduPwd = pwd
		user.IsEdu = true
		if err := models.UpdateStudent(user); err == nil {
			this.Redirect("/view/student/eduLoading", 302)
			this.StopRun()
		}
	}
	this.Redirect("/view/student/login", 302)
	this.StopRun()
}

// @Title 导入教务系统
// @router /eduLoading [get]
func (this *APIStudentController) EduLoading() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "学生" {
		return
	}
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	defer func() {
		ws.WriteJSON(&Edu.WsData{Done: true, Data: "操作结束..."})
		ws.Close()
	}()
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		models.Info("eduLoading", err)
		return
	}
	id := this.GetString("id")
	if user, err := models.GetStudentById(id); err == nil {
		if data, ok, cookies, _ := Edu.Sign_in(id, user.EduPwd, "学生"); ok {
			ws.WriteJSON(&Edu.WsData{Done: false, Data: "登录成功..."})
			ws.WriteJSON(&Edu.WsData{Done: false, Data: "正在获取个人信息..."})
			if userinfo, err := Edu.GetStudentProfile(data, cookies); err == nil {
				ws.WriteJSON(&Edu.WsData{Done: false, Data: "读取个人信息成功..."})
				if err = Edu.UpdateStudentProfile(user, userinfo); err != nil {
					ws.WriteJSON(&Edu.WsData{Done: false, Data: "更新个人信息失败..."})
					return
				}
				ws.WriteJSON(&Edu.WsData{Done: false, Data: "更新个人信息成功..."})
				ws.WriteJSON(&Edu.WsData{Done: false, Data: "正在读取个人课表..."})

				//	add edu log
				userinfo_json, err := json.Marshal(userinfo)
				if err == nil {
					Edu.EduLogCreate(&Edu.EduLog{Student: user, Content: string(userinfo_json), Type: "userProfile", Result: 1})
				}

				// get user schedules
				if schedules, err := Edu.GetStudentSchedule(user, data, cookies); err == nil {
					ws.WriteJSON(&Edu.WsData{Done: false, Data: "读取个人课表成功..."})
					ws.WriteJSON(&Edu.WsData{Done: false, Data: "正在读取历史成绩..."})
					if userscore, err := Edu.GetStudentScore(user, data, cookies); err == nil {
						//	通过遍历课程表 与 历史成绩 添加 教师课程表 添加 历史成绩
						if err := Edu.ControlSQLByScheduleAndScore(user, schedules, userscore); err == nil {
							ws.WriteJSON(&Edu.WsData{Done: false, Data: "添加历史成绩成功..."})
							return
						}
						ws.WriteJSON(&Edu.WsData{Done: false, Data: "添加历史成绩失败..."})
						return
					}
					ws.WriteJSON(&Edu.WsData{Done: false, Data: "读取历史成绩失败..."})
					return
				}
				ws.WriteJSON(&Edu.WsData{Done: false, Data: "读取个人课表失败..."})
				return
			}
			ws.WriteJSON(&Edu.WsData{Done: false, Data: "读取个人信息失败..."})
			return
		}
	}
}
