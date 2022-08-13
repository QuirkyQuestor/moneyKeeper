package category

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

func GetAllCategories(DBConnection *sql.DB) ([]datamodel.Category, error) {
	var categories []datamodel.Category
	bdStatement, err := DBConnection.Prepare("SELECT categoryId, parentId, name, description, expence FROM category")
	if err != nil {
		log.Fatal(err)
	}
	defer bdStatement.Close()
	rows, err := bdStatement.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var categoryId int
		var parentId int
		var name string
		var description sql.NullString
		var expence bool
		err := rows.Scan(&categoryId, &parentId, &name, &description)

		if err != nil {
			log.Fatal(err)
		}
		category := datamodel.Category{
			CategoryID:  int64(categoryId),
			Name:        name,
			ParentID:    int64(parentId),
			Description: description.String,
			Expence:     expence,
		}
		categories = append(categories, category)
		// log.Println(id, name)

	}
	return categories, nil
}
func AddCategory(DBConnection *sql.DB, category datamodel.Category) (datamodel.Category, error) {
	log.WithField("account", category).Info("The AccountType object")
	bdStatement, err := DBConnection.Prepare("INSERT INTO category(parentId, name, description, expence) VALUES (?, ?, ?, ?);")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatementError)
		return category, ErrCannotPrepareSQLStatementError
	}

	defer bdStatement.Close()

	result, err := bdStatement.Exec(category.ParentID, category.Name, category.Description, category.Expence)

	if err != nil {
		log.WithError(err).Error(ErrSQLExecutionError)
		return category, ErrSQLExecutionError
	}
	newId, err := result.LastInsertId()
	if err != nil {
		log.WithError(err).Error(err.Error())
		return category, ErrSQLInsertError
	}
	category.CategoryID = newId
	return category, nil
}
