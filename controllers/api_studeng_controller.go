package controllers

import (
	"learn/models"
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
