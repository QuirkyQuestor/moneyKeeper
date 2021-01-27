package datamodel

// Account is to represent account info
type Account struct {
	AccountID   int64  `json:"accountId"`
	TypeID      int64  `json:"typeId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}
