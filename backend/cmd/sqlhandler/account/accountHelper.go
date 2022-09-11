package account

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

func GetAllAccounts(DBConnection *sql.DB) ([]datamodel.Account, error) {
	var accounts []datamodel.Account

	bdStatement, err := DBConnection.Prepare("SELECT accountId, typeId, name, description, active FROM account;")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return accounts, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	rows, err := bdStatement.Query()

	if err != nil {
		log.WithError(err).Fatal(err.Error())
	}
	defer rows.Close()

	log.Info("Test")

	for rows.Next() {
		var accountId int64
		var typeId int64
		var name string
		var description sql.NullString
		var active bool

		err = rows.Scan(&accountId, &typeId, &name, &description, &active)
		if err != nil {
			log.WithError(err).Error(ErrConvertingDBResponse)
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
		log.WithField("account", account).Info("Got this...")
	}
	return accounts, nil
}
func AddAccount(DBConnection *sql.DB, account datamodel.Account) (datamodel.Account, error) {
	log.WithField("account", account).Info("The Account object")
	bdStatement, err := DBConnection.Prepare("INSERT INTO account(typeId, name, description, active) VALUES (?, ?, ?, ?);")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return account, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	result, err := bdStatement.Exec(account.TypeID, account.Name, account.Description, account.Active)

	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return account, ErrSQLExecution
	}
	newId, err := result.LastInsertId()
	if err != nil {
		log.WithError(err).Error(err.Error())
		return account, ErrSQLInsert
	}
	account.AccountID = newId
	return account, nil
}
func GetAccountByID(DBConnection *sql.DB, accountID string) (datamodel.Account, error) {
	var account datamodel.Account

	bdStatement, err := DBConnection.Prepare("SELECT accountId, typeId, name, description, active FROM account WHERE accountId = ?;")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return account, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	var accountId int64
	var typeId int64
	var name string
	var description sql.NullString
	var active bool

	err = bdStatement.QueryRow(accountID).Scan(&accountId, &typeId, &name, &description, &active)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Info(ErrNoItemResponse)
			return account, ErrNoItemResponse
		}
		log.WithError(err).Error(ErrConvertingDBResponse)
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
func UpdateAccountByID(DBConnection *sql.DB, accountUpd *datamodel.Account) error {

	bdStatement, err := DBConnection.Prepare("UPDATE account SET typeId=?, name=?, description=?, active=? WHERE accountId = ?")
	if err != nil {
		log.WithError(err).Error("cannot prepare update statement")
		return err
	}
	defer bdStatement.Close()
	result, err := bdStatement.Exec(accountUpd.TypeID, accountUpd.Name, accountUpd.Description, accountUpd.Active, accountUpd.AccountID)

	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return err
	}
	rowsUpdated, err := result.RowsAffected()
	if err != nil || rowsUpdated != 1 {
		log.WithError(err).WithField("result", result).Error(ErrUnexpectedDBExecutionResult)
		return err
	}
	log.WithField("rowsUpdated", rowsUpdated).Info("rowsUpdated")

	return nil
}

func DeleteAccountByID(DBConnection *sql.DB, accountID string) error {

	bdStatement, err := DBConnection.Prepare("DELETE FROM account WHERE accountId = ?")
	if err != nil {
		log.WithError(err).Error("cannot prepare update statement")
		return err
	}
	defer bdStatement.Close()
	result, err := bdStatement.Exec(accountID)

	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return err
	}
	rowsUpdated, err := result.RowsAffected()
	if err != nil || rowsUpdated > 1 {
		log.WithError(err).WithField("result", result).Error(ErrUnexpectedDBExecutionResult)
	}
	return nil
}
