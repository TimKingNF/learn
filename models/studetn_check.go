package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type StudentCheck struct {
	Id            int64          `orm:"auto"`
	Student       *Student       `orm:"rel(fk);null;on_delete(set_null)"`
	Result        string         `orm:"default(未知)" valid:"Range(未知, 已到, 未到, 迟到, 请假)"`
	Week          int            `orm:"default(1)" valid:"Range(1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20)"`
	TeacherCourse *TeacherCourse `orm:"rel(fk);null;on_delete(set_null)"`
}

func checkStudentCheck(u *StudentCheck) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkStudentCheck", err)
		}
	}
	return nil
}

func AddStudentCheck(Ptr *StudentCheck) error {
	if err := checkStudentCheck(Ptr); err != nil {
		return ErrorInfo("AddStudentCheck", err)
	}
	_, err := orm.NewOrm().Insert(Ptr)
	if err != nil {
		return ErrorInfo("AddStudentCheck", err)
	}
	return nil
}

func ExistStudentCheck(ptr *StudentCheck) bool {
	if ptr == nil {
		return false
	}
	return orm.NewOrm().QueryTable("StudentCheck").Filter("TeacherCourse", ptr.TeacherCourse.Id).Filter("Student", ptr.Student.Id).Filter("Week", ptr.Week).Exist()
}

func SearchStudentCheck(t_course_id int64, s_id string, week int) (*StudentCheck, error) {
	if t_course_id <= 0 || len(s_id) <= 0 || week <= 0 {
		return nil, ErrorInfo("SearchStudentCheck", "data is error")
	}
	var s_check StudentCheck
	err := orm.NewOrm().QueryTable("StudentCheck").Filter("TeacherCourse", t_course_id).Filter("Student", s_id).Filter("Week", week).RelatedSel().One(&s_check)
	if err != nil {
		return nil, ErrorInfo("SearchStudentCheck", err)
	}
	return &s_check, nil
}

func UpdateStudentCheck(Ptr *StudentCheck) error {
	if err := checkStudentCheck(Ptr); err != nil {
		return ErrorInfo("UpdateStudentCheck", err)
	}
	_, err := orm.NewOrm().Update(Ptr)
	if err != nil {
		return ErrorInfo("UpdateStudentCheck", err)
	}
	return nil
}

func GetStudentChecksByTeacherCourseAndWeek(t_course *TeacherCourse, week int) ([]*StudentCheck, error) {
	if t_course == nil {
		return nil, ErrorInfo("GetStudentChecksByTeacherCourseAndWeek", "data is nil")
	}
	var s_checks []*StudentCheck
	_, err := orm.NewOrm().QueryTable("StudentCheck").Filter("TeacherCourse", t_course.Id).Filter("Week", week).RelatedSel().All(&s_checks)
	if err != nil {
		return nil, ErrorInfo("GetStudentChecksByTeacherCourseAndWeek", err)
	}
	return s_checks, nil
}

func GetStudentChecksByTeacherCourseAndStudent(t_course *TeacherCourse, s_id string) ([]*StudentCheck, error) {
	if t_course == nil {
		return nil, ErrorInfo("GetStudentChecksByTeacherCourseAndStudent", "data is nil")
	}
	var s_checks []*StudentCheck
	_, err := orm.NewOrm().QueryTable("StudentCheck").Filter("TeacherCourse", t_course.Id).Filter("Student", s_id).RelatedSel().All(&s_checks)
	if err != nil {
		return nil, ErrorInfo("GetStudentChecksByTeacherCourseAndStudent", err)
	}
	return s_checks, nil
}

func GetStudentChecksByTeacherCourse(t_course *TeacherCourse) ([]*StudentCheck, error) {
	if t_course == nil {
		return nil, ErrorInfo("GetStudentChecksByTeacherCourse", "data is nil")
	}
	var s_checks []*StudentCheck
	_, err := orm.NewOrm().QueryTable("StudentCheck").Filter("TeacherCourse", t_course.Id).RelatedSel().All(&s_checks)
	if err != nil {
		return nil, ErrorInfo("GetStudentChecksByTeacherCourse", err)
	}
	return s_checks, nil
}
