package transaction

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/winchien/moneyKeeper/backend/cmd/datamodel"
	"github.com/winchien/moneyKeeper/backend/cmd/sqlhandler"
)

var (
	ErrSQLExecution              = errors.New("error during sql stetement execution")
	ErrSQLInsert                 = errors.New("error when getting LastInsertId")
	ErrCannotPrepareSQLStatement = errors.New("cannot prepare sql statement")
	ErrConvertingDBResponse      = errors.New("error during converting DB/Go types")
	ErrSQLUpdate                 = errors.New("error while updating the record in DB")
	ErrNoItemResponse            = errors.New("DB query returned no result")
)

func GetAllTransactions(DBConnection *sql.DB) ([]datamodel.Transaction, error) {
	var transactions = []datamodel.Transaction{}

	// json:"transaction_id"`
	// json:"account_from"`
	// json:"date"`
	// json:"amount"`
	// json:"account_to"`
	// json:"memo"`
	// json:"category_id"`
	// json:"transfer_transaction_id"` // omit empty

	bdStatement, err := DBConnection.Prepare("SELECT transaction_id, account_from, date, amount, account_to, memo, category_id, transfer_transaction_id FROM moneykeeper.transaction;")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return transactions, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()
	rows, err := bdStatement.Query()
	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return nil, ErrSQLExecution
	}
	defer rows.Close()

	for rows.Next() {
		var transactionId string
		var accountFrom string
		var accountTo string
		var date *time.Time
		var amount float64
		var memo sql.NullString
		var categoryId string
		var transferTransactionId *string

		err := rows.Scan(&transactionId, &accountFrom, &date, &amount, &accountTo, &memo, &categoryId, &transferTransactionId)
		if err != nil {
			log.WithError(err).Error("Could not parse row from the DB")
		}
		transaction := datamodel.Transaction{
			TransactionID:         transactionId,
			AccountFrom:           accountFrom,
			AccountTo:             accountTo,
			Date:                  date,
			Amount:                amount,
			CategoryID:            categoryId,
			Memo:                  memo.String,
			TransferTransactionID: transferTransactionId,
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func AddTransaction(DBConnection *sql.DB, transaction datamodel.Transaction) (datamodel.Transaction, error) {
	log.WithField("transaction", transaction).Info("The Transaction object")
	bdStatement, err := DBConnection.Prepare("INSERT INTO moneykeeper.transaction(account_from, date, amount, account_to, memo, category_id, transfer_transaction_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING transaction_id;")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return transaction, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	err = bdStatement.QueryRow(transaction.AccountFrom, transaction.Date, transaction.Amount, transaction.AccountTo, transaction.Memo, transaction.CategoryID, transaction.TransferTransactionID).Scan(&transaction.TransactionID)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == sqlhandler.PGErrUniqueViolation {
				return transaction, sqlhandler.SQLConflict
			}
		}
		log.WithError(err).Error(ErrSQLExecution)
		return transaction, ErrSQLExecution
	}

	return transaction, nil
}

func GetTransactionByID(DBConnection *sql.DB, transactionID string) (*datamodel.Transaction, error) {

	var transaction datamodel.Transaction
	bdStatement, err := DBConnection.Prepare("SELECT transaction_id, account_from, date, amount, account_to, memo, category_id, transfer_transaction_id FROM moneykeeper.transaction WHERE transaction_id = $1")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return nil, ErrCannotPrepareSQLStatement
	}
	defer bdStatement.Close()

	err = bdStatement.QueryRow(transactionID).Scan(&transaction.TransactionID, &transaction.AccountFrom, &transaction.Date,
		&transaction.Amount, &transaction.AccountTo, &transaction.Memo, &transaction.CategoryID, &transaction.TransferTransactionID)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Info(ErrNoItemResponse)
			return nil, ErrNoItemResponse
		}
		log.WithError(err).Error(ErrConvertingDBResponse)
		return nil, ErrConvertingDBResponse
	}

	return &transaction, nil
}

func UpdateTransaction(DBConnection *sql.DB, transaction *datamodel.Transaction) error {
	log.WithField("transaction", transaction).Info("The Transaction object")
	bdStatement, err := DBConnection.Prepare("UPDATE moneykeeper.transaction SET account_from=$1, date=$2, amount=$3, account_to=$4, memo=$5, category_id=$6, transfer_transaction_id=$7 WHERE transaction_id = $8")
	if err != nil {
		log.WithError(err).Error("cannot prepare update statement")
	}

	defer bdStatement.Close()
	result, err := bdStatement.Exec(transaction.AccountFrom, transaction.Date, transaction.Amount, transaction.AccountTo, transaction.Memo, transaction.CategoryID, transaction.TransferTransactionID, transaction.TransactionID)

	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return ErrSQLExecution
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithError(err).Error("Cannot get rowsAffected for UpdateTransaction")
		return ErrSQLUpdate
	}
	log.WithField("rowsAffected", rowsAffected).Info("rowsAffected")

	if rowsAffected == 0 {
		log.Error("The record does not seem to be updated.")
		return ErrNoItemResponse
	}

	return nil
}

func DeleteTransactionByID(DBConnection *sql.DB, transactionID string) error {

	bdStatement, err := DBConnection.Prepare("DELETE FROM moneykeeper.transaction WHERE transaction_id = $1")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return ErrCannotPrepareSQLStatement
	}
	defer bdStatement.Close()

	result, err := bdStatement.Exec(transactionID)
	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return ErrSQLExecution
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithError(err).Error("Cannot get rowsAffected for Delete transaction")
		return ErrSQLUpdate
	}

	if rowsAffected != 1 {
		log.WithField("rowsAffected", rowsAffected).Info("The requested transaction did not exist in the DB Table")
	}

	return nil
}
