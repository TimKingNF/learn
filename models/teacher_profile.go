package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type TeacherProfile struct {
	Id      int64    `orm:"auto"`
	Teacher *Teacher `orm:"rel(one);null;on_delete(set_null)"`
	Phone   string   `orm:"null;size(20)" valid:"Phone;Maxsize(20)"`
	Email   string   `orm:"null;size(40)" valid:"Email;Maxsize(40)"`
	Degree  string   `orm:"null;size(20)"`  // 学历 学位
	Title   string   `orm:"null;size(40)"`  // 职称
	Subject string   `orm:"null;size(40)"`  // 学科方向
	Remark  string   `orm:"null;size(600)"` // 300字的简介
}

func checkTeacherProfile(u *TeacherProfile) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkTeacherProfile", err)
		}
	}
	return nil
}

func TeacherProfileExist(id int64) bool {
	return orm.NewOrm().QueryTable("TeacherProfile").Filter("Teacher", id).Exist()
}

func AddTeacherProfile(userPtr *TeacherProfile) error {
	if err := checkTeacherProfile(userPtr); err != nil {
		return ErrorInfo("AddTeacherProfile", err)
	}
	_, err := orm.NewOrm().Insert(userPtr)
	if err != nil {
		return ErrorInfo("AddTeacherProfile", err)
	}
	return nil
}

func GetTeacherProfile(t_id int64) (*TeacherProfile, error) {
	var t_profile TeacherProfile
	err := orm.NewOrm().QueryTable("TeacherProfile").Filter("Teacher", t_id).One(&t_profile)
	if err != nil {
		return nil, ErrorInfo("GetTeacherProfile", err)
	}
	return &t_profile, nil
}

func UpdateTeacherProfile(in *Teacher) (err error) {
	if err := checkTeacherProfile(in.Profile); err != nil {
		return ErrorInfo("UpdateTeacherProfile", err)
	}
	//	teacher profile exist
	if exist := TeacherProfileExist(in.Id); !exist {
		//	add teacher profile
		if err := AddTeacherProfile(in.Profile); err != nil {
			return ErrorInfo("UpdateTeacherProfile", err)
		}
	} else {
		//	search teacher profile
		if t_profile, err := GetTeacherProfile(in.Id); err == nil {
			//	update teacher profile
			in.Profile.Id = t_profile.Id
			if _, err := orm.NewOrm().Update(in.Profile); err != nil {
				return ErrorInfo("UpdateTeacherProfile", err)
			}
		}
	}
	//	update teacher
	return UpdateTeacher(in)
}
