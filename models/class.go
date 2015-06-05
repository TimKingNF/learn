package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Class struct {
	Id             int64            `orm:"auto"`
	Name           string           `orm:"null;size(50)" valid:"MaxSize(50)"`
	Department     *Department      `orm:"rel(fk);null;on_delete(set_null)"`
	Major          *Major           `orm:"rel(fk);null;on_delete(set_null)"`
	Students       []*Student       `orm:"reverse(many)"`
	TeacherCourses []*TeacherCourse `orm:"reverse(many)"`
}

func checkClass(u *Class) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkClass", err)
		}
	}
	return nil
}

func ClassExist(name string) bool {
	return orm.NewOrm().QueryTable("Class").Filter("Name", name).Exist()
}

func GetClassById(id int64) (*Class, error) {
	var class Class
	err := orm.NewOrm().QueryTable("Class").Filter("Id", id).RelatedSel().One(&class)
	if err != nil {
		return nil, ErrorInfo("GetClassById", err)
	}
	return &class, nil
}

func GetClassByName(name string) (*Class, error) {
	var class Class
	err := orm.NewOrm().QueryTable("Class").Filter("Name", name).One(&class)
	if err != nil {
		return nil, ErrorInfo("GetClassByName", err)
	}
	return &class, nil
}

func UpdateClass(ptr *Class) error {
	_, err := orm.NewOrm().Update(ptr)
	if err != nil {
		return ErrorInfo("UpdateClass", err)
	}
	return nil
}

func AddClass(ptr *Class) (int64, error) {
	if err := checkClass(ptr); err != nil {
		return 0, ErrorInfo("AddClass", err)
	}
	id, err := orm.NewOrm().Insert(ptr)
	if err != nil {
		return 0, ErrorInfo("AddClass", err)
	}
	return id, nil
}

func GetClassesByTeacherCourse(in *TeacherCourse) ([]*Class, error) {
	if in == nil {
		return nil, ErrorInfo("GetClassesByTeacherCourse", "data is nil")
	}
	if _, err := orm.NewOrm().LoadRelated(in, "Classes"); err != nil {
		return nil, ErrorInfo("GetClassesByTeacherCourse", err)
	}
	return in.Classes, nil
}
