package account

import (
	"database/sql"
	"errors"

	"github.com/QuirkyQuestor/moneyKeeper/internal/datamodel"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
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
	var accounts = []datamodel.Account{}

	bdStatement, err := DBConnection.Prepare("SELECT account_id, type_id, name, description, active FROM moneykeeper.account;")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return accounts, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	rows, err := bdStatement.Query()

	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
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
	log.WithField("incomming_account", account).Info("Received Account object")

	bdStatement, err := DBConnection.Prepare("INSERT INTO moneykeeper.account(type_id, name, description, active) VALUES ($1, $2, $3, $4) RETURNING account_id;")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return account, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	err = bdStatement.QueryRow(account.TypeID, account.Name, account.Description, account.Active).Scan(&account.AccountID)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == sqlhandler.PGErrUniqueViolation {
				return account, sqlhandler.SQLConflict
			}
		}
		log.WithError(err).Error(ErrSQLExecution)
		return account, ErrSQLExecution
	}
	return account, nil
}

func GetAccountByID(DBConnection *sql.DB, accountID string) (datamodel.Account, error) {
	var account datamodel.Account

	bdStatement, err := DBConnection.Prepare("SELECT account_id, type_id, name, description, active FROM moneykeeper.account WHERE account_id = $1;")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return account, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	var accountId string
	var typeId string
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

	bdStatement, err := DBConnection.Prepare("UPDATE moneykeeper.account SET type_id=$1, name=$2, description=$3, active=$4 WHERE account_id = $5")
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

	bdStatement, err := DBConnection.Prepare("DELETE FROM moneykeeper.account WHERE account_id = $1")
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
