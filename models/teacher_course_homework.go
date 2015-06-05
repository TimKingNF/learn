package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"time"
)

type TeacherCourseHomework struct {
	Id                  int64          `orm:"auto"`
	TeacherCourse       *TeacherCourse `orm:"rel(fk);null;on_delete(set_null)"`
	Title               string         `orm:"null;size(30)" valid:"MaxSize(30)"`                                            //	作业标题
	Remark              string         `orm:"null;size(150)" valid:"MaxSize(150)"`                                          //	作业描述
	AsOfTime            string         `orm:"null"`                                                                         // 作业截止时间
	Attachment          string         `orm:"null"`                                                                         //	附件地址
	PublishWeek         int            `orm:"default(1)" valid:"Range(1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20)"` // 发布周
	StudentHomeworkPath string
	CreateTime          time.Time          `orm:"auto_now;type(datetime)"`
	StudentHomeworks    []*StudentHomework `orm:"reverse(many)"`
}

func checkTeacherCourseHomework(u *TeacherCourseHomework) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkTeacherCourseHomework", err)
		}
	}
	return nil
}

func ExistTeacherCourseHomework(ptr *TeacherCourseHomework) bool {
	return orm.NewOrm().QueryTable("TeacherCourseHomework").Filter("TeacherCourse", ptr.TeacherCourse.Id).Filter("PublishWeek", ptr.PublishWeek).Exist()
}

func AddTeacherCourseHomework(Ptr *TeacherCourseHomework) (int64, error) {
	if err := checkTeacherCourseHomework(Ptr); err != nil {
		return 0, ErrorInfo("AddTeacherCourseHomework", err)
	}
	id, err := orm.NewOrm().Insert(Ptr)
	if err != nil {
		return 0, ErrorInfo("AddTeacherCourseHomework", err)
	}
	return id, nil
}

func UpdateTeacherCourseHomework(Ptr *TeacherCourseHomework) error {
	if err := checkTeacherCourseHomework(Ptr); err != nil {
		return ErrorInfo("UpdateTeacherCourseHomework", err)
	}
	_, err := orm.NewOrm().Update(Ptr)
	if err != nil {
		return ErrorInfo("UpdateTeacherCourseHomework", err)
	}
	return nil
}

func GetTeacherCourseHomeworkByTeacherCourse(in *TeacherCourse) ([]*TeacherCourseHomework, error) {
	if in == nil {
		return nil, ErrorInfo("GetTeacherCourseHomeworkByTeacherCourse", "data is nil")
	}
	if _, err := orm.NewOrm().LoadRelated(in, "Homeworks"); err != nil {
		return nil, ErrorInfo("GetTeacherCourseHomeworkByTeacherCourse", err)
	}
	return in.Homeworks, nil
}

func GetTeacherCourseHomeworkById(in int64) (*TeacherCourseHomework, error) {
	var t_course_homework TeacherCourseHomework
	if err := orm.NewOrm().QueryTable("TeacherCourseHomework").Filter("Id", in).RelatedSel().One(&t_course_homework); err != nil {
		return nil, ErrorInfo("GetTeacherCourseHomeworkById", err)
	}
	return &t_course_homework, nil
}

func DelTeacherCourseHomework(ptr *TeacherCourseHomework) error {
	if ptr == nil {
		return ErrorInfo("DelTeacherCourseHomework", "data is nil")
	}
	_, err := orm.NewOrm().Delete(ptr)
	if err != nil {
		return ErrorInfo("DelTeacherCourseHomework", err)
	}
	return nil
}
