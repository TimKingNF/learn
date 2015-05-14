//	学风跟踪系统
//	version -- 0.1
package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"learn/models"
	_ "learn/routers"
)

func init() {
	orm.Debug = false                         //	数据库调试模式
	orm.RegisterDriver("mysql", orm.DR_MySQL) // 注册数据库驱动 MYSQL
	//	注册数据库
	err := orm.RegisterDataBase("default", "mysql", "root:timking@tcp(127.0.0.1:13306)/learn?charset=utf8&loc=Asia%2FShanghai")
	if err != nil {
		models.Info("main", err)
	}
	//	注册数据库表模型
	orm.RegisterModel(new(models.Student), new(models.Department), new(models.Major), new(models.Class), new(models.Teacher), new(models.Admin))
	err = orm.RunSyncdb("default", false, false) //	运行数据库
	if err != nil {
		models.Info("main", err)
	}
}

func main() {
	//	注册文件目录
	beego.SetStaticPath("UPLOADS", "uploads")

	//	启动session
	beego.SessionOn = true

	//	注册模板函数
	RegisterFuncMap()
	beego.Run()
}
