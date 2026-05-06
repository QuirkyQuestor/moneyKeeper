package accountType

import (
	"database/sql"
	"errors"

	"log/slog"

	"github.com/QuirkyQuestor/moneyKeeper/internal/datamodel"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler"
	"github.com/lib/pq"
)

var (
	ErrSQLExecution                = errors.New("error during sql stetement execution")
	ErrSQLInsert                   = errors.New("error when getting LastInsertId")
	ErrCannotPrepareSQLStatement   = errors.New("cannot prepare sql statement")
	ErrConvertingDBResponse        = errors.New("error during converting DB/Go types")
	ErrNoItemResponse              = errors.New("DB query returned no result")
	ErrUnexpectedDBExecutionResult = errors.New("unexpected DB statement ExecutionResult")
)

func GetAllAccountTypes(DBConnection *sql.DB) ([]datamodel.AccountType, error) {
	var accountTypes = []datamodel.AccountType{}

	bdStatement, err := DBConnection.Prepare("SELECT type_id, name, description FROM account_type;")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return nil, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	rows, err := bdStatement.Query()

	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return nil, ErrSQLExecution
	}
	defer rows.Close()

	for rows.Next() {
		var typeId string
		var name string
		var description sql.NullString

		err = rows.Scan(&typeId, &name, &description)
		if err != nil {
			slog.Error(ErrConvertingDBResponse.Error(), "error", err)
			return nil, ErrConvertingDBResponse
		}

		accountType := datamodel.AccountType{
			TypeID:      typeId,
			Name:        name,
			Description: description.String,
		}
		accountTypes = append(accountTypes, accountType)
	}
	return accountTypes, nil
}

func AddAccountType(DBConnection *sql.DB, accountType datamodel.AccountType) (datamodel.AccountType, error) {
	slog.Info("The AccountType object", "accountType", accountType)
	bdStatement, err := DBConnection.Prepare(`INSERT INTO account_type(name, description) VALUES ($1, $2) RETURNING type_id;`)
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return accountType, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	err = bdStatement.QueryRow(accountType.Name, accountType.Description).Scan(&accountType.TypeID)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == sqlhandler.PGErrUniqueViolation {
				return accountType, sqlhandler.SQLConflict
			}
		}
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return accountType, ErrSQLExecution
	}
	return accountType, nil
}

func GetAccountTypeByID(DBConnection *sql.DB, accountTyepID string) (*datamodel.AccountType, error) {
	var accountType *datamodel.AccountType

	bdStatement, err := DBConnection.Prepare("SELECT type_id, name, description FROM account_type WHERE type_id = $1;")
	if err != nil {
		slog.Error("Cannot prepare SELECT statement", "error", err)
		return nil, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	var typeId string
	var name string
	var description sql.NullString

	err = bdStatement.QueryRow(accountTyepID).Scan(&typeId, &name, &description)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.Info(ErrNoItemResponse.Error())
			return nil, ErrNoItemResponse
		}
		slog.Error(ErrConvertingDBResponse.Error(), "error", err)
		return nil, ErrConvertingDBResponse
	}
	accountType = &datamodel.AccountType{
		TypeID:      typeId,
		Name:        name,
		Description: description.String,
	}

	return accountType, nil
}

func UpdateAccountTypeByID(DBConnection *sql.DB, accountTypeUpd *datamodel.AccountType) error {

	bdStatement, err := DBConnection.Prepare("UPDATE account_type SET name=$1, description=$2 WHERE type_id = $3")
	if err != nil {
		slog.Error("cannot prepare update statement", "error", err)
		return ErrCannotPrepareSQLStatement
	}
	defer bdStatement.Close()
	result, err := bdStatement.Exec(accountTypeUpd.Name, accountTypeUpd.Description, accountTypeUpd.TypeID)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == sqlhandler.PGErrUniqueViolation {
				return sqlhandler.SQLConflict
			}
		}
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return ErrSQLExecution
	}
	rowsUpdated, err := result.RowsAffected()
	if err != nil || rowsUpdated != 1 {
		slog.Error(ErrUnexpectedDBExecutionResult.Error(), "error", err, "result", result)
		return err
	}
	slog.Info("rowsUpdated", "rowsUpdated", rowsUpdated)

	return nil
}

func DeleteAccountTypeByID(DBConnection *sql.DB, accountTypeID string) error {

	bdStatement, err := DBConnection.Prepare("DELETE FROM account_type WHERE type_id = $1")
	if err != nil {
		slog.Error("cannot prepare delete statement", "error", err)
		return ErrCannotPrepareSQLStatement
	}
	defer bdStatement.Close()
	result, err := bdStatement.Exec(accountTypeID)

	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return ErrSQLExecution
	}
	rowsUpdated, err := result.RowsAffected()
	if err != nil || rowsUpdated > 1 {
		slog.Error(ErrUnexpectedDBExecutionResult.Error(), "error", err, "result", result)
	}
	return nil
}
