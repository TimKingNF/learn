package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Student struct {
	Id              string      `orm:"pk"`
	Name            string      `orm:"null;size(50)" valid:"MaxSize(50)"`
	Sex             int         `orm:"default(0)" valid:"Range(0,1,2)"`
	Password        string      `orm:"null;size(20)" valid:"MaxSize(20)"`
	Headimgurl      string      `orm:"default(/static/img/avatar.jpeg)"`
	EnterSchoolYear int         `orm:"null;size(4)" valid:"MaxSize(4)"`
	EnterSchoolDate int         `orm:"null;size(8)" valid:"MaxSize(8)"`
	Department      *Department `orm:"rel(fk);null;on_delete(set_null)"`
	Major           *Major      `orm:"rel(fk);null;on_delete(set_null)"`
	Class           *Class      `orm:"rel(fk);null;on_delete(set_null)"`
	IdCard          int64       `orm:"null;size(20)" valid:"MaxSize(20)"`
}

func checkStudent(u *Student) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkStudent", err)
		}
	}
	return nil
}

func StudentExist(id string) bool {
	return orm.NewOrm().QueryTable("Student").Filter("Id", id).Exist()
}

func GetStudentById(id string) (*Student, error) {
	user := &Student{Id: id}
	err := orm.NewOrm().Read(user, "Id")
	if err != nil {
		return nil, ErrorInfo("GetStudentById", err)
	}
	return user, nil
}

func UpdateStudent(ptr *Student) error {
	_, err := orm.NewOrm().Update(ptr)
	if err != nil {
		return ErrorInfo("UpdateStudent", err)
	}
	return nil
}
