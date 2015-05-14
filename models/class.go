package models

import (
// "github.com/astaxie/beego/orm"
// "github.com/astaxie/beego/validation"
)

type Class struct {
	Id         int64       `orm:"auto"`
	Name       string      `orm:"null;size(50)" valid:"MaxSize(50)"`
	Department *Department `orm:"rel(fk);null;on_delete(set_null)"`
	Major      *Major      `orm:"rel(fk);null;on_delete(set_null)"`
	Students   []*Student  `orm:"reverse(many)"`
}
