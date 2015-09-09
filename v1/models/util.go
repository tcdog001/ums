package models

import (
	"github.com/astaxie/beego/orm"
)

func cutLastChar(s string) string {
	c := []byte(s)
	
	Len := len(s)
	if c[Len-1] == '_' {
		c = c[:Len-1]
	}
	
	return string(c)
}

type IEntry interface {
	TableName() string
	KeyName() string
	Key() string
}

func EntryExist(o orm.Ormer, e IEntry) bool {
	return o.QueryTable(e.TableName()).
			Filter(e.KeyName(), e.Key()).
			Exist()
}

func EntryOne(o orm.Ormer, e IEntry, one interface{}) error {
	return o.QueryTable(e.TableName()).
			Filter(e.KeyName(), e.Key()).
			One(one)
}
