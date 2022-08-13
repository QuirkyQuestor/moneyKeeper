package accountType

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

func GetAllAccountTypes(DBConnection *sql.DB) ([]datamodel.AccountType, error) {
	var accountTypes []datamodel.AccountType

	bdStatement, err := DBConnection.Prepare("SELECT typeId, name, description, active FROM accountType;")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatementError)
		return accountTypes, ErrCannotPrepareSQLStatementError
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
			log.WithError(err).Error(ErrConvertingDBResponseError)
			return accountTypes, ErrConvertingDBResponseError
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
		log.WithError(err).Error(ErrCannotPrepareSQLStatementError)
		return accountType, ErrCannotPrepareSQLStatementError
	}

	defer bdStatement.Close()

	result, err := bdStatement.Exec(accountType.Name, accountType.Description)

	if err != nil {
		log.WithError(err).Error(ErrSQLExecutionError)
		return accountType, ErrSQLExecutionError
	}
	newId, err := result.LastInsertId()
	if err != nil {
		log.WithError(err).Error(err.Error())
		return accountType, ErrSQLInsertError
	}
	accountType.TypeID = newId
	return accountType, nil
}
