package transaction

import (
	"database/sql"
	"errors"
	"time"

	"github.com/QuirkyQuestor/moneyKeeper/internal/datamodel"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler"
	"github.com/lib/pq"
	"log/slog"
)

var (
	ErrSQLExecution              = errors.New("error during sql stetement execution")
	ErrSQLInsert                 = errors.New("error when getting LastInsertId")
	ErrCannotPrepareSQLStatement = errors.New("cannot prepare sql statement")
	ErrConvertingDBResponse      = errors.New("error during converting DB/Go types")
	ErrSQLUpdate                 = errors.New("error while updating the record in DB")
	ErrNoItemResponse            = errors.New("DB query returned no result")
)

func GetAllTransactions(DBConnection *sql.DB, userID string) ([]datamodel.Transaction, error) {
	var transactions = []datamodel.Transaction{}

	bdStatement, err := DBConnection.Prepare("SELECT transaction_id, account_from, date, amount, account_to, memo, category_id, transfer_transaction_id FROM transaction WHERE user_id = $1;")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return transactions, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()
	rows, err := bdStatement.Query(userID)
	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
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
			slog.Error("Could not parse row from the DB", "error", err)
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

func GetTransactionsByAccountId(DBConnection *sql.DB, userID string, accountFrom string) ([]datamodel.Transaction, error) {
	var transactions = []datamodel.Transaction{}

	bdStatement, err := DBConnection.Prepare("SELECT transaction_id, account_from, date, amount, account_to, memo, category_id, transfer_transaction_id FROM transaction WHERE account_from = $1 AND user_id = $2;")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return transactions, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()
	rows, err := bdStatement.Query(accountFrom, userID)
	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
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
			slog.Error("Could not parse row from the DB", "error", err)
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

func AddTransaction(DBConnection *sql.DB, userID string, transaction datamodel.Transaction) (datamodel.Transaction, error) {
	slog.Info("The Transaction object", "transaction", transaction)
	bdStatement, err := DBConnection.Prepare("INSERT INTO transaction(user_id, account_from, date, amount, account_to, memo, category_id, transfer_transaction_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING transaction_id;")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return transaction, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	err = bdStatement.QueryRow(userID, transaction.AccountFrom, transaction.Date, transaction.Amount, transaction.AccountTo, transaction.Memo, transaction.CategoryID, transaction.TransferTransactionID).Scan(&transaction.TransactionID)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == sqlhandler.PGErrUniqueViolation {
				return transaction, sqlhandler.SQLConflict
			}
		}
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return transaction, ErrSQLExecution
	}

	return transaction, nil
}

func GetTransactionByID(DBConnection *sql.DB, userID string, transactionID string) (*datamodel.Transaction, error) {

	var transaction datamodel.Transaction
	bdStatement, err := DBConnection.Prepare("SELECT transaction_id, account_from, date, amount, account_to, memo, category_id, transfer_transaction_id FROM transaction WHERE transaction_id = $1 AND user_id = $2")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return nil, ErrCannotPrepareSQLStatement
	}
	defer bdStatement.Close()

	err = bdStatement.QueryRow(transactionID, userID).Scan(&transaction.TransactionID, &transaction.AccountFrom, &transaction.Date,
		&transaction.Amount, &transaction.AccountTo, &transaction.Memo, &transaction.CategoryID, &transaction.TransferTransactionID)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.Info(ErrNoItemResponse.Error())
			return nil, ErrNoItemResponse
		}
		slog.Error(ErrConvertingDBResponse.Error(), "error", err)
		return nil, ErrConvertingDBResponse
	}

	return &transaction, nil
}

func UpdateTransaction(DBConnection *sql.DB, userID string, transaction *datamodel.Transaction) error {
	slog.Info("The Transaction object", "transaction", transaction)
	bdStatement, err := DBConnection.Prepare("UPDATE transaction SET account_from=$1, date=$2, amount=$3, account_to=$4, memo=$5, category_id=$6, transfer_transaction_id=$7 WHERE transaction_id = $8 AND user_id = $9")
	if err != nil {
		slog.Error("cannot prepare update statement", "error", err)
	}

	defer bdStatement.Close()
	result, err := bdStatement.Exec(transaction.AccountFrom, transaction.Date, transaction.Amount, transaction.AccountTo, transaction.Memo, transaction.CategoryID, transaction.TransferTransactionID, transaction.TransactionID, userID)

	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return ErrSQLExecution
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Cannot get rowsAffected for UpdateTransaction", "error", err)
		return ErrSQLUpdate
	}
	slog.Info("rowsAffected", "rowsAffected", rowsAffected)

	if rowsAffected == 0 {
		slog.Error("The record does not seem to be updated.")
		return ErrNoItemResponse
	}

	return nil
}

func DeleteTransactionByID(DBConnection *sql.DB, userID string, transactionID string) error {

	bdStatement, err := DBConnection.Prepare("DELETE FROM transaction WHERE transaction_id = $1 AND user_id = $2")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return ErrCannotPrepareSQLStatement
	}
	defer bdStatement.Close()

	result, err := bdStatement.Exec(transactionID, userID)
	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return ErrSQLExecution
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Cannot get rowsAffected for Delete transaction", "error", err)
		return ErrSQLUpdate
	}

	if rowsAffected != 1 {
		slog.Info("The requested transaction did not exist in the DB Table", "rowsAffected", rowsAffected)
	}

	return nil
}
