package controllers

import (
	"github.com/gorilla/websocket"
	"learn/models"
	"net/http"
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
	if _, ok, _, _ := models.EduLogin(id, pwd, "学生"); !ok {
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
		ws.WriteJSON(&models.WsData{Done: true, Data: "操作结束..."})
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
		if data, ok, cookies, _ := models.EduLogin(id, user.EduPwd, "学生"); ok {
			ws.WriteJSON(&models.WsData{Done: false, Data: "登录成功..."})
			ws.WriteJSON(&models.WsData{Done: false, Data: "正在获取个人信息..."})
			if userinfo, err := models.EduGetStudentProfile(data, cookies); err == nil {
				ws.WriteJSON(&models.WsData{Done: false, Data: "读取个人信息成功..."})
				err = models.EduUpdateStudentProfile(user, userinfo)
				if err != nil {
					ws.WriteJSON(&models.WsData{Done: false, Data: "更新个人信息失败..."})
					return
				}
				ws.WriteJSON(&models.WsData{Done: false, Data: "更新个人信息成功..."})
				ws.WriteJSON(&models.WsData{Done: false, Data: "正在读取历史成绩..."})
				if userscore, err := models.EduGetStudentScore(user, data, cookies); err == nil {
					ws.WriteJSON(&models.WsData{Done: false, Data: "读取历史成绩成功..."})
					ws.WriteJSON(&models.WsData{Done: false, Data: "正在读取个人课表..."})
					//	获取当前课表
					if schedules, err := models.EduGetStudentSchedule(user, data, cookies); err == nil {
						ws.WriteJSON(&models.WsData{Done: false, Data: "读取个人课表成功..."})
						ws.WriteJSON(&models.WsData{Done: false, Data: "正在添加历史成绩..."})
						//	通过遍历课程表 与 历史成绩 添加 教师课程表 添加 历史成绩
						if err := models.EduControlSQLByScheduleAndScore(user, schedules, userscore); err == nil {
							ws.WriteJSON(&models.WsData{Done: false, Data: "添加历史成绩成功..."})
							return
						}
						ws.WriteJSON(&models.WsData{Done: false, Data: "添加历史成绩失败..."})
						return
					}
					ws.WriteJSON(&models.WsData{Done: false, Data: "读取个人课表失败..."})
					return
				}
				ws.WriteJSON(&models.WsData{Done: false, Data: "读取历史成绩失败..."})
				return
			}
			ws.WriteJSON(&models.WsData{Done: false, Data: "读取个人信息失败..."})
			return
		}
	}
}
