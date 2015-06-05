package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type StudentHomework struct {
	Id                    int64                  `orm:"auto"`
	Student               *Student               `orm:"rel(fk);null;on_delete(set_null)"`
	Grade                 string                 `orm:"default(E)" valid:"Range(A, B, C, D, E)"`
	Attachment            string                 `orm:"null"`
	IsUpload              int                    `orm:"default(0)" valid:"Range(0,1)"`
	TeacherCourseHomework *TeacherCourseHomework `orm:"rel(fk);null;on_delete(set_null)"`
}

func checkStudentHomework(u *StudentHomework) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkStudentHomework", err)
		}
	}
	return nil
}

func AddStudentHomework(Ptr *StudentHomework) error {
	if err := checkStudentHomework(Ptr); err != nil {
		return ErrorInfo("AddStudentHomework", err)
	}
	_, err := orm.NewOrm().Insert(Ptr)
	if err != nil {
		return ErrorInfo("AddStudentHomework", err)
	}
	return nil
}

func UpdateStudentHomework(Ptr *StudentHomework) error {
	if err := checkStudentHomework(Ptr); err != nil {
		return ErrorInfo("UpdateStudentHomework", err)
	}
	_, err := orm.NewOrm().Update(Ptr)
	if err != nil {
		return ErrorInfo("UpdateStudentHomework", err)
	}
	return nil
}

func GetStudentHomeworkById(id int64) (*StudentHomework, error) {
	var s_homework StudentHomework
	err := orm.NewOrm().QueryTable("StudentHomework").Filter("Id", id).RelatedSel().One(&s_homework)
	if err != nil {
		return nil, ErrorInfo("GetStudentHomeworkById", err)
	}
	return &s_homework, nil
}

func GetStudentHomeworkByTeacherCourseAndStudent(in *TeacherCourse, s_id string) ([]*StudentHomework, error) {
	if in == nil {
		return nil, ErrorInfo("GetStudentHomeworkByTeacherCourseAndStudent", "data is nil")
	}
	var s_homeworks []*StudentHomework
	_, err := orm.NewOrm().QueryTable("StudentHomework").Filter("TeacherCourseHomework__TeacherCourse__Id", in.Id).Filter("Student", s_id).RelatedSel().All(&s_homeworks)
	if err != nil {
		return nil, ErrorInfo("GetStudentHomeworkByTeacherCourseAndStudent", err)
	}
	return s_homeworks, nil
}

func GetStudentHomeworkByStudentCourse(in *StudentCourse, s_id string) ([]*StudentHomework, error) {
	if in == nil {
		return nil, ErrorInfo("GetStudentHomeworkByStudentCourse", "data is nil")
	}
	var s_homeworks []*StudentHomework
	_, err := orm.NewOrm().QueryTable("StudentHomework").Filter("TeacherCourseHomework__TeacherCourse__StudentCourses__Id", in.Id).Filter("Student", s_id).RelatedSel().All(&s_homeworks)
	if err != nil {
		return nil, ErrorInfo("GetStudentHomeworkByStudentCourse", err)
	}
	return s_homeworks, nil
}

func DelStudentHomeworkByTeacherCourseHomework(in *TeacherCourseHomework) error {
	if in == nil {
		return ErrorInfo("DelStudentHomeworkByTeacherCourseHomework", "data is nil")
	}
	_, err := orm.NewOrm().QueryTable("StudentHomework").Filter("TeacherCourseHomework", in.Id).Delete()
	if err != nil {
		return ErrorInfo("DelStudentHomeworkByTeacherCourseHomework", err)
	}
	return nil
}

func GetUploadedStudentHomeworkByTeacherCourseHomework(in *TeacherCourseHomework) ([]*StudentHomework, error) {
	if in == nil {
		return nil, ErrorInfo("GetUploadedStudentHomeworkByTeacherCourseHomework", "data is nil")
	}
	var s_homeworks []*StudentHomework
	_, err := orm.NewOrm().QueryTable("StudentHomework").Filter("TeacherCourseHomework", in.Id).Filter("IsUpload", 1).RelatedSel().All(&s_homeworks)
	if err != nil {
		return nil, ErrorInfo("GetUploadedStudentHomeworkByTeacherCourseHomework", err)
	}
	return s_homeworks, nil
}
