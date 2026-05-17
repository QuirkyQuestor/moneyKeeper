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

// AccountBalance represents the current balance of an account
type AccountBalance struct {
	AccountID string  `json:"accountId"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
}
