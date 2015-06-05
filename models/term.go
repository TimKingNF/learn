package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Term struct {
	Id        int64 `orm:"auto"`
	Number    int   `orm:"null;size(1);default(1)" valid:"Range(1,2)"`
	StartYear int   `orm:"null;size(11)"`
	EndYear   int   `orm:"null;size(11)"`
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
