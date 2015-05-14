package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Teacher struct {
	Id         string      `orm:"pk"`
	Name       string      `orm:"null;size(50)" valid:"MaxSize(50)"`
	Sex        int         `orm:"default(0)" valid:"Range(0,1,2)"`
	Password   string      `orm:"null;size(20)" valid:"MaxSize(20)"`
	Headimgurl string      `orm:"default(/static/img/avatar.jpeg)"`
	Department *Department `orm:"rel(fk);null;on_delete(set_null)"`
}

func checkTeacher(u *Teacher) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkTeacher", err)
		}
	}
	return nil
}

func TeacherExist(id string) bool {
	return orm.NewOrm().QueryTable("Teacher").Filter("Id", id).Exist()
}

func GetTeacherById(id string) (*Teacher, error) {
	user := &Teacher{Id: id}
	err := orm.NewOrm().Read(user, "Id")
	if err != nil {
		return nil, ErrorInfo("GetTeacherById", err)
	}
	return user, nil
}
