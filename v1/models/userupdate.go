package models

type Userupdate struct {
	Usermac  string `json:"usermac"`
	Flowup   uint64 `json:"flowup"`
	Flowdown uint64 `json:"flowdown"`
}
