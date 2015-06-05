package controllers

import (
	// "github.com/astaxie/beego"
	"fmt"
	Edu "learn/EDU"
	"learn/models"
	"strconv"
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
		edu := this.GetString("edu")
		if edu == "on" {
			if !models.StudentExist(account) {
				if data, ok, cookies, _ := Edu.Sign_in(account, pwd, user_type); ok {
					if _, err := Edu.GetStudentProfile(data, cookies); err == nil {
						if err = models.AddStudent(&models.Student{Id: account, EduPwd: pwd}); err == nil {
							//	设置session
							if this.GetSession("id") != nil {
								this.DelSession("id")
							}
							this.SetSession("id", account)
							if this.GetSession("type") != nil {
								this.DelSession("type")
							}
							this.SetSession("type", user_type)

							this.Redirect("/view/student/eduLoading", 302)
							this.StopRun()
						}
					}
				}
			} else {
				user, err := models.GetStudentById(account)
				if err == nil {
					if user.EduPwd == pwd {
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
		}
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
		id, _ := strconv.ParseInt(account, 10, 64)
		if models.TeacherExist(id) == true {
			user, err := models.GetTeacherById(id)
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
					this.Redirect("/view/teacher/table", 302)
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

func (this *ViewController) Finish() {
	this.Render() //	加载模板
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
