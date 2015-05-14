package controllers

import (
	"learn/models"
)

type ViewAdminController struct {
	ViewController
}

func (this *ViewAdminController) Prepare() {
	//	get session
	user_type := this.GetSession("type").(string)
	if user_type == "教务" {
		id := this.GetSession("id").(string)
		if len(id) > 0 {
			admin, err := models.GetAdminById(id)
			if err == nil {
				this.Data["admin"] = admin
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

// @Title 默认视图
// @router /index [get]
func (this *ViewAdminController) Index() {
	change, _ := this.GetBool("change")
	this.Data["change"] = change
	this.Layout = "admin/base.html"
	this.TplNames = "admin/index.html"
}

// @Title 设置视图
// @router /setting [get]
func (this *ViewAdminController) Setting() {
	this.Data["key"] = models.Str2Sha1(this.Ctx.Input.Cookie("beegosessionID"))
	this.Layout = "admin/base.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Scripts"] = "student/scripts/signature_scripts.html"
	this.TplNames = "admin/setting.html"
}

// @Title 学生视图
// @router /student [get]
func (this *ViewAdminController) Student() {
	this.Layout = "admin/base.html"
	this.TplNames = "admin/student.html"
}

// @Title 数据备份视图
// @router /backup [get]
func (this *ViewAdminController) Backup() {
	this.Layout = "admin/base.html"
	this.TplNames = "admin/backup.html"
}

// @Title 教师视图
// @router /teacher [get]
func (this *ViewAdminController) Teacher() {
	this.Layout = "admin/base.html"
	this.TplNames = "admin/teacher.html"
}

// @Title 课程视图
// @router /course [get]
func (this *ViewAdminController) Course() {
	this.Layout = "admin/base.html"
	this.TplNames = "admin/course.html"
}
