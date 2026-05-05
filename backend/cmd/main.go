package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/QuirkyQuestor/moneyKeeper/internal/datamodel"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler/account"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler/accountType"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler/category"
	"github.com/QuirkyQuestor/moneyKeeper/internal/sqlhandler/transaction"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

var DBConnection *sql.DB

func apiHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Welcome!"})
}

func getAllAccountsHandler(w http.ResponseWriter, r *http.Request) {
	accounts, err := account.GetAllAccounts(DBConnection)
	if err != nil {
		log.WithError(err).Error("An error has happened during GetAllAccounts DB query")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusOK, accounts)
}

func addNewAccountHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	log.WithField("d", d).Info("Incoming body string")
	newAccount := datamodel.Account{}
	err := d.Decode(&newAccount)
	if err != nil {
		log.WithError(err).Error("Could not Decode incoming request body")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if newAccount.TypeID == "" || newAccount.Name == "" {
		log.Error("Account object missing mandatory fields")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	newAccount, err = account.AddAccount(DBConnection, newAccount)
	if err != nil {
		if err == sqlhandler.SQLConflict {
			log.WithError(err).Error("Account with this name already exists")
			respondWithError(w, http.StatusConflict, http.StatusText(http.StatusConflict))
			return
		}
		log.WithError(err).Error("AddAccount DB operation returned an error")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusCreated, newAccount)
}

// Get account from DB by the account ID
func getAccountByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	log.WithField("id", idParam).Info("accountID to process")

	getAccount, err := account.GetAccountByID(DBConnection, idParam)
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
	log.WithField("id", idParam).Info("accountID to process")

	d := json.NewDecoder(r.Body)
	accountUpd := &datamodel.Account{}
	err := d.Decode(accountUpd)
	if err != nil {
		log.WithError(err).Error("Could not decode incoming json body to the Account type")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	accountUpd.AccountID = idParam

	err = account.UpdateAccountByID(DBConnection, accountUpd)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	respondWithJSON(w, http.StatusOK, accountUpd)
}

func deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	log.WithField("id", idParam).Info("accountID to process")

	err := account.DeleteAccountByID(DBConnection, idParam)
	if err != nil {
		log.WithError(err).WithField("accountID", idParam).Error("Could not delete the Account")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}

// getAllaccountTypesHandler: Get all accountType from DB
func getAllaccountTypesHandler(w http.ResponseWriter, r *http.Request) {
	accountTypes, err := accountType.GetAllAccountTypes(DBConnection)
	if err != nil {
		log.WithError(err).Error("An error has happened during GetAllAccountTypes DB query")
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
		log.WithError(err).Error("Could not Decode incoming request body")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if newAccountType.Name == "" {
		log.Error("AccountType object missing mandatory fields")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	accountType, err := accountType.AddAccountType(DBConnection, newAccountType)
	if err != nil {
		if err == sqlhandler.SQLConflict {
			log.WithError(err).Error("AccountType with this name already exists")
			respondWithError(w, http.StatusConflict, http.StatusText(http.StatusConflict))
			return
		}
		log.WithError(err).Error("AddAccountType DB operation returned an error")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusCreated, accountType)
}

func getAccountTypeByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	log.WithField("id", idParam).Info("Processing accountTypeID")

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
	log.WithField("id", idParam).Info("Processing accountTypeID")

	d := json.NewDecoder(r.Body)
	accountTypeUpd := &datamodel.AccountType{}
	err := d.Decode(accountTypeUpd)
	if err != nil {
		log.WithError(err).Error("Could not decode incoming json body to the AccountType type")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	accountTypeUpd.TypeID = idParam

	err = accountType.UpdateAccountTypeByID(DBConnection, accountTypeUpd)
	if err != nil {
		if err == sqlhandler.SQLConflict {
			log.WithError(err).Error("AddAccountType with this name already exists")
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
	log.WithField("id", idParam).Info("Processing accountTypeID")

	err := accountType.DeleteAccountTypeByID(DBConnection, idParam)
	if err != nil {
		log.WithError(err).WithField("accountTypeID", idParam).Error("Could not delete the AccountType")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// Get All Categories from DB
func getAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := category.GetAllCategories(DBConnection)
	if err != nil {
		log.WithError(err).Error("Could not do GetAllCategories")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

func addCategoryHandler(w http.ResponseWriter, r *http.Request) {
	body := new(bytes.Buffer)
	_, err := body.ReadFrom(r.Body)
	if err != nil {
		log.WithError(err).Error("Could not read the response body")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	newCategory := datamodel.Category{}
	err = json.Unmarshal(body.Bytes(), &newCategory)
	if err != nil {
		log.WithError(err).WithField("IncominBody", body).Error("Could not parse the incomming request body")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	err = category.AddCategory(DBConnection, &newCategory)
	if err != nil {
		if err == sqlhandler.SQLConflict {
			log.WithError(err).Error("Category with this name already exists")
			respondWithError(w, http.StatusConflict, http.StatusText(http.StatusConflict))
			return
		}
		log.WithError(err).WithField("IncominBody", body).Error("An error has happened during the AddCategory DB operation")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusCreated, newCategory)
}

// Get category from DB by the category ID
func getCategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	log.WithField("id", idParam).Info("Processing categoryID")

	getCategory, err := category.GetCategoryByID(DBConnection, idParam)

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
	log.WithField("id", idParam).Info("Processing categoryID")

	d := json.NewDecoder(r.Body)
	categoryUpd := &datamodel.Category{}
	err := d.Decode(categoryUpd)
	if err != nil {
		log.WithError(err).WithField("IncominBody", r.Body).Error("Could not do Decode the response body")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	*categoryUpd.CategoryID = idParam
	if err != nil {
		log.WithError(err).WithField("categoryUpd", categoryUpd).Error("Could not parse the Category ID")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	err = category.UpdateCategory(DBConnection, categoryUpd)

	if err != nil {
		log.WithError(err).WithField("categoryUpd", categoryUpd).Error("could not update the Category")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusOK, &categoryUpd)
}

func deleteCategorybyIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	log.WithField("id", idParam).Info("Processing categoryID")

	err := category.DeleteCategoryByID(DBConnection, idParam)
	if err != nil {
		log.WithError(err).WithField("categoryID", idParam).Error("could not delete the Category")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func getTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	if accountFrom, ok := queries["accountFrom"]; ok {
		// Check here if `accountFrom` is in valid format
		var re = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

		if !re.Match([]byte(accountFrom[0])) {
			log.WithField("accountFrom", accountFrom[0]).Error("accountFrom has incorrect format")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		transactions, err := transaction.GetTransactionsByAccountId(DBConnection, accountFrom[0])
		if err != nil {
			log.WithField("accountFrom", accountFrom[0]).WithError(err).Error("Could not do GetTransactionByID")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		respondWithJSON(w, http.StatusOK, transactions)
		return

	} else {
		// Get transaction from DB
		transactions, err := transaction.GetAllTransactions(DBConnection)
		if err != nil {
			log.WithError(err).Error("Could not do GetAllTransactions")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		respondWithJSON(w, http.StatusOK, transactions)
		return

	}
}

func addTransactionHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	newTransaction := datamodel.Transaction{}
	err := d.Decode(&newTransaction)
	if err != nil {
		log.WithError(err).WithField("IncominBody", r.Body).Error("Could not do Decode the response body")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	log.WithField("IncominBody", r.Body).Info("Need the response body")

	if newTransaction.AccountFrom == "" ||
		newTransaction.AccountTo == "" ||
		newTransaction.CategoryID == "" ||
		newTransaction.Amount == 0 ||
		newTransaction.Date == nil {
		log.Error("Transaction object missing mandatory fields")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	newTransaction, err = transaction.AddTransaction(DBConnection, newTransaction)
	if err != nil {
		log.WithError(err).WithField("IncominBody", r.Body).Error("An error has happened during the AddTransaction DB operation")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondWithJSON(w, http.StatusCreated, newTransaction)
}

func getTransactionByIDHandler(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")
	log.WithField("id", idParam).Info("Processing transactionID")

	getTransaction, err := transaction.GetTransactionByID(DBConnection, idParam)

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
	log.WithField("id", idParam).Info("Processing transactionID")

	d := json.NewDecoder(r.Body)
	transactionUpd := &datamodel.Transaction{}
	err := d.Decode(transactionUpd)
	if err != nil {
		log.WithError(err).WithField("IncominBody", r.Body).Error("Could not do Decode the response body")
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	transactionUpd.TransactionID = idParam

	err = transaction.UpdateTransaction(DBConnection, transactionUpd)

	if err != nil {
		log.WithError(err).WithField("transactionUpd", transactionUpd).Error("Could not update the Transaction")
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
	log.WithField("id", idParam).Info("Processing transactionID")

	err := transaction.DeleteTransactionByID(DBConnection, idParam)
	if err != nil {
		log.WithError(err).WithField("transactionID", idParam).Error("Could not delete the transaction")
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func main() {
	DBConnection = sqlhandler.DBConnect()
	defer DBConnection.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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

	log.Println("Starting MoneyKeeper backend!")
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func respondWithJSON(response http.ResponseWriter, statusCode int, data interface{}) {
	result, err := json.Marshal(data)
	if err != nil {
		log.WithError(err).Error("Could not marshal data object for response")
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	response.Write(result)
}

func respondWithError(response http.ResponseWriter, statusCode int, msg string) {
	respondWithJSON(response, statusCode, map[string]string{"error": msg})
}
