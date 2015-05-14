package controllers

type ViewTeatherController struct {
	ViewController
}

// @Title 文件视图
// @router /file [get]
func (this *ViewTeatherController) File() {
	this.TplNames = "teacher/file.html"
}
