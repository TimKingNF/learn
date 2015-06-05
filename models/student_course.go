package models

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"strconv"
	"strings"
)

type StudentCourse struct {
	Id            int64                  `orm:"auto"`
	Student       *Student               `orm:"rel(fk);null;on_delete(set_null)"`
	TeacherCourse *TeacherCourse         `orm:"rel(fk);null;on_delete(set_null)"`
	Score         *StudentGrade          `orm:"reverse(one)"`
	Orgs          map[string]interface{} `orm:"-"`
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

func GetStudentCourseById(id int64) (*StudentCourse, error) {
	if id <= 0 {
		return nil, ErrorInfo("GetStudentCourseById", "id is error")
	}
	var studentCourse StudentCourse
	err := orm.NewOrm().QueryTable("StudentCourse").Filter("Id", id).RelatedSel().One(&studentCourse)
	if err != nil {
		return nil, ErrorInfo("GetStudentCourseById", err)
	}
	return &studentCourse, nil
}

func SearchStudentCourse(id string, t_course_id int64) (*StudentCourse, error) {
	var studentCourse StudentCourse
	err := orm.NewOrm().QueryTable("StudentCourse").Filter("Student", id).Filter("TeacherCourse", t_course_id).RelatedSel().One(&studentCourse)
	if err != nil {
		return nil, ErrorInfo("SearchStudentCourse", err)
	}
	return &studentCourse, nil
}

func GetStudentCourseByTerm(term *Term, id string) ([]*StudentCourse, error) {
	var list []*StudentCourse
	_, err := orm.NewOrm().QueryTable("StudentCourse").Filter("TeacherCourse__Term__Id", term.Id).Filter("Student", id).RelatedSel().All(&list)
	if err != nil {
		return nil, ErrorInfo("GetStudentCourseByTerm", err)
	}
	return list, nil
}

func RangeToMapStudentCourse(in []*StudentCourse) SortMap {
	if len(in) <= 0 {
		return nil
	}
	var info = make(map[int]map[int]map[string]*StudentCourse, 15)
	for _, v := range in {
		var time map[string]string
		if err := json.Unmarshal([]byte(v.TeacherCourse.Time), &time); err != nil {
			continue
		}
		list := strings.Split(time["week_time"], ",")
		for _, v1 := range list {
			week_time, _ := strconv.Atoi(v1)              // 第几节
			week_day, _ := strconv.Atoi(time["week_day"]) // 周几
			if week_time > 0 && week_day > 0 {
				if info[week_time] == nil {
					info[week_time] = make(map[int]map[string]*StudentCourse, 5)
				}
				if info[week_time][week_day] == nil {
					info[week_time][week_day] = make(map[string]*StudentCourse, 3)
				}
				if v.TeacherCourse.Spacing == 1 {
					info[week_time][week_day]["单周"] = v
				} else if v.TeacherCourse.Spacing == 2 {
					info[week_time][week_day]["双周"] = v
				} else {
					info[week_time][week_day]["单,双周"] = v
				}
			}
		}
	}
	for i := 1; i <= 15; i++ {
		var iswait = false
		for k, _ := range info {
			if i == k {
				for i1 := 1; i1 <= 5; i1++ {
					var iswait1 = false
					for k1, _ := range info[k] {
						if i1 == k1 {
							iswait1 = true
							break
						}
					}
					if !iswait1 {
						info[k][i1] = nil
					}
				}
				iswait = true
				break
			}
		}
		if !iswait {
			info[i] = map[int]map[string]*StudentCourse{1: nil, 2: nil, 3: nil, 4: nil, 5: nil}
		}
	}
	var sort = Sort(info)
	for k, v := range sort {
		sort[k] = SortMapItem{v.Key, SortMapItem{v.Key, Sort(v.Val)}}
	}
	return sort
}

func GetStudentCoursesByTeacherCourse(in *TeacherCourse) ([]*StudentCourse, error) {
	if in == nil {
		return nil, ErrorInfo("GetStudentCoursesByTeacherCourse", "data is nil")
	}
	if _, err := orm.NewOrm().LoadRelated(in, "StudentCourses"); err != nil {
		return nil, ErrorInfo("GetStudentCoursesByTeacherCourse", err)
	}
	return in.StudentCourses, nil
}

func FilterStudentCourse(in []*StudentCourse) []*StudentCourse {
	if in == nil {
		return nil
	}
	var list []*StudentCourse
	for k, v := range in {
		if v.Score != nil {
			if v.Score.Grade == 0 && v.Score.GradePointAverage == 0 {
				continue
			} else {
				list = append(list, in[k])
			}
		}
	}
	return list
}
