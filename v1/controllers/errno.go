package controllers

import (
	. "asdf"
)

type EUmsError 	int

const (
	ErrUmsBegin						EUmsError = 0
	
	ErrUmsOk						EUmsError = 0
	ErrUmsSmsError					EUmsError = 1
	ErrUmsInputError 				EUmsError = 2
	ErrUmsUserInfoRegistered		EUmsError = 3
	ErrUmsUserInfoNotRegistered 	EUmsError = 4
	ErrUmsUserInfoRegisterError		EUmsError = 5
	ErrUmsUserStatusNotExist 		EUmsError = 6
	ErrUmsUserStatusDeleteError 	EUmsError = 7
	ErrUmsUserStatusRegisterError 	EUmsError = 8
	ErrUmsRadError					EUmsError = 9
	ErrUmsRadAuthError				EUmsError = 10
	ErrUmsRadAcctStartError			EUmsError = 11
	ErrUmsRadAcctUpdateError		EUmsError = 12
	ErrUmsRadAcctStopError			EUmsError = 13
	
	ErrUmsEnd 						EUmsError = 14
)


func (me EUmsError) Tag() string {
	return "Ums Error"
}

func (me EUmsError) Begin() int {
	return int(ErrUmsBegin)
}

func (me EUmsError) End() int {
	return int(ErrUmsEnd)
}

func (me EUmsError) Int() int {
	return int(me)
}

func (me EUmsError) IsGood() bool {
	if !IsGoodEnum(me) {
		Log.Error("bad attr(%s) value(%d)", me.Tag(), me)

		return false
	} else if 0 == len(errUmsBind[me]) {
		Log.Error("no support attr(%s) value(%d)", me.Tag(), me)

		return false
	}

	return true
}

func (me EUmsError) ToString() string {
	var b EnumBinding = errUmsBind[:]

	return b.EntryShow(me)
}

var errUmsBind = [ErrUmsEnd]string{
	ErrUmsOk:						"ok",
	ErrUmsSmsError:					"sms error",
	ErrUmsInputError:				"input error",
	ErrUmsUserInfoRegistered:		"user info have registered",
	ErrUmsUserInfoNotRegistered:	"user info NOT registered",
	ErrUmsUserInfoRegisterError:	"user info register error",
	ErrUmsUserStatusNotExist:		"user status NOT exist",
	ErrUmsUserStatusDeleteError:	"user status delete error",
	ErrUmsUserStatusRegisterError:	"user status register error",
	ErrUmsRadError:					"radius error",
	ErrUmsRadAuthError:				"radius auth error",
	ErrUmsRadAcctStartError:		"radius acct start error",
	ErrUmsRadAcctUpdateError:		"radius acct update error",
	ErrUmsRadAcctStopError:			"radius acct stop error",
}
