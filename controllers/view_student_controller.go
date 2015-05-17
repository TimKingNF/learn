package controllers

import (
	"learn/models"
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
	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "student/scripts/notice_scripts.html"
	this.LayoutSections["Head_html"] = "student/head/notice_head.html"
	this.TplNames = "student/notice.html"
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

	this.Layout = "student/base.html"
	this.TplNames = "student/eduManage.html"
}
