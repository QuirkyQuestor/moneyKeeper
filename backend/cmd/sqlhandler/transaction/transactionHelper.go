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
		var date time.Time
		var amount float64
		var memo sql.NullString
		var categoryId string
		var transferTransactionId string

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
	bdStatement, err := DBConnection.Prepare("INSERT INTO moneykeeper.transaction(account_from, date, amount, account_to, memo, category_id, transfer_ref_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING transaction_id;")
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

// func GetCategoryByID(DBConnection *sql.DB, categoryID string) (*datamodel.Category, error) {

// 	var category datamodel.Category
// 	bdStatement, err := DBConnection.Prepare("SELECT categoryId, parentId, name, description, expence FROM category WHERE categoryID = $1")
// 	if err != nil {
// 		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
// 		return nil, ErrCannotPrepareSQLStatement
// 	}
// 	defer bdStatement.Close()

// 	var categoryId int
// 	var parentId int
// 	var name string
// 	var description sql.NullString
// 	var expence bool

// 	err = bdStatement.QueryRow(categoryID).Scan(&categoryId, &parentId, &name, &description, &expence)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Info(ErrNoItemResponse)
// 			return nil, ErrNoItemResponse
// 		}
// 		log.WithError(err).Error(ErrConvertingDBResponse)
// 		return nil, ErrConvertingDBResponse
// 	}

// 	category = datamodel.Category{
// 		CategoryID:  int64(categoryId),
// 		ParentID:    int64(parentId),
// 		Name:        name,
// 		Description: description.String,
// 		Expence:     expence,
// 	}
// 	return &category, nil
// }

// func UpdateCategory(DBConnection *sql.DB, category *datamodel.Category) error {
// 	log.WithField("category", category).Info("The Category object")
// 	bdStatement, err := DBConnection.Prepare("UPDATE category SET parentID=?, name=?, description=?, expence=? WHERE categoryId = ?")
// 	if err != nil {
// 		log.WithError(err).Error("cannot prepare update statement")
// 	}

// 	defer bdStatement.Close()
// 	result, err := bdStatement.Exec(category.ParentID, category.Name, category.Description, category.Expence, category.CategoryID)

// 	if err != nil {
// 		log.WithError(err).Error(ErrSQLExecution)
// 		return ErrSQLExecution
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		log.WithError(err).Error("Cannot get rowsAffected for Delete")
// 		return ErrSQLUpdate
// 	}
// 	log.WithField("rowsAffected", rowsAffected).Info("rowsAffected")

// 	if rowsAffected != 1 {
// 		log.Error("The record does not seem to be updated.")
// 		return ErrSQLUpdate
// 	}

// 	return nil
// }

// func DeleteCategoryByID(DBConnection *sql.DB, categoryID int64) error {

// 	bdStatement, err := DBConnection.Prepare("DELETE FROM category WHERE categoryId = ?")
// 	if err != nil {
// 		log.WithError(err).Error("Cannot prepare SQL statement")
// return ErrCannotPrepareSQLStatement
// 	}
// 	defer bdStatement.Close()

// 	result, err := bdStatement.Exec(categoryID)
// 	if err != nil {
// 		log.WithError(err).Error(ErrSQLExecution)
// 		return ErrSQLExecution
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		log.WithError(err).Error("Cannot get rowsAffected for Delete")
// 		return ErrSQLUpdate
// 	}
// 	log.WithField("rowsAffected", rowsAffected).Info("rowsAffected")

// 	if rowsAffected != 1 {
// 		log.Error("The record does not seem to be updated.")
// 		return ErrSQLUpdate
// 	}

// 	return nil
// }
