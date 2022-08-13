package account

import (
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/winchien/moneyKeeper/backend/cmd/datamodel"
)

var (
	ErrSQLExecutionError              = errors.New("error during sql stetement execution")
	ErrSQLInsertError                 = errors.New("error when getting LastInsertId")
	ErrCannotPrepareSQLStatementError = errors.New("cannot prepare sql statement")
	ErrConvertingDBResponseError      = errors.New("error during converting DB/Go types")
)

func GetAllAccounts(DBConnection *sql.DB) ([]datamodel.Account, error) {
	var accounts []datamodel.Account

	bdStatement, err := DBConnection.Prepare("SELECT accountId, typeId, name, description, active FROM account;")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatementError)
		return accounts, ErrCannotPrepareSQLStatementError
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
			log.WithError(err).Error(ErrConvertingDBResponseError)
			return accounts, ErrConvertingDBResponseError
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
		log.WithError(err).Error(ErrCannotPrepareSQLStatementError)
		return account, ErrCannotPrepareSQLStatementError
	}

	defer bdStatement.Close()

	result, err := bdStatement.Exec(account.TypeID, account.Name, account.Description, account.Active)

	if err != nil {
		log.WithError(err).Error(ErrSQLExecutionError)
		return account, ErrSQLExecutionError
	}
	newId, err := result.LastInsertId()
	if err != nil {
		log.WithError(err).Error(err.Error())
		return account, ErrSQLInsertError
	}
	account.AccountID = newId
	return account, nil
}
