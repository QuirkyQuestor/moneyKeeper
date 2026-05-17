package datamodel

import (
	"encoding/json"

	"log/slog"
)

// Category is to represent category info
type Category struct {
	CategoryID  *string `json:"categoryId,omitempty"`
	ParentID    *string `json:"parentId,omitempty"`
	Name        string  `json:"name"`
	FullName    string  `json:"fullName,omitempty"`
	Description string  `json:"description"`
	Expense     bool    `json:"expense"`
}

func (c *Category) UnmarshalJSON(data []byte) error {

	type categoryTmp struct {
		CategoryID  string `json:"categoryId,omitempty"`
		ParentID    string `json:"parentId,omitempty"`
		Name        string `json:"name"`
		FullName    string `json:"fullName,omitempty"`
		Description string `json:"description"`
		Expense     bool   `json:"expense"`
	}

	var tmp categoryTmp

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if tmp.CategoryID == "" {
		c.CategoryID = nil
	} else {
		c.CategoryID = &tmp.CategoryID
	}
	if tmp.ParentID == "" {
		c.ParentID = nil
	} else {
		c.ParentID = &tmp.ParentID
	}
	c.Name = tmp.Name
	c.FullName = tmp.FullName
	c.Description = tmp.Description
	c.Expense = tmp.Expense

	slog.Info("unmarshaling", "c", c)
	return nil
}
