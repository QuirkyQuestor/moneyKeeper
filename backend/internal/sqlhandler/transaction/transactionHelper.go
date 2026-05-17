package transaction

import (
	"database/sql"
	"errors"
	"time"

	"github.com/QuirkyQuestor/moneyKeeper/internal/datamodel"
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

func GetAllTransactions(DBConnection *sql.DB, userID string, limit, offset int) ([]datamodel.Transaction, int, error) {
	var transactions = []datamodel.Transaction{}

	var totalCount int
	err := DBConnection.QueryRow("SELECT COUNT(*) FROM transaction WHERE user_id = $1", userID).Scan(&totalCount)
	if err != nil {
		slog.Error("Could not get total count of transactions", "error", err)
		return nil, 0, err
	}

	bdStatement, err := DBConnection.Prepare("SELECT transaction_id, account_from, date, amount, account_to, memo, category_id, transfer_transaction_id FROM transaction WHERE user_id = $1 ORDER BY date DESC LIMIT $2 OFFSET $3;")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return transactions, 0, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()
	rows, err := bdStatement.Query(userID, limit, offset)
	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return nil, 0, ErrSQLExecution
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
	return transactions, totalCount, nil
}

func GetTransactionsByAccountId(DBConnection *sql.DB, userID string, accountFrom string, limit, offset int) ([]datamodel.Transaction, int, error) {
	var transactions = []datamodel.Transaction{}

	var totalCount int
	err := DBConnection.QueryRow("SELECT COUNT(*) FROM transaction WHERE account_from = $1 AND user_id = $2", accountFrom, userID).Scan(&totalCount)
	if err != nil {
		slog.Error("Could not get total count of transactions for account", "error", err)
		return nil, 0, err
	}

	bdStatement, err := DBConnection.Prepare("SELECT transaction_id, account_from, date, amount, account_to, memo, category_id, transfer_transaction_id FROM transaction WHERE account_from = $1 AND user_id = $2 ORDER BY date DESC LIMIT $3 OFFSET $4;")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return transactions, 0, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()
	rows, err := bdStatement.Query(accountFrom, userID, limit, offset)
	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return nil, 0, ErrSQLExecution
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
	return transactions, totalCount, nil
}

func AddTransaction(DBConnection *sql.DB, userID string, txData datamodel.Transaction) (datamodel.Transaction, error) {
	slog.Info("Adding transaction", "transaction", txData)

	// Start a SQL transaction
	dbTx, err := DBConnection.Begin()
	if err != nil {
		slog.Error("Could not start DB transaction", "error", err)
		return txData, err
	}
	defer dbTx.Rollback()

	// Check if AccountTo is internal
	var isToInternal bool
	err = dbTx.QueryRow("SELECT NOT is_external FROM account WHERE account_id = $1 AND user_id = $2", txData.AccountTo, userID).Scan(&isToInternal)
	if err != nil && err != sql.ErrNoRows {
		slog.Error("Error checking if AccountTo is internal", "error", err)
		return txData, err
	}

	if isToInternal {
		// 1. Insert Transaction A (Source)
		var txID_A string
		err = dbTx.QueryRow(`
			INSERT INTO transaction(user_id, account_from, date, amount, account_to, memo, category_id, transfer_transaction_id) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING transaction_id`,
			userID, txData.AccountFrom, txData.Date, txData.Amount, txData.AccountTo, txData.Memo, txData.CategoryID, nil,
		).Scan(&txID_A)
		if err != nil {
			slog.Error("Error inserting source transaction", "error", err)
			return txData, err
		}

		// 2. Insert Transaction B (Destination)
		var txID_B string
		err = dbTx.QueryRow(`
			INSERT INTO transaction(user_id, account_from, date, amount, account_to, memo, category_id, transfer_transaction_id) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING transaction_id`,
			userID, txData.AccountTo, txData.Date, -txData.Amount, txData.AccountFrom, txData.Memo, txData.CategoryID, txID_A,
		).Scan(&txID_B)
		if err != nil {
			slog.Error("Error inserting destination transaction", "error", err)
			return txData, err
		}

		// 3. Update Transaction A with ID of B
		_, err = dbTx.Exec("UPDATE transaction SET transfer_transaction_id = $1 WHERE transaction_id = $2", txID_B, txID_A)
		if err != nil {
			slog.Error("Error updating source transaction link", "error", err)
			return txData, err
		}

		if err := dbTx.Commit(); err != nil {
			return txData, err
		}
		txData.TransactionID = txID_A
		transferID := txID_B
		txData.TransferTransactionID = &transferID
		return txData, nil
	}

	// Regular (non-transfer) transaction
	err = dbTx.QueryRow(`
		INSERT INTO transaction(user_id, account_from, date, amount, account_to, memo, category_id, transfer_transaction_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING transaction_id`,
		userID, txData.AccountFrom, txData.Date, txData.Amount, txData.AccountTo, txData.Memo, txData.CategoryID, txData.TransferTransactionID,
	).Scan(&txData.TransactionID)

	if err != nil {
		slog.Error("Error inserting transaction", "error", err)
		return txData, err
	}

	if err := dbTx.Commit(); err != nil {
		return txData, err
	}

	return txData, nil
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

func UpdateTransaction(DBConnection *sql.DB, userID string, txData *datamodel.Transaction) error {
	slog.Info("Updating transaction", "transaction", txData)

	dbTx, err := DBConnection.Begin()
	if err != nil {
		slog.Error("Could not start DB transaction", "error", err)
		return err
	}
	defer dbTx.Rollback()

	// 1. Get current state from DB to check for linked transaction
	var currentTransferID *string
	err = dbTx.QueryRow("SELECT transfer_transaction_id FROM transaction WHERE transaction_id = $1 AND user_id = $2", txData.TransactionID, userID).Scan(&currentTransferID)
	if err != nil && err != sql.ErrNoRows {
		slog.Error("Error fetching current transaction state", "error", err)
		return err
	}

	// Use the DB's transfer ID if the payload one is missing
	if (txData.TransferTransactionID == nil || *txData.TransferTransactionID == "") && currentTransferID != nil {
		txData.TransferTransactionID = currentTransferID
	}

	// 2. Update the primary transaction
	_, err = dbTx.Exec(`
		UPDATE transaction 
		SET account_from=$1, date=$2, amount=$3, account_to=$4, memo=$5, category_id=$6, transfer_transaction_id=$7 
		WHERE transaction_id = $8 AND user_id = $9`,
		txData.AccountFrom, txData.Date, txData.Amount, txData.AccountTo, txData.Memo, txData.CategoryID, txData.TransferTransactionID, txData.TransactionID, userID,
	)
	if err != nil {
		slog.Error("Error updating primary transaction", "error", err)
		return err
	}

	// 3. Update the linked transaction if it exists
	if txData.TransferTransactionID != nil && *txData.TransferTransactionID != "" {
		_, err = dbTx.Exec(`
			UPDATE transaction 
			SET date=$1, amount=$2, memo=$3, category_id=$4, account_from=$5, account_to=$6
			WHERE transaction_id = $7 AND user_id = $8`,
			txData.Date, -txData.Amount, txData.Memo, txData.CategoryID, txData.AccountTo, txData.AccountFrom, *txData.TransferTransactionID, userID,
		)
		if err != nil {
			slog.Error("Error updating linked transaction", "error", err)
			return err
		}
	}

	if err := dbTx.Commit(); err != nil {
		return err
	}

	return nil
}

func DeleteTransactionByID(DBConnection *sql.DB, userID string, transactionID string) error {
	slog.Info("Deleting transaction", "transactionID", transactionID)

	dbTx, err := DBConnection.Begin()
	if err != nil {
		slog.Error("Could not start DB transaction", "error", err)
		return err
	}
	defer dbTx.Rollback()

	// Get the linked transaction ID if it exists
	var transferID *string
	err = dbTx.QueryRow("SELECT transfer_transaction_id FROM transaction WHERE transaction_id = $1 AND user_id = $2", transactionID, userID).Scan(&transferID)
	if err != nil && err != sql.ErrNoRows {
		slog.Error("Error checking for linked transaction before delete", "error", err)
		return err
	}

	// If there's a linked transaction, we need to clear the foreign key references first
	// to avoid "violates foreign key constraint" due to the circular dependency.
	if transferID != nil && *transferID != "" {
		_, err = dbTx.Exec("UPDATE transaction SET transfer_transaction_id = NULL WHERE (transaction_id = $1 OR transaction_id = $2) AND user_id = $3", transactionID, *transferID, userID)
		if err != nil {
			slog.Error("Error clearing transfer links before delete", "error", err)
			return err
		}

		// Delete both
		_, err = dbTx.Exec("DELETE FROM transaction WHERE (transaction_id = $1 OR transaction_id = $2) AND user_id = $3", transactionID, *transferID, userID)
		if err != nil {
			slog.Error("Error deleting linked transaction pair", "error", err)
			return err
		}
	} else {
		// Delete only the primary transaction
		_, err = dbTx.Exec("DELETE FROM transaction WHERE transaction_id = $1 AND user_id = $2", transactionID, userID)
		if err != nil {
			slog.Error("Error deleting primary transaction", "error", err)
			return err
		}
	}

	if err := dbTx.Commit(); err != nil {
		return err
	}

	return nil
}
