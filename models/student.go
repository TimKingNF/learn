package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Student struct {
	Id              string             `orm:"pk"`
	Name            string             `orm:"null;size(50)" valid:"MaxSize(50)"`
	Sex             int                `orm:"default(0)" valid:"Range(0,1,2)"`
	Password        string             `orm:"null;size(20)" valid:"MaxSize(20)"`
	EduPwd          string             `orm:"null;size(20)" valid:"MaxSize(20)"`
	Headimgurl      string             `orm:"default(/static/img/avatar.jpeg)"`
	EnterSchoolYear int                `orm:"null;size(4)" valid:"MaxSize(4)"`
	EnterSchoolDate int                `orm:"null;size(8)" valid:"MaxSize(8)"`
	Department      *Department        `orm:"rel(fk);null;on_delete(set_null)"`
	Major           *Major             `orm:"rel(fk);null;on_delete(set_null)"`
	Class           *Class             `orm:"rel(fk);null;on_delete(set_null)"`
	IdCard          int64              `orm:"null;size(20)" valid:"MaxSize(20)"`
	IsEdu           bool               `orm:"default(0)"`
	StudentCourses  []*StudentCourse   `orm:"reverse(many)"`
	Homeworks       []*StudentHomework `orm:"reverse(many)"`
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
	var user Student
	err := orm.NewOrm().QueryTable("Student").Filter("Id", id).RelatedSel().One(&user)
	if err != nil {
		return nil, ErrorInfo("GetStudentById", err)
	}
	return &user, nil
}

func UpdateStudent(ptr *Student) error {
	_, err := orm.NewOrm().Update(ptr)
	if err != nil {
		return ErrorInfo("UpdateStudent", err)
	}
	return nil
}

func AddStudent(userPtr *Student) error {
	if err := checkStudent(userPtr); err != nil {
		return ErrorInfo("AddStudent", err)
	}
	_, err := orm.NewOrm().Insert(userPtr)
	if err != nil {
		return ErrorInfo("AddStudent", err)
	}
	return nil
}

func GetStudentsByTeacherCourse(in *TeacherCourse) ([]*Student, error) {
	if in == nil {
		return nil, ErrorInfo("GetStudentsByTeacherCourse", "data is nil")
	}
	var students []*Student
	if _, err := orm.NewOrm().QueryTable("Student").Filter("StudentCourses__TeacherCourse__Id", in.Id).RelatedSel().All(&students); err != nil {
		return nil, ErrorInfo("GetStudentsByTeacherCourse", err)
	}
	return students, nil
}

func GetStudentsByClass(in *Class) ([]*Student, error) {
	if in == nil {
		return nil, ErrorInfo("GetStudentsByClass", "data is nil")
	}
	var students []*Student
	if _, err := orm.NewOrm().QueryTable("Student").Filter("Class", in.Id).All(&students); err != nil {
		return nil, ErrorInfo("GetStudentsByClass", err)
	}
	return students, nil
}
