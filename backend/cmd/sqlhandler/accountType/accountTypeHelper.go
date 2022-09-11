package accountType

import (
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/winchien/moneyKeeper/backend/cmd/datamodel"
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
	var accountTypes []datamodel.AccountType

	bdStatement, err := DBConnection.Prepare("SELECT typeId, name, description, active FROM accountType;")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return nil, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	rows, err := bdStatement.Query()

	if err != nil {
		log.WithError(err).Fatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var typeId int64
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
	log.WithField("account", accountType).Info("The AccountType object")
	bdStatement, err := DBConnection.Prepare("INSERT INTO accountType(name, description) VALUES (?, ?);")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return accountType, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	result, err := bdStatement.Exec(accountType.Name, accountType.Description)

	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return accountType, ErrSQLExecution
	}
	newId, err := result.LastInsertId()
	if err != nil {
		log.WithError(err).Error(err.Error())
		return accountType, ErrSQLInsert
	}
	accountType.TypeID = newId
	return accountType, nil
}

func GetAccountTypeByID(DBConnection *sql.DB, accountTyepID string) (*datamodel.AccountType, error) {
	var accountType *datamodel.AccountType

	bdStatement, err := DBConnection.Prepare("SELECT typeId, name, description FROM accountType WHERE typeId = ?;")
	if err != nil {
		log.WithError(err).Error("Cannot prepare SELECT statement")
		return nil, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	var typeId int64
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

	bdStatement, err := DBConnection.Prepare("UPDATE accountType SET name=?, description=? WHERE typeId = ?")
	if err != nil {
		log.WithError(err).Error("cannot prepare update statement")
		return ErrCannotPrepareSQLStatement
	}
	defer bdStatement.Close()
	result, err := bdStatement.Exec(accountTypeUpd.Name, accountTypeUpd.Description, accountTypeUpd.TypeID)

	if err != nil {
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

	bdStatement, err := DBConnection.Prepare("DELETE FROM accountType WHERE typeId = ?")
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
