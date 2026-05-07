package category

import (
	"database/sql"
	"errors"

	"github.com/QuirkyQuestor/moneyKeeper/internal/datamodel"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler"
	"github.com/lib/pq"
	"log/slog"
)

var (
	ErrSQLExecution              = errors.New("error during sql stetement execution")
	ErrSQLInsert                 = errors.New("error when getting LastInsertId")
	ErrCannotPrepareSQLStatement = errors.New("cannot prepare sql statement")
	ErrConvertingDBResponse      = errors.New("error during converting DB/Go types")
	ErrSQLUpdate                 = errors.New("error while updating the record in DB")
	ErrNoItemResponse            = errors.New("DB query returned no result")
)

func GetAllCategories(DBConnection *sql.DB, userID string) ([]datamodel.Category, error) {
	var categories = []datamodel.Category{}
	query := `
		WITH RECURSIVE category_path AS (
			SELECT category_id, parent_id, name, name::TEXT AS full_name, description, expence, user_id
			FROM category
			WHERE parent_id IS NULL AND user_id = $1
			UNION ALL
			SELECT c.category_id, c.parent_id, c.name, cp.full_name || ': ' || c.name, c.description, c.expence, c.user_id
			FROM category c
			JOIN category_path cp ON c.parent_id = cp.category_id
			WHERE c.user_id = $1
		)
		SELECT category_id, parent_id, name, full_name, description, expence
		FROM category_path
		ORDER BY full_name;
	`
	bdStatement, err := DBConnection.Prepare(query)
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return categories, ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()
	rows, err := bdStatement.Query(userID)
	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return nil, ErrSQLExecution
	}
	defer rows.Close()

	for rows.Next() {
		var categoryId *string
		var parentId *string
		var name string
		var fullName string
		var description sql.NullString
		var expence bool
		err := rows.Scan(&categoryId, &parentId, &name, &fullName, &description, &expence)
		if err != nil {
			slog.Error("Could not parse row from the DB", "error", err)
			continue
		}
		category := datamodel.Category{
			CategoryID:  categoryId,
			ParentID:    parentId,
			Name:        name,
			FullName:    fullName,
			Description: description.String,
			Expence:     expence,
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func AddCategory(DBConnection *sql.DB, userID string, category *datamodel.Category) error {
	slog.Info("The Category object", "category", category)

	bdStatement, err := DBConnection.Prepare("INSERT INTO category(user_id, parent_id, name, description, expence) VALUES ($1, $2, $3, $4, $5) RETURNING category_id;")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()

	err = bdStatement.QueryRow(userID, category.ParentID, category.Name, category.Description, category.Expence).Scan(&category.CategoryID)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == sqlhandler.PGErrUniqueViolation {
				return sqlhandler.SQLConflict
			}
		}

		slog.Error(ErrSQLExecution.Error(), "error", err)
		return ErrSQLExecution
	}
	return nil
}

func GetCategoryByID(DBConnection *sql.DB, userID string, categoryID string) (datamodel.Category, error) {

	var category datamodel.Category
	bdStatement, err := DBConnection.Prepare("SELECT category_id, parent_id, name, description, expence FROM category WHERE category_id = $1 AND user_id = $2")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return category, ErrCannotPrepareSQLStatement
	}
	defer bdStatement.Close()

	err = bdStatement.QueryRow(categoryID, userID).Scan(&category.CategoryID, &category.ParentID, &category.Name, &category.Description, &category.Expence)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.Info(ErrNoItemResponse.Error())
			return category, ErrNoItemResponse
		}
		slog.Error(ErrConvertingDBResponse.Error(), "error", err)
		return category, ErrConvertingDBResponse
	}

	return category, nil
}

func UpdateCategory(DBConnection *sql.DB, userID string, category *datamodel.Category) error {
	slog.Info("The Category object", "category", category)
	bdStatement, err := DBConnection.Prepare("UPDATE category SET parent_id=$1, name=$2, description=$3, expence=$4 WHERE category_id = $5 AND user_id = $6")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return ErrCannotPrepareSQLStatement
	}

	defer bdStatement.Close()
	result, err := bdStatement.Exec(category.ParentID, category.Name, category.Description, category.Expence, category.CategoryID, userID)

	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return ErrSQLExecution
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Cannot get rowsAffected", "error", err)
		return ErrSQLUpdate
	}
	slog.Info("rowsAffected", "rowsAffected", rowsAffected)

	if rowsAffected != 1 {
		slog.Error("The record does not seem to be updated.")
		return ErrSQLUpdate
	}

	return nil
}

func DeleteCategoryByID(DBConnection *sql.DB, userID string, categoryID string) error {

	bdStatement, err := DBConnection.Prepare("DELETE FROM category WHERE category_id = $1 AND user_id = $2")
	if err != nil {
		slog.Error(ErrCannotPrepareSQLStatement.Error(), "error", err)
		return ErrCannotPrepareSQLStatement
	}
	defer bdStatement.Close()

	result, err := bdStatement.Exec(categoryID, userID)
	if err != nil {
		slog.Error(ErrSQLExecution.Error(), "error", err)
		return ErrSQLExecution
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Cannot get rowsAffected for Delete", "error", err)
		return ErrSQLUpdate
	}

	if rowsAffected != 1 {
		slog.Info("The requested category did not exist in the DB Table.", "rowsAffected", rowsAffected)
	}

	return nil
}
