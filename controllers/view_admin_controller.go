package controllers

type ViewAdminController struct {
	ViewController
}

// @Title 默认视图
// @router /index [get]
func (this *ViewAdminController) Index() {
	this.Layout = "admin/base.html"
	this.TplNames = "admin/index.html"
}

// @Title 设置视图
// @router /setting [get]
func (this *ViewAdminController) Setting() {
	this.Layout = "admin/base.html"
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
