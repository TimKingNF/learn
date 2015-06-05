package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Teacher struct {
	Id         int64       `orm:"auto"`
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

func TeacherExistByDep(name string, dep_id int64) bool {
	return orm.NewOrm().QueryTable("Teacher").Filter("Name", name).Filter("Department", dep_id).Exist()
}

func TeacherExist(id int64) bool {
	return orm.NewOrm().QueryTable("Teacher").Filter("Id", id).Exist()
}

func GetTeacherById(id int64) (*Teacher, error) {
	var user Teacher
	err := orm.NewOrm().QueryTable("Teacher").Filter("Id", id).RelatedSel().One(&user)
	if err != nil {
		return nil, ErrorInfo("GetTeacherById", err)
	}
	return &user, nil
}

func AddTeacher(userPtr *Teacher) error {
	if err := checkTeacher(userPtr); err != nil {
		return ErrorInfo("AddTeacher", err)
	}
	_, err := orm.NewOrm().Insert(userPtr)
	if err != nil {
		return ErrorInfo("AddTeacher", err)
	}
	return nil
}

func SearchTeacher(name string, dep_id int64) (*Teacher, error) {
	var user Teacher
	err := orm.NewOrm().QueryTable("Teacher").Filter("Department", dep_id).Filter("Name", name).RelatedSel().One(&user)
	if err != nil {
		return nil, ErrorInfo("SearchTeacher", err)
	}
	return &user, nil
}
