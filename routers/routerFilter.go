package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"learn/models"
	"strings"
)

var noNeedSessionRouter = [...]string{
	"/",
	"/error",
}

func CORSSupport(ctx *context.Context) {
	ctx.Output.Header("Access-Control-Allow-Origin", "*")
}

func isNeedSession(url string) bool {
	for _, r := range noNeedSessionRouter {
		if strings.Contains(strings.ToLower(url), strings.ToLower(r)) {
			return false
		}
	}
	return true
}

func isConfrimView(user_type, url string) bool {
	switch user_type {
	case "学生":
		if strings.Contains(strings.ToLower(url), "student") {
			return true
		}
	case "教师":
		if strings.Contains(strings.ToLower(url), "teacher") {
			return true
		}
	case "教务":
		if strings.Contains(strings.ToLower(url), "admin") {
			return true
		}
	}
	return false
}

func checkLogin(ctx *context.Context) {
	_, ok := ctx.Input.Session("id").(string)
	if !ok && !isNeedSession(ctx.Request.RequestURI) {
		ctx.Output.Json(models.GetFailedResponse(1, "您的会话已过期，请重新登录"), true, false)
		return
	}
	user_type, ok := ctx.Input.Session("type").(string)
	if !ok || !isConfrimView(user_type, ctx.Request.RequestURI) {
		ctx.Output.Json(models.GetFailedResponse(1, "您的会话已过期，请重新登录"), true, false)
		return
	}
}

func AddRouterFilter() {
	beego.InsertFilter("*", beego.BeforeRouter, CORSSupport)
	//	添加路由限制
	beego.InsertFilter("/view/*", beego.BeforeRouter, checkLogin)
}
