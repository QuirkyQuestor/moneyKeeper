package datamodel

import "time"

// Transaction is to represent the transaction info
type Transaction struct {
	TransactionID         string     `json:"transactionId"`
	AccountFrom           string     `json:"accountFrom"`
	Date                  *time.Time `json:"date"`
	Amount                float64    `json:"amount"`
	AccountTo             string     `json:"accountTo"`
	Memo                  string     `json:"memo"`
	CategoryID            string     `json:"categoryId"`
	TransferTransactionID *string    `json:"transferTransactionId,omitempty"`
}
