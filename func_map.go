//	注册模板函数
package main

import (
	"github.com/astaxie/beego"
)

func RegisterFuncMap() {
	beego.AddFuncMap("getSex", func(n int) string {
		switch n {
		case 1:
			return "男"
		case 2:
			return "女"
		default:
			return "未知"
		}
	})
}
