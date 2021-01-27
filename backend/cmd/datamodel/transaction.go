package datamodel

import "time"

// Transaction is to represent the transaction info
type Transaction struct {
	TransactionID int64     `json:"transactionId"`
	AccountID     int64     `json:"accountId"`
	Date          time.Time `json:"date"`
	Amount        int64     `json:"amount"`
	RefAccount    int64     `json:"refAccount"`
	Memo          string    `json:"memo"`
	CategoryID    int64     `json:"categoryId"`
	TransferRefID int64     `json:"transferRefId"`
}
