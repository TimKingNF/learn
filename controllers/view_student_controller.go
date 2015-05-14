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
	this.Data["key"] = models.Str2Sha1(this.Ctx.Input.Cookie("beegosessionID"))
	this.Layout = "student/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "student/scripts/signature_scripts.html"
	this.TplNames = "student/setting.html"
}
