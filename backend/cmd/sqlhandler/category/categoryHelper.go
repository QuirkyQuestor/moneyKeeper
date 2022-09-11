package category

import (
	"database/sql"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/winchien/moneyKeeper/backend/cmd/datamodel"
)

var (
	ErrSQLExecution              = errors.New("error during sql stetement execution")
	ErrSQLInsert                 = errors.New("error when getting LastInsertId")
	ErrCannotPrepareSQLStatement = errors.New("cannot prepare sql statement")
	ErrConvertingDBResponse      = errors.New("error during converting DB/Go types")
	ErrSQLUpdate                 = errors.New("error while updating the record in DB")
	ErrNoItemResponse            = errors.New("DB query returned no result")
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
			log.WithError(err).Error("Cold not parse row from the DB")
		}
		category := datamodel.Category{
			CategoryID:  int64(categoryId),
			Name:        name,
			ParentID:    int64(parentId),
			Description: description.String,
			Expence:     expence,
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func AddCategory(DBConnection *sql.DB, category datamodel.Category) (datamodel.Category, error) {
	log.WithField("category", category).Info("The Category object")
	bdStatement, err := DBConnection.Prepare("INSERT INTO category(parentId, name, description, expence) VALUES (?, ?, ?, ?);")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return category, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	result, err := bdStatement.Exec(category.ParentID, category.Name, category.Description, category.Expence)

	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return category, ErrSQLExecution
	}
	newId, err := result.LastInsertId()
	if err != nil {
		log.WithError(err).Error(err.Error())
		return category, ErrSQLInsert
	}
	category.CategoryID = newId
	return category, nil
}

func GetCategoryByID(DBConnection *sql.DB, categoryID string) (*datamodel.Category, error) {

	var category datamodel.Category
	bdStatement, err := DBConnection.Prepare("SELECT categoryId, parentId, name, description, expence FROM category WHERE categoryID = ?")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return nil, ErrCannotPrepareSQLStatement
	}
	defer bdStatement.Close()

	var categoryId int
	var parentId int
	var name string
	var description sql.NullString
	var expence bool

	err = bdStatement.QueryRow(categoryID).Scan(&categoryId, &parentId, &name, &description, &expence)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Info(ErrNoItemResponse)
			return nil, ErrNoItemResponse
		}
		log.WithError(err).Error(ErrConvertingDBResponse)
		return nil, ErrConvertingDBResponse
	}

	category = datamodel.Category{
		CategoryID:  int64(categoryId),
		ParentID:    int64(parentId),
		Name:        name,
		Description: description.String,
		Expence:     expence,
	}
	return &category, nil
}

func UpdateCategory(DBConnection *sql.DB, category *datamodel.Category) error {
	log.WithField("category", category).Info("The Category object")
	bdStatement, err := DBConnection.Prepare("UPDATE category SET parentID=?, name=?, description=?, expence=? WHERE categoryId = ?")
	if err != nil {
		log.WithError(err).Error("cannot prepare update statement")
	}

	defer bdStatement.Close()
	result, err := bdStatement.Exec(category.ParentID, category.Name, category.Description, category.Expence, category.CategoryID)

	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return ErrSQLExecution
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithError(err).Error("Cannot get rowsAffected for Delete")
		return ErrSQLUpdate
	}
	log.WithField("rowsAffected", rowsAffected).Info("rowsAffected")

	if rowsAffected != 1 {
		log.Error("The record does not seem to be updated.")
		return ErrSQLUpdate
	}

	return nil
}
func DeleteCategoryByID(DBConnection *sql.DB, categoryID int64) error {

	bdStatement, err := DBConnection.Prepare("DELETE FROM category WHERE categoryId = ?")
	if err != nil {
		log.WithError(err).Fatal("Cannot prepare SQL statement")
	}
	defer bdStatement.Close()

	result, err := bdStatement.Exec(categoryID)
	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return ErrSQLExecution
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithError(err).Error("Cannot get rowsAffected for Delete")
		return ErrSQLUpdate
	}
	log.WithField("rowsAffected", rowsAffected).Info("rowsAffected")

	if rowsAffected != 1 {
		log.Error("The record does not seem to be updated.")
		return ErrSQLUpdate
	}

	return nil
}
