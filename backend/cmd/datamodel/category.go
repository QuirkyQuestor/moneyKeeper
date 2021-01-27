package datamodel

// Category is to represent category info
type Category struct {
	CategoryID  int64  `json:"categoryId"`
	ParentID    int64  `json:"parentId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Expence     bool   `json:"expence"`
}
