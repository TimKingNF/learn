package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type StudentGrade struct {
	Id                int64          `orm:"auto"`
	StudentCourse     *StudentCourse `orm:"rel(one);null;on_delete(set_null)"`
	Grade             int
	GradePointAverage float64
}

func checkStudentGrade(u *StudentGrade) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkStudentGrade", err)
		}
	}
	return nil
}

func StudentGradeExist(s_course_id int64) bool {
	return orm.NewOrm().QueryTable("StudentGrade").Filter("StudentCourse", s_course_id).Exist()
}

func AddStudentGrade(Ptr *StudentGrade) error {
	if err := checkStudentGrade(Ptr); err != nil {
		return ErrorInfo("AddStudentGrade", err)
	}
	_, err := orm.NewOrm().Insert(Ptr)
	if err != nil {
		return ErrorInfo("AddStudentGrade", err)
	}
	return nil
}

func UpdateStudentGrade(Ptr *StudentGrade) error {
	if err := checkStudentGrade(Ptr); err != nil {
		return ErrorInfo("UpdateStudentGrade", err)
	}
	_, err := orm.NewOrm().Update(Ptr)
	if err != nil {
		return ErrorInfo("UpdateStudentGrade", err)
	}
	return nil
}

func SearchStudentGradeByStudentCourse(in *StudentCourse) ([]*StudentGrade, error) {
	var student_grades []*StudentGrade
	_, err := orm.NewOrm().QueryTable("StudentGrade").Filter("StudentCourse", in.Id).All(&student_grades)
	if err != nil {
		return nil, ErrorInfo("SearchStudentGradeByStudentCourse", err)
	}
	return student_grades, nil
}

func SearchStudentGrade(s_course_id int64) (*StudentGrade, error) {
	var studentGrade StudentGrade
	err := orm.NewOrm().QueryTable("StudentGrade").Filter("StudentCourse", s_course_id).RelatedSel().One(&studentGrade)
	if err != nil {
		return nil, ErrorInfo("SearchStudentGrade", err)
	}
	return &studentGrade, nil
}
