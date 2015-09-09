package models

type UserUpdate struct {
	UserMac  string `json:"usermac"`
	FlowUp   uint64 `json:"flowup"`
	FlowDown uint64 `json:"flowdown"`
}
