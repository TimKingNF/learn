package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type TeacherCourse struct {
	Id        int64    `orm:"auto"`
	Course    *Course  `orm:"size(20);rel(fk);null;on_delete(set_null)"`
	Term      *Term    `orm:"rel(fk);null;on_delete(set_null)"`
	Teacher   *Teacher `orm:"rel(fk);null;on_delete(set_null)"`
	Classes   []*Class `orm:"rel(m2m)"`
	Place     string   `orm:"null;size(50)" valid:"MaxSize(50)"`
	Time      string   `orm:"null;size(50)" valid:"MaxSize(50)"`
	StartWeek int      `orm:"default(1)" valid:"Range(1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20)"`
	EndWeek   int      `orm:"default(20)" valid:"Range(1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20)"`
}

func checkTeacherCourse(u *TeacherCourse) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkTeacherCourse", err)
		}
	}
	return nil
}

func TeacherCourseExist(couse_id string, term_id, teacher_id int64) bool {
	return orm.NewOrm().QueryTable("TeacherCourse").Filter("Course", couse_id).Filter("Term", term_id).Filter("Teacher", teacher_id).Exist()
}

func AddTeacherCourse(Ptr *TeacherCourse) error {
	if err := checkTeacherCourse(Ptr); err != nil {
		return ErrorInfo("AddTeacherCourse", err)
	}
	_, err := orm.NewOrm().Insert(Ptr)
	if err != nil {
		return ErrorInfo("AddTeacherCourse", err)
	}
	return nil
}

func SearchTeacherCourse(couse_id string, term_id, teacher_id int64) (*TeacherCourse, error) {
	var teacherCourse TeacherCourse
	err := orm.NewOrm().QueryTable("TeacherCourse").Filter("Course", couse_id).Filter("Term", term_id).Filter("Teacher", teacher_id).RelatedSel().One(&teacherCourse)
	if err != nil {
		return nil, ErrorInfo("SearchTeacherCourse", err)
	}
	return &teacherCourse, nil
}

func ExistTeacherCourseAndClass(tc_id, c_id int64) bool {
	m2m := orm.NewOrm().QueryM2M(&TeacherCourse{Id: tc_id}, "Classes")
	return m2m.Exist(&Class{Id: c_id})
}

func AddM2MBetweenTeacherCourseAndClass(tc_id, c_id int64) error {
	m2m := orm.NewOrm().QueryM2M(&TeacherCourse{Id: tc_id}, "Classes")
	_, err := m2m.Add(&Class{Id: c_id})
	if err != nil {
		return ErrorInfo("AddM2MBetweenTeacherCourseAndClasstc_id", err)
	}
	return nil
}
