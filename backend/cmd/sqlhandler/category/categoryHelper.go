package category

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/winchien/moneyKeeper/backend/cmd/datamodel"
	"github.com/winchien/moneyKeeper/backend/cmd/sqlhandler"
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
	var categories = []datamodel.Category{}
	bdStatement, err := DBConnection.Prepare("SELECT category_id, parent_id, name, description, expence FROM moneykeeper.category")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return categories, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()
	rows, err := bdStatement.Query()
	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return nil, ErrSQLExecution
	}
	defer rows.Close()

	for rows.Next() {
		var categoryId string
		var parentId *string
		var name string
		var description sql.NullString
		var expence bool
		err := rows.Scan(&categoryId, &parentId, &name, &description, &expence)
		if err != nil {
			log.WithError(err).Error("Cold not parse row from the DB")
		}
		category := datamodel.Category{
			CategoryID:  categoryId,
			ParentID:    parentId,
			Name:        name,
			Description: description.String,
			Expence:     expence,
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func AddCategory(DBConnection *sql.DB, category *datamodel.Category) error {
	log.WithField("category", category).Info("The Category object")
	bdStatement, err := DBConnection.Prepare("INSERT INTO moneykeeper.category(parent_id, name, description, expence) VALUES ($1, $2, $3, $4) RETURNING category_id;")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	err = bdStatement.QueryRow(category.ParentID, category.Name, category.Description, category.Expence).Scan(&category.CategoryID)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == sqlhandler.PGErrUniqueViolation {
				return sqlhandler.SQLConflict
			}
		}

		log.WithError(err).Error(ErrSQLExecution)
		return ErrSQLExecution
	}
	return nil
}

func GetCategoryByID(DBConnection *sql.DB, categoryID string) (datamodel.Category, error) {

	var category datamodel.Category
	bdStatement, err := DBConnection.Prepare("SELECT category_id, parent_id, name, description, expence FROM moneykeeper.category WHERE category_id = $1")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return category, ErrCannotPrepareSQLStatement
	}
	defer bdStatement.Close()

	err = bdStatement.QueryRow(categoryID).Scan(&category.CategoryID, &category.ParentID, &category.Name, &category.Description, &category.Expence)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Info(ErrNoItemResponse)
			return category, ErrNoItemResponse
		}
		log.WithError(err).Error(ErrConvertingDBResponse)
		return category, ErrConvertingDBResponse
	}

	return category, nil
}

func UpdateCategory(DBConnection *sql.DB, category *datamodel.Category) error {
	log.WithField("category", category).Info("The Category object")
	bdStatement, err := DBConnection.Prepare("UPDATE moneykeeper.category SET parent_id=$1, name=$2, description=$3, expence=$4 WHERE category_id = $5")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()
	result, err := bdStatement.Exec(category.ParentID, category.Name, category.Description, category.Expence, category.CategoryID)

	if err != nil {
		log.WithError(err).Error(ErrSQLExecution)
		return ErrSQLExecution
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.WithError(err).Error("Cannot get rowsAffected")
		return ErrSQLUpdate
	}
	log.WithField("rowsAffected", rowsAffected).Info("rowsAffected")

	if rowsAffected != 1 {
		log.Error("The record does not seem to be updated.")
		return ErrSQLUpdate
	}

	return nil
}

func DeleteCategoryByID(DBConnection *sql.DB, categoryID string) error {

	bdStatement, err := DBConnection.Prepare("DELETE FROM moneykeeper.category WHERE category_id = $1")
	if err != nil {
		log.WithError(err).Error(ErrCannotPrepareSQLStatement)
		return ErrCannotPrepareSQLStatement
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

	if rowsAffected != 1 {
		log.WithField("rowsAffected", rowsAffected).Info("The requested category did not exist in the DB Table.")
	}

	return nil
}
