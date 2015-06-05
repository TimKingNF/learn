package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"time"
)

type StudentNotice struct {
	Id         int64    `orm:"auto"`
	IsRead     int      `orm:"defalut(0)" valid:"Range(0,1)`
	Student    *Student `orm:"rel(fk);null;on_delete(set_null)"`
	SenderType string
	Teacher    *Teacher  `orm:"rel(fk);null;on_delete(set_null)"`
	Admin      *Admin    `orm:"rel(fk);null;on_delete(set_null)"`
	Content    string    `orm:"null;size(2000)" valid:"MaxSize(2000)"`
	CreateTime time.Time `orm:"auto_now;type(datetime)"`
}

func checkStudentNotice(u *StudentNotice) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkStudentNotice", err)
		}
	}
	return nil
}

func AddStudentNotice(Ptr *StudentNotice) error {
	if err := checkStudentNotice(Ptr); err != nil {
		return ErrorInfo("AddStudentNotice", err)
	}
	_, err := orm.NewOrm().Insert(Ptr)
	if err != nil {
		return ErrorInfo("AddStudentNotice", err)
	}
	return nil
}

func GetStudentNoticeById(notice_id int64) (*StudentNotice, error) {
	var notice StudentNotice
	err := orm.NewOrm().QueryTable("StudentNotice").Filter("Id", notice_id).RelatedSel().One(&notice)
	if err != nil {
		return nil, ErrorInfo("GetStudentNoticeById", err)
	}
	return &notice, nil
}

func CountNotReadStudentNotice(s_id string) int64 {
	if len(s_id) <= 0 {
		return 0
	}
	count, err := orm.NewOrm().QueryTable("StudentNotice").Filter("Student", s_id).Filter("IsRead", 0).Count()
	if err != nil {
		return 0
	}
	return count
}

func DelStudentNotice(ptr *StudentNotice) error {
	if ptr == nil {
		return ErrorInfo("DelStudentNotice", "data is nil")
	}
	_, err := orm.NewOrm().Delete(ptr)
	if err != nil {
		return ErrorInfo("DelStudentNotice", err)
	}
	return nil
}

func ReadNoticeByStudent(s_id, sender_type string, sender_id int64) error {
	if len(s_id) <= 0 || len(sender_type) <= 0 || sender_id <= 0 {
		return ErrorInfo("ReadNoticeByStudent", "data is error")
	}
	switch sender_type {
	case "teacher":
		_, err := orm.NewOrm().QueryTable("StudentNotice").Filter("Student", s_id).Filter("SenderType", sender_type).Filter("Teacher", sender_id).Update(orm.Params{
			"IsRead": 1,
		})
		if err != nil {
			return ErrorInfo("ReadNoticeByStudent", err)
		}
		return nil
	case "admin":
		_, err := orm.NewOrm().QueryTable("StudentNotice").Filter("Student", s_id).Filter("SenderType", sender_type).Filter("Admin", sender_id).Update(orm.Params{
			"IsRead": 1,
		})
		if err != nil {
			return ErrorInfo("ReadNoticeByStudent", err)
		}
		return nil
	}
	return ErrorInfo("ReadNoticeByStudent", "sender type is error："+sender_type)
}

func GetStudentNoticeByStudent(s_id string) ([]*StudentNotice, error) {
	if len(s_id) <= 0 {
		return nil, ErrorInfo("GetStudentNoticeByStudent", "s_id is nil")
	}
	var s_notices []*StudentNotice
	_, err := orm.NewOrm().QueryTable("StudentNotice").Filter("Student", s_id).RelatedSel().OrderBy("-CreateTime").All(&s_notices)
	if err != nil {
		return nil, ErrorInfo("GetStudentNoticeByStudent", err)
	}
	return s_notices, nil
}

func RankingNotice(in []*StudentNotice) (info []map[string]interface{}, err error) {
	if in == nil {
		return nil, ErrorInfo("RankingNotice", "data is nil")
	}
	for k, _ := range in {
		var tmp = make(map[string]interface{})
		if in[k].SenderType == "teacher" {
			tmp["sender"] = in[k].Teacher
			tmp["type"] = "teacher"
		} else if in[k].SenderType == "admin" {
			tmp["sender"] = in[k].Admin
			tmp["type"] = "admin"
		}
		tmp["notice"] = make([]*StudentNotice, 0)
		info = append(info, tmp)
	}
	info = FilterRepeat(info).([]map[string]interface{})
	for k, _ := range in {
		for k1, _ := range info {
			if in[k].SenderType == "teacher" && info[k1]["type"] == "teacher" {
				if in[k].Teacher.Id == info[k1]["sender"].(*Teacher).Id {
					info[k1]["notice"] = append(info[k1]["notice"].([]*StudentNotice), in[k])
				}
			} else if in[k].SenderType == "admin" && info[k1]["type"] == "admin" {
				if in[k].Admin.Id == info[k1]["sender"].(*Admin).Id {
					info[k1]["notice"] = append(info[k1]["notice"].([]*StudentNotice), in[k])
				}
			}
		}
	}
	return info, nil
}
