package datamodel

// Category is to represent category info
type Category struct {
	CategoryID  string  `json:"categoryId"`
	ParentID    *string `json:"parentId,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Expence     bool    `json:"expence"`
}
