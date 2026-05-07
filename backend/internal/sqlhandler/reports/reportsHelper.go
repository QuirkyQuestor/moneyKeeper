package reports

import (
	"database/sql"

	"log/slog"
)

type CategoryReport struct {
	CategoryName string  `json:"categoryName"`
	Amount       float64 `json:"amount"`
}

type MonthlyReport struct {
	Month    string  `json:"month"`
	Income   float64 `json:"income"`
	Expenses float64 `json:"expenses"`
}

type NetWorthPoint struct {
	Date    string  `json:"date"`
	Balance float64 `json:"balance"`
}

func GetExpensesByCategory(DBConnection *sql.DB, userID string) ([]CategoryReport, error) {
	var reports []CategoryReport

	query := `
		SELECT c.name, SUM(ABS(t.amount)) as total_amount
		FROM transaction t
		JOIN category c ON t.category_id = c.category_id
		WHERE t.user_id = $1 AND t.amount < 0
		GROUP BY c.name
		ORDER BY total_amount DESC;
	`

	rows, err := DBConnection.Query(query, userID)
	if err != nil {
		slog.Error("Error querying expenses by category", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var report CategoryReport
		if err := rows.Scan(&report.CategoryName, &report.Amount); err != nil {
			slog.Error("Error scanning category report row", "error", err)
			continue
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func GetMonthlyIncomeVsExpenses(DBConnection *sql.DB, userID string) ([]MonthlyReport, error) {
	var reports []MonthlyReport

	query := `
		SELECT 
			TO_CHAR(date, 'YYYY-MM') as month,
			SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END) as income,
			SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END) as expenses
		FROM transaction
		WHERE user_id = $1
		GROUP BY month
		ORDER BY month ASC;
	`

	rows, err := DBConnection.Query(query, userID)
	if err != nil {
		slog.Error("Error querying monthly reports", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var report MonthlyReport
		if err := rows.Scan(&report.Month, &report.Income, &report.Expenses); err != nil {
			slog.Error("Error scanning monthly report row", "error", err)
			continue
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func GetNetWorthTrend(DBConnection *sql.DB, userID string) ([]NetWorthPoint, error) {
	var points []NetWorthPoint

	// Cumulative sum over all transactions
	query := `
		SELECT 
			TO_CHAR(date, 'YYYY-MM-DD') as day,
			SUM(amount) OVER (ORDER BY date) as balance
		FROM transaction
		WHERE user_id = $1
		ORDER BY date ASC;
	`

	rows, err := DBConnection.Query(query, userID)
	if err != nil {
		slog.Error("Error querying net worth trend", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var point NetWorthPoint
		if err := rows.Scan(&point.Date, &point.Balance); err != nil {
			slog.Error("Error scanning net worth point", "error", err)
			continue
		}
		points = append(points, point)
	}

	// For the chart, we might want to sample this if there are too many points,
	// but for a start, returning all is fine.
	return points, nil
}

type ReportsSummary struct {
	ExpensesByCategory []CategoryReport `json:"expensesByCategory"`
	MonthlyComparison  []MonthlyReport  `json:"monthlyComparison"`
	NetWorthTrend      []NetWorthPoint  `json:"netWorthTrend"`
}
