package datamodel

// AccountType is to represent account info
type AccountType struct {
	TypeID      int64  `json:"typeId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
