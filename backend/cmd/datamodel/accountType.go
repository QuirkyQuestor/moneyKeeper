package datamodel

// AccountType is to represent account info
type AccountType struct {
	TypeID      string `json:"typeId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
