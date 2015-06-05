package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"time"
)

type Term struct {
	Id            int64                  `orm:"auto"`
	Number        int                    `orm:"null;size(1);default(1)" valid:"Range(1,2)"`
	StartYear     int                    `orm:"null;size(11)"`
	EndYear       int                    `orm:"null;size(11)"`
	TeacherCourse []*TeacherCourse       `orm:"reverse(many)"`
	Orgs          map[string]interface{} `orm:"-"`
}

func checkTerm(u *Term) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkTerm", err)
		}
	}
	return nil
}

func TermExist(num, staryear, endyear int) bool {
	return orm.NewOrm().QueryTable("Term").Filter("Number", num).Filter("StartYear", staryear).Filter("EndYear", endyear).Exist()
}

func AddTerm(Ptr *Term) error {
	if err := checkTerm(Ptr); err != nil {
		return ErrorInfo("AddTerm", err)
	}
	_, err := orm.NewOrm().Insert(Ptr)
	if err != nil {
		return ErrorInfo("AddTerm", err)
	}
	return nil
}

func SearchTerm(num, staryear, endyear int) (*Term, error) {
	var term Term
	err := orm.NewOrm().QueryTable("Term").Filter("Number", num).Filter("StartYear", staryear).Filter("EndYear", endyear).One(&term)
	if err != nil {
		return nil, ErrorInfo("SearchTerm", err)
	}
	return &term, nil
}

func GetTermByTimeNow() (*Term, error) {
	now := time.Now()
	year := now.Year()
	mouth := now.Month()
	if mouth < 9 {
		return SearchTerm(2, year-1, year)
	} else {
		return SearchTerm(1, year, year+1)
	}
	return nil, ErrorInfo("GetTermByTimeNow", errors.New("data is error"))
}

func GetTermListByStudentCourse(id string) ([]*Term, error) {
	var list []*Term
	_, err := orm.NewOrm().QueryTable("Term").Filter("TeacherCourse__StudentCourses__Student__Id", id).RelatedSel().All(&list)
	if err != nil {
		return nil, ErrorInfo("GetTermListByStudentCourse", err)
	}
	//	去除重复数据
	return FilterRepeat(list).([]*Term), nil
}

func GetTermListByTeacherCourse(id int64) ([]*Term, error) {
	var list []*Term
	_, err := orm.NewOrm().QueryTable("Term").Filter("TeacherCourse__Teacher__Id", id).RelatedSel().All(&list)
	if err != nil {
		return nil, ErrorInfo("GetTermListByTeacherCourse", err)
	}
	//	去除重复数据
	return FilterRepeat(list).([]*Term), nil
}

func GetTermListByStudentHomework(id string) ([]*Term, error) {
	var list []*Term
	_, err := orm.NewOrm().QueryTable("Term").Filter("TeacherCourse__StudentCourses__Student__Id", id).RelatedSel().All(&list)
	if err != nil {
		return nil, ErrorInfo("GetTermListByStudentHomework", err)
	}
	//	去除重复数据
	return FilterRepeat(list).([]*Term), nil
}

func RankingTerm(in []*Term) []*Term {
	if in == nil {
		return nil
	}
	for k, _ := range in {
		for k1 := k + 1; k1 < len(in); k1++ {
			if in[k].StartYear > in[k1].StartYear {
				in[k], in[k1] = in[k1], in[k]
			} else if in[k].StartYear == in[k1].StartYear {
				if in[k].Number > in[k1].Number {
					in[k], in[k1] = in[k1], in[k]
				}
			}
		}
	}
	return in
}

func GetTermNumberByTeacher(t_id int64, startYear, endYear int) ([]*Term, error) {
	if t_id == 0 || startYear == 0 || endYear == 0 {
		return nil, ErrorInfo("GetTermNumberByTeacher", "data is zero")
	}
	var terms []*Term
	_, err := orm.NewOrm().QueryTable("Term").Filter("TeacherCourse__Teacher__Id", t_id).Filter("StartYear", startYear).Filter("EndYear", endYear).All(&terms)
	if err != nil {
		return nil, ErrorInfo("GetTermNumberByTeacher", err)
	}
	return FilterRepeat(terms).([]*Term), nil
}

func GetTermNumberByStudent(s_id string, startYear, endYear int) ([]*Term, error) {
	if len(s_id) <= 0 || startYear == 0 || endYear == 0 {
		return nil, ErrorInfo("GetTermNumberByStudent", "data is zero")
	}
	var terms []*Term
	_, err := orm.NewOrm().QueryTable("Term").Filter("TeacherCourse__StudentCourses__Student__Id", s_id).Filter("StartYear", startYear).Filter("EndYear", endYear).All(&terms)
	if err != nil {
		return nil, ErrorInfo("GetTermNumberByStudent", err)
	}
	return FilterRepeat(terms).([]*Term), nil
}
