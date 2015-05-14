package models

import (
// "github.com/astaxie/beego/orm"
// "github.com/astaxie/beego/validation"
)

type Department struct {
	Id       int64      `orm:"auto"`
	Name     string     `orm:"null;size(50)" valid:"MaxSize(50)"`
	Students []*Student `orm:"reverse(many)"`
	Teachers []*Teacher `orm:"reverse(many)"`
	Majors   []*Major   `orm:"reverse(many)"`
	Classes  []*Class   `orm:"reverse(many)"`
}
