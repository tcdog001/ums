package controllers

import (
	"errors"
	"ums/v1/models"
)

type DatabaseCheck struct {
}

//对数据库做健康检查
func (dc *DatabaseCheck) Check() error {
	if dc.isConnected() {
		return nil
	} else {
		return errors.New("can't connect the database")
	}
}

func (dc *DatabaseCheck) isConnected() bool {
	return models.CheckDatabase()
}
