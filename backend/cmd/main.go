package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/QuirkyQuestor/moneyKeeper/internal/auth"
	"github.com/QuirkyQuestor/moneyKeeper/internal/datamodel"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler/account"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler/accountType"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler/category"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler/reports"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler/transaction"
	"github.com/lib/pq"
	"log/slog"
)

func init() {
	// Configure slog to use JSON handler and output to stdout.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

var DBConnection *sql.DB

func apiHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Welcome!"})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	hash, err := auth.HashPassword(creds.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password")
		return
	}

	newUserID, err := uuid.NewV7()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate user ID")
		return
	}

	var userID string
	err = DBConnection.QueryRow("INSERT INTO users(user_id, email, password_hash) VALUES($1, $2, $3) RETURNING user_id", newUserID, creds.Email, hash).Scan(&userID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == sqlhandler.PGErrUniqueViolation {
				respondWithError(w, http.StatusConflict, "Email already exists")
				return
			}
		}
		slog.Error("Database error during registration", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Registration failed due to server error")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"user_id": userID})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		slog.Error("Could not decode login body", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	slog.Info("Login attempt received", "email", creds.Email)

	var userID string
	var hash string
	err := DBConnection.QueryRow("SELECT user_id, password_hash FROM users WHERE email = $1", creds.Email).Scan(&userID, &hash)
	if err != nil {
		slog.Warn("Login failed: user not found", "email", creds.Email, "error", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !auth.CheckPasswordHash(creds.Password, hash) {
		slog.Warn("Login failed: incorrect password", "email", creds.Email)
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, expiresAt, err := auth.GenerateJWT(userID)
	if err != nil {
		slog.Error("Could not generate JWT", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Could not generate token")
		return
	}

	slog.Info("Login successful, setting cookie", "userID", userID)
	auth.SetAuthCookie(w, token, expiresAt)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Logged in successfully"})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	var email string
	err := DBConnection.QueryRow("SELECT email FROM users WHERE user_id = $1", userID).Scan(&email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"user_id": userID, "email": email})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		userID, err := auth.GetUserIDFromToken(cookie.Value)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAllAccountsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	accounts, err := account.GetAllAccounts(DBConnection, userID)
	if err != nil {
		slog.Error("An error has happened during GetAllAccounts DB query", "error", err)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusOK, accounts)
}

func addNewAccountHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	slog.Info("Incoming body string", "d", d)
	newAccount := datamodel.Account{}
	err := d.Decode(&newAccount)
	if err != nil {
		slog.Error("Could not Decode incoming request body", "error", err)
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if newAccount.TypeID == "" || newAccount.Name == "" {
		slog.Error("Account object missing mandatory fields")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	userID := r.Context().Value("user_id").(string)
	newAccount, err = account.AddAccount(DBConnection, userID, newAccount)
	if err != nil {
		if err == sqlhandler.SQLConflict {
			slog.Error("Account with this name already exists", "error", err)
			respondWithError(w, http.StatusConflict, http.StatusText(http.StatusConflict))
			return
		}
		slog.Error("AddAccount DB operation returned an error", "error", err)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusCreated, newAccount)
}

// Get account from DB by the account ID
func getAccountByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	slog.Info("accountID to process", "id", idParam)
	userID := r.Context().Value("user_id").(string)

	getAccount, err := account.GetAccountByID(DBConnection, userID, idParam)
	if err == account.ErrNoItemResponse {
		respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusOK, getAccount)
}

func updateAccountHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	slog.Info("accountID to process", "id", idParam)

	d := json.NewDecoder(r.Body)
	accountUpd := &datamodel.Account{}
	err := d.Decode(accountUpd)
	if err != nil {
		slog.Error("Could not decode incoming json body to the Account type", "error", err)
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	accountUpd.AccountID = idParam
	userID := r.Context().Value("user_id").(string)

	err = account.UpdateAccountByID(DBConnection, userID, accountUpd)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	respondWithJSON(w, http.StatusOK, accountUpd)
}

func deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	slog.Info("accountID to process", "id", idParam)
	userID := r.Context().Value("user_id").(string)

	err := account.DeleteAccountByID(DBConnection, userID, idParam)
	if err != nil {
		slog.Error("Could not delete the Account", "error", err, "accountID", idParam)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}

// getAllaccountTypesHandler: Get all accountType from DB
func getAllaccountTypesHandler(w http.ResponseWriter, r *http.Request) {
	accountTypes, err := accountType.GetAllAccountTypes(DBConnection)
	if err != nil {
		slog.Error("An error has happened during GetAllAccountTypes DB query", "error", err)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusOK, accountTypes)
}

func addAccountTypeHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	newAccountType := datamodel.AccountType{}
	err := d.Decode(&newAccountType)
	if err != nil {
		slog.Error("Could not Decode incoming request body", "error", err)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if newAccountType.Name == "" {
		slog.Error("AccountType object missing mandatory fields")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	accountType, err := accountType.AddAccountType(DBConnection, newAccountType)
	if err != nil {
		if err == sqlhandler.SQLConflict {
			slog.Error("AccountType with this name already exists", "error", err)
			respondWithError(w, http.StatusConflict, http.StatusText(http.StatusConflict))
			return
		}
		slog.Error("AddAccountType DB operation returned an error", "error", err)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusCreated, accountType)
}

func getAccountTypeByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	slog.Info("Processing accountTypeID", "id", idParam)

	getAccountType, err := accountType.GetAccountTypeByID(DBConnection, idParam)
	if err == accountType.ErrNoItemResponse {
		respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusOK, getAccountType)
}

func updateAccountTypeByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	slog.Info("Processing accountTypeID", "id", idParam)

	d := json.NewDecoder(r.Body)
	accountTypeUpd := &datamodel.AccountType{}
	err := d.Decode(accountTypeUpd)
	if err != nil {
		slog.Error("Could not decode incoming json body to the AccountType type", "error", err)
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	accountTypeUpd.TypeID = idParam

	err = accountType.UpdateAccountTypeByID(DBConnection, accountTypeUpd)
	if err != nil {
		if err == sqlhandler.SQLConflict {
			slog.Error("AddAccountType with this name already exists", "error", err)
			respondWithError(w, http.StatusConflict, http.StatusText(http.StatusConflict))
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusOK, accountTypeUpd)
}

func deleteAccountTypeByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	slog.Info("Processing accountTypeID", "id", idParam)

	err := accountType.DeleteAccountTypeByID(DBConnection, idParam)
	if err != nil {
		slog.Error("Could not delete the AccountType", "error", err, "accountTypeID", idParam)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// Get All Categories from DB
func getAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	categories, err := category.GetAllCategories(DBConnection, userID)
	if err != nil {
		slog.Error("Could not do GetAllCategories", "error", err)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

func addCategoryHandler(w http.ResponseWriter, r *http.Request) {
	body := new(bytes.Buffer)
	_, err := body.ReadFrom(r.Body)
	if err != nil {
		slog.Error("Could not read the response body", "error", err)
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	newCategory := datamodel.Category{}
	err = json.Unmarshal(body.Bytes(), &newCategory)
	if err != nil {
		slog.Error("Could not parse the incomming request body", "error", err, "IncominBody", body)
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	userID := r.Context().Value("user_id").(string)
	err = category.AddCategory(DBConnection, userID, &newCategory)
	if err != nil {
		if err == sqlhandler.SQLConflict {
			slog.Error("Category with this name already exists", "error", err)
			respondWithError(w, http.StatusConflict, http.StatusText(http.StatusConflict))
			return
		}
		slog.Error("An error has happened during the AddCategory DB operation", "error", err, "IncominBody", body)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusCreated, newCategory)
}

// Get category from DB by the category ID
func getCategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	slog.Info("Processing categoryID", "id", idParam)
	userID := r.Context().Value("user_id").(string)

	getCategory, err := category.GetCategoryByID(DBConnection, userID, idParam)

	if err != nil {
		if err == category.ErrNoItemResponse {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}

	respondWithJSON(w, http.StatusOK, getCategory)
}

func updateCategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	slog.Info("Processing categoryID", "id", idParam)

	d := json.NewDecoder(r.Body)
	categoryUpd := &datamodel.Category{}
	err := d.Decode(categoryUpd)
	if err != nil {
		slog.Error("Could not do Decode the response body", "error", err, "IncominBody", r.Body)
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	*categoryUpd.CategoryID = idParam
	userID := r.Context().Value("user_id").(string)
	err = category.UpdateCategory(DBConnection, userID, categoryUpd)

	if err != nil {
		slog.Error("could not update the Category", "error", err, "categoryUpd", categoryUpd)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusOK, &categoryUpd)
}

func deleteCategorybyIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	slog.Info("Processing categoryID", "id", idParam)
	userID := r.Context().Value("user_id").(string)

	err := category.DeleteCategoryByID(DBConnection, userID, idParam)
	if err != nil {
		slog.Error("could not delete the Category", "error", err, "categoryID", idParam)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func getTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	userID := r.Context().Value("user_id").(string)

	limitStr := queries.Get("limit")
	offsetStr := queries.Get("offset")

	limit := 50
	offset := 0

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	type paginatedResponse struct {
		Transactions []datamodel.Transaction `json:"transactions"`
		TotalCount   int                    `json:"totalCount"`
	}

	if accountFrom, ok := queries["accountFrom"]; ok && accountFrom[0] != "" {
		// Check here if `accountFrom` is in valid format
		var re = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

		if !re.Match([]byte(accountFrom[0])) {
			slog.Error("accountFrom has incorrect format", "accountFrom", accountFrom[0])
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		transactions, totalCount, err := transaction.GetTransactionsByAccountId(DBConnection, userID, accountFrom[0], limit, offset)
		if err != nil {
			slog.Error("Could not do GetTransactionByID", "accountFrom", accountFrom[0], "error", err)
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		respondWithJSON(w, http.StatusOK, paginatedResponse{Transactions: transactions, TotalCount: totalCount})
		return

	} else {
		// Get transaction from DB
		transactions, totalCount, err := transaction.GetAllTransactions(DBConnection, userID, limit, offset)
		if err != nil {
			slog.Error("Could not do GetAllTransactions", "error", err)
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		respondWithJSON(w, http.StatusOK, paginatedResponse{Transactions: transactions, TotalCount: totalCount})
		return

	}
}

func addTransactionHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	newTransaction := datamodel.Transaction{}
	err := d.Decode(&newTransaction)
	if err != nil {
		slog.Error("Could not do Decode the response body", "error", err, "IncominBody", r.Body)
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	slog.Info("Need the response body", "IncominBody", r.Body)

	if newTransaction.AccountFrom == "" ||
		newTransaction.CategoryID == "" ||
		newTransaction.Amount == 0 ||
		newTransaction.Date == nil {
		slog.Error("Transaction object missing mandatory fields")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	userID := r.Context().Value("user_id").(string)
	newTransaction, err = transaction.AddTransaction(DBConnection, userID, newTransaction)
	if err != nil {
		slog.Error("An error has happened during the AddTransaction DB operation", "error", err, "IncominBody", r.Body)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusCreated, newTransaction)
}

func getTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")
	slog.Info("Processing transactionID", "id", idParam)
	userID := r.Context().Value("user_id").(string)

	getTransaction, err := transaction.GetTransactionByID(DBConnection, userID, idParam)

	if err != nil {
		if err == transaction.ErrNoItemResponse {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}

	respondWithJSON(w, http.StatusOK, getTransaction)
}

func updateTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")
	slog.Info("Processing transactionID", "id", idParam)

	d := json.NewDecoder(r.Body)
	transactionUpd := &datamodel.Transaction{}
	err := d.Decode(transactionUpd)
	if err != nil {
		slog.Error("Could not do Decode the response body", "error", err, "IncominBody", r.Body)
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	transactionUpd.TransactionID = idParam
	userID := r.Context().Value("user_id").(string)

	err = transaction.UpdateTransaction(DBConnection, userID, transactionUpd)

	if err != nil {
		slog.Error("Could not update the Transaction", "error", err, "transactionUpd", transactionUpd)
		if err == transaction.ErrNoItemResponse {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusOK, &transactionUpd)
}

func deleteTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")
	slog.Info("Processing transactionID", "id", idParam)
	userID := r.Context().Value("user_id").(string)

	err := transaction.DeleteTransactionByID(DBConnection, userID, idParam)
	if err != nil {
		slog.Error("Could not delete the transaction", "error", err, "transactionID", idParam)
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func getReportsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	expenses, err := reports.GetExpensesByCategory(DBConnection, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get category reports")
		return
	}

	monthly, err := reports.GetMonthlyIncomeVsExpenses(DBConnection, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get monthly reports")
		return
	}

	networth, err := reports.GetNetWorthTrend(DBConnection, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get net worth reports")
		return
	}

	summary := reports.ReportsSummary{
		ExpensesByCategory: expenses,
		MonthlyComparison:  monthly,
		NetWorthTrend:      networth,
	}

	respondWithJSON(w, http.StatusOK, summary)
}

func main() {
	DBConnection = sqlhandler.DBConnect()
	defer DBConnection.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))

	})
	r.Get("/", apiHandler)

	// Auth routes
	r.Post("/api/register", registerHandler)
	r.Post("/api/login", loginHandler)
	r.Post("/api/logout", logoutHandler)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Get("/api/me", meHandler)

		r.Get("/api/account", getAllAccountsHandler)                                                                                    // Handle Accout requests: GET ALL accounts
		r.Post("/api/account", addNewAccountHandler)                                                                                    // Handle Accout requests: POST new account
		r.Get("/api/account/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", getAccountByIDHandler)   // Handle Account requests: GET single specified account
		r.Put("/api/account/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", updateAccountHandler)    // Handle Account requests: PUT to update account
		r.Delete("/api/account/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", deleteAccountHandler) // Handle Account requests: DELETE specified account

		r.Get("/api/account_type", getAllaccountTypesHandler)                                                                                        // Handle AccountType requests: GET all account type
		r.Post("/api/account_type", addAccountTypeHandler)                                                                                           // Handle AccountType requests: POST new account type
		r.Get("/api/account_type/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", getAccountTypeByIDHandler)       // Handle AccountType requests: GET single specified account type
		r.Put("/api/account_type/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", updateAccountTypeByIDHandler)    // Handle AccountType requests: PUT to update account type
		r.Delete("/api/account_type/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", deleteAccountTypeByIDHandler) // Handle AccountType requests: DELETE specified account type

		r.Get("/api/category", getAllCategoriesHandler)                                                                                       // GET ALL categories
		r.Post("/api/category", addCategoryHandler)                                                                                           // POST new category
		r.Get("/api/category/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", getCategoryByIDHandler)       // GET single category info
		r.Put("/api/category/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", updateCategoryByIDHandler)    // PUT update category
		r.Delete("/api/category/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", deleteCategorybyIDHandler) //  DELETE category

		r.Get("/api/transaction", getTransactionsHandler)                                                                                           // GET ALL transactions, SEARCH transactions by some criteria
		r.Post("/api/transaction", addTransactionHandler)                                                                                           // POST new transaction
		r.Get("/api/transaction/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", getTransactionByIDHandler)       // GET single transaction info
		r.Put("/api/transaction/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", updateTransactionByIDHandler)    // PUT update transaction
		r.Delete("/api/transaction/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", deleteTransactionByIDHandler) // DELETE transaction

		r.Get("/api/reports", getReportsHandler)
	})

	slog.Info("Starting MoneyKeeper backend!")
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func respondWithJSON(response http.ResponseWriter, statusCode int, data interface{}) {
	result, err := json.Marshal(data)
	if err != nil {
		slog.Error("Could not marshal data object for response", "error", err)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	response.Write(result)
}

func respondWithError(response http.ResponseWriter, statusCode int, msg string) {
	respondWithJSON(response, statusCode, map[string]string{"error": msg})
}
