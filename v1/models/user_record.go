package models

import (
	"time"
)

type UserRecord struct {
	Id 		   uint64 `orm:"auto"`
	UserName   string
	UserMac    string
	DevMac     string
	AuthTime   time.Time
	DeauthTime time.Time
}

func (this *UserRecord) TableName() string {
	return "Userrecord"
}

func (this *UserRecord) KeyName() string {
	return "username"
}

func (this *UserRecord) Key() string {
	return this.UserName
}

func LogUserRecord(user *UserStatus) {
	record := &UserRecord {
		UserName : user.UserName,
		UserMac : user.UserMac,
		DevMac : user.DevMac,
		AuthTime : user.AuthTime,
		DeauthTime : time.Now(),
	}
	
	dbEntryInsert(nil, record)
}
