package accountType

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/winchien/moneyKeeper/backend/cmd/datamodel"
	"github.com/winchien/moneyKeeper/backend/cmd/sqlhandler"
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

	bdStatement, err := DBConnection.Prepare("SELECT type_id, name, description FROM moneykeeper.account_type;")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return nil, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	rows, err := bdStatement.Query()

	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return nil, ErrSQLExecution
	}
	defer rows.Close()

	for rows.Next() {
		var typeId string
		var name string
		var description sql.NullString

		err = rows.Scan(&typeId, &name, &description)
		if err != nil {
			log.WithError(err).Error(ErrConvertingDBResponse)
			return nil, ErrConvertingDBResponse
		}

		account := datamodel.AccountType{
			TypeID:      typeId,
			Name:        name,
			Description: description.String,
		}
		accountTypes = append(accountTypes, account)
		log.WithField("account", account).Info("Got this...")
	}
	return accountTypes, nil
}

func AddAccountType(DBConnection *sql.DB, accountType datamodel.AccountType) (datamodel.AccountType, error) {
	log.WithField("accountType", accountType).Info("The AccountType object")
	bdStatement, err := DBConnection.Prepare(`INSERT INTO moneykeeper.account_type(name, description) VALUES ($1, $2) RETURNING type_id;`)
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
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
		log.WithError(err).Error(ErrSQLExecution)
		return accountType, ErrSQLExecution
	}
	return accountType, nil
}

func GetAccountTypeByID(DBConnection *sql.DB, accountTyepID string) (*datamodel.AccountType, error) {
	var accountType *datamodel.AccountType

	bdStatement, err := DBConnection.Prepare("SELECT type_id, name, description FROM moneykeeper.account_type WHERE type_id = $1;")
	if err != nil {
		log.WithError(err).Error("Cannot prepare SELECT statement")
		return nil, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	var typeId string
	var name string
	var description sql.NullString

	err = bdStatement.QueryRow(accountTyepID).Scan(&typeId, &name, &description)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Info(ErrNoItemResponse)
			return nil, ErrNoItemResponse
		}
		log.WithError(err).Error(ErrConvertingDBResponse)
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

	bdStatement, err := DBConnection.Prepare("UPDATE moneykeeper.account_type SET name=$1, description=$2 WHERE type_id = $3")
	if err != nil {
		log.WithError(err).Error("cannot prepare update statement")
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
		log.WithError(err).Error(ErrSQLExecution)
		return ErrSQLExecution
	}
	rowsUpdated, err := result.RowsAffected()
	if err != nil || rowsUpdated != 1 {
		log.WithError(err).WithField("result", result).Error(ErrUnexpectedDBExecutionResult)
		return err
	}
	log.WithField("rowsUpdated", rowsUpdated).Info("rowsUpdated")

	return nil
}

func DeleteAccountTypeByID(DBConnection *sql.DB, accountTypeID string) error {

	bdStatement, err := DBConnection.Prepare("DELETE FROM moneykeeper.account_type WHERE type_id = $1")
	if err != nil {
		log.WithError(err).Error("cannot prepare delete statement")
		return ErrCannotPrepareSQLStatement
	}
	defer bdStatement.Close()
	result, err := bdStatement.Exec(accountTypeID)

	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return ErrSQLExecution
	}
	rowsUpdated, err := result.RowsAffected()
	if err != nil || rowsUpdated > 1 {
		log.WithError(err).WithField("result", result).Error(ErrUnexpectedDBExecutionResult)
	}
	return nil
}
