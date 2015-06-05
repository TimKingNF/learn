package models

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"strconv"
	"strings"
)

type TeacherCourse struct {
	Id             int64                    `orm:"auto"`
	Course         *Course                  `orm:"size(20);rel(fk);null;on_delete(set_null)"`
	Term           *Term                    `orm:"rel(fk);null;on_delete(set_null)"`
	Teacher        *Teacher                 `orm:"rel(fk);null;on_delete(set_null)"`
	Classes        []*Class                 `orm:"rel(m2m)"`
	StudentCourses []*StudentCourse         `orm:"reverse(many)"`
	Place          string                   `orm:"null;size(50)" valid:"MaxSize(50)"`
	Time           string                   `orm:"null;size(50)" valid:"MaxSize(50)"`
	StartWeek      int                      `orm:"default(1)" valid:"Range(1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20)"`
	EndWeek        int                      `orm:"default(20)" valid:"Range(1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20)"`
	Remark         string                   `orm:"size(20);null" valid:"MaxSize(20)"`
	Spacing        int                      `orm:"default(0)" valid:"Range(0,1,2)"`
	Homeworks      []*TeacherCourseHomework `orm:"reverse(many)"`
	StudentChecks  []*StudentCheck          `orm:"reverse(many)"`
	Orgs           map[string]interface{}   `orm:"-"`
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

func TeacherCourseExist(couse_id, time string, term_id, teacher_id int64) bool {
	return orm.NewOrm().QueryTable("TeacherCourse").Filter("Time", time).Filter("Course", couse_id).Filter("Term", term_id).Filter("Teacher", teacher_id).Exist()
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

func SearchTeacherCourse(couse_id, time string, term_id, teacher_id int64) (*TeacherCourse, error) {
	var teacherCourse TeacherCourse
	err := orm.NewOrm().QueryTable("TeacherCourse").Filter("Time", time).Filter("Course", couse_id).Filter("Term", term_id).Filter("Teacher", teacher_id).RelatedSel().One(&teacherCourse)
	if err != nil {
		return nil, ErrorInfo("SearchTeacherCourse", err)
	}
	return &teacherCourse, nil
}

func GetTeacherCourseById(t_course_id int64) (*TeacherCourse, error) {
	var (
		teacherCourse TeacherCourse
		o             = orm.NewOrm()
	)
	err := o.QueryTable("TeacherCourse").Filter("Id", t_course_id).RelatedSel().One(&teacherCourse)
	if err != nil {
		return nil, ErrorInfo("GetTeacherCourseById", err)
	}
	teacherCourse.Orgs = make(map[string]interface{})
	class_num, _ := o.QueryTable("Class").Filter("TeacherCourses__Id", t_course_id).Count()
	teacherCourse.Orgs["class_num"] = class_num
	student_num, _ := o.QueryTable("StudentCourse").Filter("TeacherCourse__Id", t_course_id).Count()
	teacherCourse.Orgs["student_num"] = student_num
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

func GetTeacherCourseByTerm(term *Term, id int64) ([]*TeacherCourse, error) {
	var list []*TeacherCourse
	_, err := orm.NewOrm().QueryTable("TeacherCourse").Filter("Term__Id", term.Id).Filter("Teacher", id).RelatedSel().All(&list)
	if err != nil {
		return nil, ErrorInfo("GetTeacherCourseByTerm", err)
	}
	//	search classes
	for k, v := range list {
		if classes, err := GetClassesByTeacherCourse(v); err == nil {
			list[k].Classes = classes
			if s_courses, err := GetStudentCoursesByTeacherCourse(v); err == nil {
				list[k].StudentCourses = s_courses
			} else {
				continue
			}
		} else {
			continue
		}
	}
	return list, nil
}

func RangeToMapTeacherCourse(in []*TeacherCourse) SortMap {
	if len(in) <= 0 {
		return nil
	}
	var info = make(map[int]map[int]map[string]*TeacherCourse, 15)
	for _, v := range in {
		var time map[string]string
		if err := json.Unmarshal([]byte(v.Time), &time); err != nil {
			continue
		}
		list := strings.Split(time["week_time"], ",")
		for _, v1 := range list {
			week_time, _ := strconv.Atoi(v1)              // 第几节
			week_day, _ := strconv.Atoi(time["week_day"]) // 周几
			if week_time > 0 && week_day > 0 {
				if info[week_time] == nil {
					info[week_time] = make(map[int]map[string]*TeacherCourse, 5)
				}
				if info[week_time][week_day] == nil {
					info[week_time][week_day] = make(map[string]*TeacherCourse, 3)
				}
				if v.Spacing == 1 {
					info[week_time][week_day]["单周"] = v
				} else if v.Spacing == 2 {
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
			info[i] = map[int]map[string]*TeacherCourse{1: nil, 2: nil, 3: nil, 4: nil, 5: nil}
		}
	}
	var sort = Sort(info)
	for k, v := range sort {
		sort[k] = SortMapItem{v.Key, SortMapItem{v.Key, Sort(v.Val)}}
	}
	return sort
}
