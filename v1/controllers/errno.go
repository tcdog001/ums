package controllers

import (
	. "asdf"
)

type EUmsError int

const (
	ErrUmsBegin EUmsError = 0

	ErrUmsOk                      EUmsError = 0
	ErrUmsSmsError                EUmsError = 1
	ErrUmsInputError              EUmsError = 2
	ErrUmsUserInfoRegistered      EUmsError = 3
	ErrUmsUserInfoNotRegistered   EUmsError = 4
	ErrUmsUserInfoRegisterError   EUmsError = 5
	ErrUmsUserStatusNotExist      EUmsError = 6
	ErrUmsUserHaveBeenDeauthed    EUmsError = 7
	ErrUmsUserStatusDeleteError   EUmsError = 8
	ErrUmsUserStatusRegisterError EUmsError = 9
	ErrUmsRadError                EUmsError = 10
	ErrUmsRadAuthError            EUmsError = 11
	ErrUmsRadAcctStartError       EUmsError = 12
	ErrUmsRadAcctUpdateError      EUmsError = 13
	ErrUmsRadAcctStopError        EUmsError = 14

	ErrUmsEnd EUmsError = 15
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
	ErrUmsOk:                      "ok",
	ErrUmsSmsError:                "sms error",
	ErrUmsInputError:              "input error",
	ErrUmsUserInfoRegistered:      "user info have registered",
	ErrUmsUserInfoNotRegistered:   "user info NOT registered",
	ErrUmsUserInfoRegisterError:   "user info register error",
	ErrUmsUserStatusNotExist:      "user status NOT exist",
	ErrUmsUserHaveBeenDeauthed:    "user status have been deleted",
	ErrUmsUserStatusDeleteError:   "user status delete error",
	ErrUmsUserStatusRegisterError: "user status register error",
	ErrUmsRadError:                "radius error",
	ErrUmsRadAuthError:            "radius auth error",
	ErrUmsRadAcctStartError:       "radius acct start error",
	ErrUmsRadAcctUpdateError:      "radius acct update error",
	ErrUmsRadAcctStopError:        "radius acct stop error",
}
