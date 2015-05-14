package controllers

import (
	"fmt"
)

type ViewStudentController struct {
	ViewController
}

// @Title 课程表视图
// @router /table [get]
func (this *ViewStudentController) Table() {
	fmt.Println(this.GetSession("id"))
	this.TplNames = "student/table.html"
}
