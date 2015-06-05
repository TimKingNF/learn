package EDU

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"learn/models"
	"time"
)

type EduLog struct {
	Id          int64           `orm:"auto"`
	Student     *models.Student `orm:"rel(fk);null;on_delete(set_null)"`
	Content     string          `orm:"null;size(1000)" valid:"MaxSize(1000)"`
	Type        string          `orm:"null;size(20)" valid:"MaxSize(20)`
	Result      int             `orm:"default(0)" valid:"Range(0,1)"`
	ControlTime time.Time       `orm:"auto_now;type(datetime)"`
}

func checkEduLog(u *EduLog) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return models.ErrorInfo("checkEduLog", err)
		}
	}
	return nil
}

func AddEduLog(userPtr *EduLog) error {
	if err := checkEduLog(userPtr); err != nil {
		return models.ErrorInfo("AddEduLog", err)
	}
	_, err := orm.NewOrm().Insert(userPtr)
	if err != nil {
		return models.ErrorInfo("AddEduLog", err)
	}
	return nil
}

func UpdateEduLog(userPtr *EduLog) error {
	if err := checkEduLog(userPtr); err != nil {
		return models.ErrorInfo("UpdateEduLog", err)
	}
	_, err := orm.NewOrm().Update(userPtr)
	if err != nil {
		return models.ErrorInfo("UpdateEduLog", err)
	}
	return nil
}

func EduLogCreate(userPtr *EduLog) error {
	if err := checkEduLog(userPtr); err != nil {
		return models.ErrorInfo("EduLogCreate", err)
	}
	switch userPtr.Type {
	case "userProfile":
		if edu_log, _ := GetUserinfoLog(userPtr.Student.Id); edu_log != nil {
			userPtr.ControlTime = time.Now()
			return models.ErrorInfo("EduLogCreate", UpdateEduLog(userPtr))
		} else {
			return models.ErrorInfo("EduLogCreate", AddEduLog(userPtr))
		}
	default:
		return models.ErrorInfo("EduLogCreate", AddEduLog(userPtr))
	}
	return models.ErrorInfo("EduLogCreate", "no change")
}

func GetUserinfoLog(s_id string) (*EduLog, error) {
	var result EduLog
	err := orm.NewOrm().QueryTable("EduLog").Filter("Student", s_id).Filter("Type", "userProfile").Filter("Result", 1).OrderBy("-ControlTime").Limit(1).One(&result)
	if err != nil {
		return nil, models.ErrorInfo("GetUserinfoLog", err)
	}
	return &result, nil
}

func GetUserDepartmentLog(s_id string) ([]*EduLog, error) {
	var result []*EduLog
	_, err := orm.NewOrm().QueryTable("EduLog").Filter("Student", s_id).Filter("Type", "department").Filter("Result", 1).OrderBy("-ControlTime").All(&result)
	if err != nil {
		return nil, models.ErrorInfo("GetUserDepartmentLog", err)
	}
	return result, nil
}

func GetUserMajorLog(s_id string) ([]*EduLog, error) {
	var result []*EduLog
	_, err := orm.NewOrm().QueryTable("EduLog").Filter("Student", s_id).Filter("Type", "major").Filter("Result", 1).OrderBy("-ControlTime").All(&result)
	if err != nil {
		return nil, models.ErrorInfo("GetUserMajorLog", err)
	}
	return result, nil
}

func GetUserClassLog(s_id string) ([]*EduLog, error) {
	var result []*EduLog
	_, err := orm.NewOrm().QueryTable("EduLog").Filter("Student", s_id).Filter("Type", "class").Filter("Result", 1).OrderBy("-ControlTime").All(&result)
	if err != nil {
		return nil, models.ErrorInfo("GetUserClassLog", err)
	}
	return result, nil
}

func GetUserFailedLog(s_id string) ([]*EduLog, error) {
	var result []*EduLog
	_, err := orm.NewOrm().QueryTable("EduLog").Filter("Student", s_id).Filter("Result", 0).OrderBy("-ControlTime").All(&result)
	if err != nil {
		return nil, models.ErrorInfo("GetUserFailedLog", err)
	}
	return result, nil
}

func GetUserTermLog(s_id string) ([]*EduLog, error) {
	var result []*EduLog
	_, err := orm.NewOrm().QueryTable("EduLog").Filter("Student", s_id).Filter("Type", "term").Filter("Result", 1).OrderBy("-ControlTime").All(&result)
	if err != nil {
		return nil, models.ErrorInfo("GetUserTermLog", err)
	}
	return result, nil
}

func GetUserCourseLog(s_id string) ([]*EduLog, error) {
	var result []*EduLog
	_, err := orm.NewOrm().QueryTable("EduLog").Filter("Student", s_id).Filter("Type", "course").Filter("Result", 1).OrderBy("-ControlTime").All(&result)
	if err != nil {
		return nil, models.ErrorInfo("GetUserCourseLog", err)
	}
	return result, nil
}

func GetUserTeacherLog(s_id string) ([]*EduLog, error) {
	var result []*EduLog
	_, err := orm.NewOrm().QueryTable("EduLog").Filter("Student", s_id).Filter("Type", "teacher").Filter("Result", 1).OrderBy("-ControlTime").All(&result)
	if err != nil {
		return nil, models.ErrorInfo("GetUserTeacherLog", err)
	}
	return result, nil
}

func GetUserTeacherCourseLog(s_id string) ([]*EduLog, error) {
	var result []*EduLog
	_, err := orm.NewOrm().QueryTable("EduLog").Filter("Student", s_id).Filter("Type", "t_course").Filter("Result", 1).OrderBy("-ControlTime").All(&result)
	if err != nil {
		return nil, models.ErrorInfo("GetUserTeacherCourseLog", err)
	}
	return result, nil
}

func GetUserStudentCourseLog(s_id string) ([]*EduLog, error) {
	var result []*EduLog
	_, err := orm.NewOrm().QueryTable("EduLog").Filter("Student", s_id).Filter("Type", "s_course").Filter("Result", 1).OrderBy("-ControlTime").All(&result)
	if err != nil {
		return nil, models.ErrorInfo("GetUserStudentCourseLog", err)
	}
	return result, nil
}

func GetUserTeacherCourseClassLog(s_id string) ([]*EduLog, error) {
	var result []*EduLog
	_, err := orm.NewOrm().QueryTable("EduLog").Filter("Student", s_id).Filter("Type", "t_course_class").Filter("Result", 1).OrderBy("-ControlTime").All(&result)
	if err != nil {
		return nil, models.ErrorInfo("GetUserTeacherCourseClassLog", err)
	}
	return result, nil
}

func GetUserStudentCourseGradeLog(s_id string) ([]*EduLog, error) {
	var result []*EduLog
	_, err := orm.NewOrm().QueryTable("EduLog").Filter("Student", s_id).Filter("Type", "s_course_grade").Filter("Result", 1).OrderBy("-ControlTime").All(&result)
	if err != nil {
		return nil, models.ErrorInfo("GetUserStudentCourseGradeLog", err)
	}
	return result, nil
}
