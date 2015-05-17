package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Department struct {
	Id       int64      `orm:"auto"`
	Name     string     `orm:"null;size(50)" valid:"MaxSize(50)"`
	Students []*Student `orm:"reverse(many)"`
	Teachers []*Teacher `orm:"reverse(many)"`
	Majors   []*Major   `orm:"reverse(many)"`
	Classes  []*Class   `orm:"reverse(many)"`
	Courses  []*Course  `orm:"reverse(many)"`
}

func checkDepartment(u *Department) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(u)
	if !b {
		for _, err := range valid.Errors {
			return ErrorInfo("checkDepartment", err)
		}
	}
	return nil
}

func DepartmentExist(name string) bool {
	return orm.NewOrm().QueryTable("Department").Filter("Name", name).Exist()
}

func GetDepartmentById(id int64) (*Department, error) {
	department := &Department{Id: id}
	err := orm.NewOrm().Read(department, "Id")
	if err != nil {
		return nil, ErrorInfo("GetDepartmentById", err)
	}
	return department, nil
}

func GetDepartmentByName(name string) (*Department, error) {
	var department Department
	err := orm.NewOrm().QueryTable("Department").Filter("Name", name).One(&department)
	if err != nil {
		return nil, ErrorInfo("GetDepartmentByName", err)
	}
	return &department, nil
}

func UpdateDepartment(ptr *Department) error {
	_, err := orm.NewOrm().Update(ptr)
	if err != nil {
		return ErrorInfo("UpdateDepartment", err)
	}
	return nil
}

func AddDepartment(ptr *Department) (int64, error) {
	if err := checkDepartment(ptr); err != nil {
		return 0, ErrorInfo("AddDepartment", err)
	}
	id, err := orm.NewOrm().Insert(ptr)
	if err != nil {
		return 0, ErrorInfo("AddDepartment", err)
	}
	return id, nil
}
