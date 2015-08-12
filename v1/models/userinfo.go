package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Userinfo struct {
	Username         string    `orm:"pk";json:"username"`
	LastRegisterTime time.Time `json:"-"`
}

func (account *Userinfo) TableName() string {
	return "userinfo"
}

func RegisterUserinfo(account *Userinfo) bool {
	o := orm.NewOrm()
	account.LastRegisterTime = time.Now()
	beego.Debug("regiteraccount table=", account.TableName())
	//查找对应的username是否存在
	exist := o.QueryTable(account.TableName()).Filter("username", account.Username).Exist()
	if exist {
		//account存在，则更新account信息
		//return UpdateUserinfo(account)
		return true
	} else {
		//account不存在，则插入account信息
		_, err := o.Insert(account)
		if err != nil {
			beego.Error(err)
			return false
		}
		return true
	}
}

func UpdateUserinfo(account *Userinfo) bool {
	beego.Debug("UpdateUserinfo table=", account.TableName())
	o := orm.NewOrm()

	var acc Userinfo
	err := o.QueryTable(account.TableName()).Filter("username", account.Username).One(&acc)
	if err != nil {
		return false
	} else {
		acc.Username = account.Username
		acc.LastRegisterTime = account.LastRegisterTime

		beego.Debug("UpdateUserinfo Username =", acc.Username)

		_, err := o.Update(&acc)
		if err != nil {
			beego.Error(err)
			return false
		}
		return true
	}
}
func (this *Userinfo) Init_obj() {
	len := len(this.Username)
	var b []byte = []byte(this.Username)
	
	if b[len-1] == '_' {
		b = b[:len-1]
		this.Username = string(b)
	}
	return
}