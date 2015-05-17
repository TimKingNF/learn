package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Major struct {
	Id         int64       `orm:"auto"`
	Name       string      `orm:"null;size(50)" valid:"MaxSize(50)"`
	Department *Department `orm:"rel(fk);null;on_delete(set_null)"`
	Classes    []*Class    `orm:"reverse(many)"`
	Students   []*Student  `orm:"reverse(many)"`
}

func checkMajor(u *Major) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkMajor", err)
		}
	}
	return nil
}

func MajorExist(name string) bool {
	return orm.NewOrm().QueryTable("Major").Filter("Name", name).Exist()
}

func GetMajorById(id int64) (*Major, error) {
	Major := &Major{Id: id}
	err := orm.NewOrm().Read(Major, "Id")
	if err != nil {
		return nil, ErrorInfo("GetMajorById", err)
	}
	return Major, nil
}

func GetMajorByName(name string) (*Major, error) {
	var major Major
	err := orm.NewOrm().QueryTable("Major").Filter("Name", name).One(&major)
	if err != nil {
		return nil, ErrorInfo("GetMajorByName", err)
	}
	return &major, nil
}

func UpdateMajor(ptr *Major) error {
	_, err := orm.NewOrm().Update(ptr)
	if err != nil {
		return ErrorInfo("UpdateMajor", err)
	}
	return nil
}

func AddMajor(ptr *Major) (int64, error) {
	if err := checkMajor(ptr); err != nil {
		return 0, ErrorInfo("AddMajor", err)
	}
	id, err := orm.NewOrm().Insert(ptr)
	if err != nil {
		return 0, ErrorInfo("AddMajor", err)
	}
	return id, nil
}
