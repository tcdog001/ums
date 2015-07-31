package models

import (
	// "asdf"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Userstatus struct {
	Username string    `json:"username"`
	Usermac  string    `orm:"pk";json:"usermac"`
	Devmac   string    `json:"devmac"`
	Authcode string    `json:"authcode"`
	Flowup   uint64    `json:"flowup"`
	Flowdown uint64    `json:"flowdown"`
	AuthTime time.Time `orm:"type(datetime)";json:"-"`
	session []byte
}

func (user *Userstatus) TableName() string {
	return "userstatus"
}
func RegisterUserstatus(user *Userstatus) bool {
	o := orm.NewOrm()
	user.AuthTime = time.Now()
	beego.Debug("userstatus table=", user.TableName())
	//查找对应的mac地址是否存在，存在的话要求状态为离线
	exist := o.QueryTable(user.TableName()).Filter("usermac", user.Usermac).Exist()
	if exist {
		//用户存在，则更新用户信息
		var u Userstatus
		err := o.QueryTable(user.TableName()).Filter("usermac", user.Usermac).One(&u)
		if err != nil {
			beego.Debug("get item in userstatus failed!!")
			return false
		} else {
			u.Username = user.Username
			u.Usermac = user.Usermac
			u.Authcode = user.Authcode
			u.Devmac = user.Devmac
			_, err := o.Update(&u)
			if err != nil {
				beego.Error(err)
				return false
			}
		}
	} else {
		//用户不存在，则插入用户信息
		_, err := o.Insert(user)
		if err != nil {
			beego.Error(err)
			return false
		}
	}
	return true
}

func UpdateUserstatusBymac(user *Userstatus) bool {
	beego.Debug("Update userstatus table=", user.TableName())
	o := orm.NewOrm()

	var u Userstatus
	err := o.QueryTable(user.TableName()).Filter("usermac", user.Usermac).One(&u)
	if err != nil {
		return false
	} else {
		u.Flowup = user.Flowup
		u.Flowdown = user.Flowdown
		beego.Debug("Update userstatus usermac = ", u.Usermac)

		_, err := o.Update(&u)
		if err != nil {
			beego.Error(err)
			return false
		}
		return true
	}
}

func FindUserstatusByMac(user *Userstatus) bool {
	o := orm.NewOrm()
	exist := o.QueryTable(user.TableName()).Filter("usermac", user.Usermac).Exist()
	return exist
}

func DelUserStatusByMac(user *Userstatus) bool {
	o := orm.NewOrm()

	var u Userstatus
	err := o.QueryTable(user.TableName()).Filter("usermac", user.Usermac).One(&u)
	if err != nil {
		beego.Error(err)
		return false
	}

	_, err = o.Delete(&u)
	if err != nil {
		beego.Error(err)
		return false
	}

	return true
}
