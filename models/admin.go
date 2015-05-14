package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Admin struct {
	Id       string `orm:"pk"`
	Name     string `orm:"null;size(50)" valid:"MaxSize(50)"`
	Sex      int    `orm:"default(0)" valid:"Range(0,1,2)"`
	Password string `orm:"null;size(20)" valid:"MaxSize(20)"`
}

func checkAdmin(u *Admin) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkAdmin", err)
		}
	}
	return nil
}

func AdminExist(id string) bool {
	return orm.NewOrm().QueryTable("Admin").Filter("Id", id).Exist()
}

func GetAdminById(id string) (*Admin, error) {
	user := &Admin{Id: id}
	err := orm.NewOrm().Read(user, "Id")
	if err != nil {
		return nil, ErrorInfo("GetAdminById", err)
	}
	return user, nil
}

func UpdateAdmin(ptr *Admin) error {
	_, err := orm.NewOrm().Update(ptr)
	if err != nil {
		return ErrorInfo("UpdateAdmin", err)
	}
	return nil
}
