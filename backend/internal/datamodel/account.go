package datamodel

// Account is to represent account info
type Account struct {
	AccountID   string `json:"accountId"`
	TypeID      string `json:"typeId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	IsExternal  bool   `json:"isExternal"`
}
