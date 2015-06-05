package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type StudentCourse struct {
	Id            int64          `orm:"auto"`
	Student       *Student       `orm:"rel(fk);null;on_delete(set_null)"`
	TeacherCourse *TeacherCourse `orm:"rel(fk);null;on_delete(set_null)"`
}

func checkStudentCourse(u *StudentCourse) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkStudentCourse", err)
		}
	}
	return nil
}

func StudentCourseExist(id string, t_course_id int64) bool {
	return orm.NewOrm().QueryTable("StudentCourse").Filter("Student", id).Filter("TeacherCourse", t_course_id).Exist()
}

func AddStudentCourse(Ptr *StudentCourse) error {
	if err := checkStudentCourse(Ptr); err != nil {
		return ErrorInfo("AddStudentCourse", err)
	}
	_, err := orm.NewOrm().Insert(Ptr)
	if err != nil {
		return ErrorInfo("AddStudentCourse", err)
	}
	return nil
}

func SearchStudentCourse(id string, t_course_id int64) (*StudentCourse, error) {
	var studentCourse StudentCourse
	err := orm.NewOrm().QueryTable("StudentCourse").Filter("Student", id).Filter("TeacherCourse", t_course_id).RelatedSel().One(&studentCourse)
	if err != nil {
		return nil, ErrorInfo("SearchStudentCourse", err)
	}
	return &studentCourse, nil
}
