package account

import (
	"database/sql"
	"errors"

	"github.com/QuirkyQuestor/moneyKeeper/internal/datamodel"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler"
	"github.com/lib/pq"
	"log/slog"
)

var (
	ErrSQLExecution                = errors.New("error during sql stetement execution")
	ErrSQLInsert                   = errors.New("error when getting LastInsertId")
	ErrCannotPrepareSQLStatement   = errors.New("cannot prepare sql statement")
	ErrConvertingDBResponse        = errors.New("error during converting DB/Go types")
	ErrNoItemResponse              = errors.New("DB query returned no result")
	ErrUnexpectedDBExecutionResult = errors.New("unexpected DB statement ExecutionResult")
)

func GetAllAccounts(DBConnection *sql.DB, userID string) ([]datamodel.Account, error) {
	var accounts = []datamodel.Account{}

	bdStatement, err := DBConnection.Prepare("SELECT account_id, type_id, name, description, active FROM account WHERE user_id = $1;")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return accounts, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	rows, err := bdStatement.Query(userID)

	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return nil, ErrSQLExecution
	}
	defer rows.Close()

	for rows.Next() {
		var accountId string
		var typeId string
		var name string
		var description sql.NullString
		var active bool

		err = rows.Scan(&accountId, &typeId, &name, &description, &active)
		if err != nil {
			slog.Error(ErrConvertingDBResponse.Error(), "error", err)
			return accounts, ErrConvertingDBResponse
		}

		account := datamodel.Account{
			AccountID:   accountId,
			TypeID:      typeId,
			Name:        name,
			Description: description.String,
			Active:      active,
		}
		accounts = append(accounts, account)
		slog.Info("Got this...", "account", account)
	}
	return accounts, nil
}
func AddAccount(DBConnection *sql.DB, userID string, account datamodel.Account) (datamodel.Account, error) {
	slog.Info("Received Account object", "incomming_account", account)

	bdStatement, err := DBConnection.Prepare("INSERT INTO account(user_id, type_id, name, description, active) VALUES ($1, $2, $3, $4, $5) RETURNING account_id;")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return account, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	err = bdStatement.QueryRow(userID, account.TypeID, account.Name, account.Description, account.Active).Scan(&account.AccountID)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == sqlhandler.PGErrUniqueViolation {
				return account, sqlhandler.SQLConflict
			}
		}
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return account, ErrSQLExecution
	}
	return account, nil
}

func GetAccountByID(DBConnection *sql.DB, userID string, accountID string) (datamodel.Account, error) {
	var account datamodel.Account

	bdStatement, err := DBConnection.Prepare("SELECT account_id, type_id, name, description, active FROM account WHERE account_id = $1 AND user_id = $2;")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return account, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	var accountId string
	var typeId string
	var name string
	var description sql.NullString
	var active bool

	err = bdStatement.QueryRow(accountID, userID).Scan(&accountId, &typeId, &name, &description, &active)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.Info(ErrNoItemResponse.Error())
			return account, ErrNoItemResponse
		}
		slog.Error(ErrConvertingDBResponse.Error(), "error", err)
		return account, ErrConvertingDBResponse
	}
	account = datamodel.Account{
		AccountID:   accountId,
		TypeID:      typeId,
		Name:        name,
		Description: description.String,
		Active:      active,
	}

	return account, nil
}

func UpdateAccountByID(DBConnection *sql.DB, userID string, accountUpd *datamodel.Account) error {

	bdStatement, err := DBConnection.Prepare("UPDATE account SET type_id=$1, name=$2, description=$3, active=$4 WHERE account_id = $5 AND user_id = $6")
	if err != nil {
		slog.Error("cannot prepare update statement", "error", err)
		return err
	}
	defer bdStatement.Close()
	result, err := bdStatement.Exec(accountUpd.TypeID, accountUpd.Name, accountUpd.Description, accountUpd.Active, accountUpd.AccountID, userID)

	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return err
	}
	rowsUpdated, err := result.RowsAffected()
	if err != nil || rowsUpdated != 1 {
		slog.Error(ErrUnexpectedDBExecutionResult.Error(), "error", err, "result", result)
		return err
	}
	slog.Info("rowsUpdated", "rowsUpdated", rowsUpdated)

	return nil
}

func DeleteAccountByID(DBConnection *sql.DB, userID string, accountID string) error {

	bdStatement, err := DBConnection.Prepare("DELETE FROM account WHERE account_id = $1 AND user_id = $2")
	if err != nil {
		slog.Error("cannot prepare update statement", "error", err)
		return err
	}
	defer bdStatement.Close()
	result, err := bdStatement.Exec(accountID, userID)

	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return err
	}
	rowsUpdated, err := result.RowsAffected()
	if err != nil || rowsUpdated > 1 {
		slog.Error(ErrUnexpectedDBExecutionResult.Error(), "error", err, "result", result)
	}
	return nil
}
