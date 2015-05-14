package controllers

import (
	"learn/models"
)

type APIAdminController struct {
	baseController
}

func (this *APIAdminController) Prepare() {
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
func (this *APIAdminController) UpdatePassword() {
	//	recevie data
	if user_type := this.GetSession("type").(string); user_type != "教务" {
		this.Redirect("/error", 302)
		this.StopRun()
	}
	id := this.GetSession("id").(string)
	if user, err := models.GetAdminById(id); err == nil {
		newpwd := this.GetString("newpwd")
		chkpwd := this.GetString("chkpwd")
		if newpwd != chkpwd {
			this.Redirect("/error", 302)
			this.StopRun()
		}
		if oldpwd := this.GetString("oldpwd"); user.Password == oldpwd {
			//	update pwd
			user.Password = newpwd
			if err := models.UpdateAdmin(user); err != nil {
				this.Redirect("/error", 302)
				this.StopRun()
			}
			this.Redirect("/view/admin/index?change=true", 302)
			this.StopRun()
		}
	}
	this.Redirect("/error", 302)
	this.StopRun()
}
