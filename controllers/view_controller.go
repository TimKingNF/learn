package controllers

import (
	// "github.com/astaxie/beego"
	"fmt"
	"learn/models"
)

type ViewController struct {
	baseController
}

// @Title 首页
// @router /index [get]
func (this *ViewController) Index() {
	//	取出 refer
	refer := this.Ctx.Request.Referer()
	if len(refer) > 0 && refer == "http://"+models.Host+":"+models.Port+"/" {
		this.Data["fail"] = true
	}

	this.Data["login_key"] = models.Str2Sha1(this.Ctx.Input.Cookie("beegosessionID"))
	this.TplNames = "login.html"
}

// @Title 错误页面
// @router /error [get]
func (this *ViewController) Error() {
	this.TplNames = "error.html"
}

// @Title 登录验证
// @router /index [post]
func (this *ViewController) Login() {
	//	recevie the login key
	login_key := this.GetString("login_key")
	if len(login_key) <= 0 && login_key != models.Str2Sha1(this.Ctx.Input.Cookie("beegosessionID")) {
		this.Redirect("/index", 302)
		this.StopRun()
	}
	account := this.GetString("account")
	pwd := this.GetString("pwd")
	user_type := this.GetString("type")

	//	植入cookie
	token := fmt.Sprintf("%d", models.GetMathRand(10240))
	this.Ctx.SetCookie("token", token, 3600)
	if this.GetSession("token") != nil {
		this.DelSession("token")
	}
	this.SetSession("token", token)

	switch user_type {
	case "学生":
		if models.StudentExist(account) == true {
			user, err := models.GetStudentById(account)
			if err == nil {
				if user.Password == pwd {
					//	设置session
					if this.GetSession("id") != nil {
						this.DelSession("id")
					}
					this.SetSession("id", account)
					if this.GetSession("type") != nil {
						this.DelSession("type")
					}
					this.SetSession("type", user_type)

					//	login success
					this.Redirect("/view/student/table", 302)
					this.StopRun()
				}
			}
		}
	case "教师":
		if models.TeacherExist(account) == true {
			user, err := models.GetTeacherById(account)
			if err == nil {
				if user.Password == pwd {
					//	设置session
					if this.GetSession("id") != nil {
						this.DelSession("id")
					}
					this.SetSession("id", account)
					if this.GetSession("type") != nil {
						this.DelSession("type")
					}
					this.SetSession("type", user_type)

					//	login success
					this.Redirect("/view/teacher/file", 302)
					this.StopRun()
				}
			}
		}
	case "教务":
		if models.AdminExist(account) == true {
			user, err := models.GetAdminById(account)
			if err == nil {
				if user.Password == pwd {
					//	设置session
					if this.GetSession("id") != nil {
						this.DelSession("id")
					}
					this.SetSession("id", account)
					if this.GetSession("type") != nil {
						this.DelSession("type")
					}
					this.SetSession("type", user_type)

					//	login success
					this.Redirect("/view/admin/index", 302)
					this.StopRun()
				}
			}
		}
	}

	this.Redirect("/", 302)
	this.StopRun()
}

// @Title 退出
// @router /signOut [get]
func (this *ViewController) SignOut() {
	//	删除session
	this.DelSession("id")
	this.DelSession("type")
	this.Redirect("/", 302)
	this.StopRun()
}
