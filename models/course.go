package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Course struct {
	Id         string      `orm:"pk"`
	Name       string      `orm:"null;size(50)" valid:"MaxSize(50)"`
	Department *Department `orm:"rel(fk);null;on_delete(set_null)"`
	Remark     string      `orm:"null;size(50)" valid:"MaxSize(50)"`
	Type       string      `orm:"null;size(50)" valid:"MaxSize(50)"`
	Credit     float64
}

func checkCourse(u *Course) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkCourse", err)
		}
	}
	return nil
}

func CourseExist(name string) bool {
	return orm.NewOrm().QueryTable("Course").Filter("Name", name).Exist()
}

func GetCourseById(id string) (*Course, error) {
	var course Course
	err := orm.NewOrm().QueryTable("Course").Filter("Id", id).RelatedSel().One(&course)
	if err != nil {
		return nil, ErrorInfo("GetCourseById", err)
	}
	return &course, nil
}

func GetCourseByName(name string) (*Course, error) {
	var course Course
	err := orm.NewOrm().QueryTable("Course").Filter("Name", name).RelatedSel().One(&course)
	if err != nil {
		return nil, ErrorInfo("GetCourseById", err)
	}
	return &course, nil
}

func UpdateCourse(ptr *Course) error {
	_, err := orm.NewOrm().Update(ptr)
	if err != nil {
		return ErrorInfo("UpdateCourse", err)
	}
	return nil
}

func AddCourse(userPtr *Course) error {
	if err := checkCourse(userPtr); err != nil {
		return ErrorInfo("AddCourse", err)
	}
	_, err := orm.NewOrm().Insert(userPtr)
	if err != nil {
		return ErrorInfo("AddCourse", err)
	}
	return nil
}
