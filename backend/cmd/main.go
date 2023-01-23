package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/winchien/moneyKeeper/backend/cmd/datamodel"
	"github.com/winchien/moneyKeeper/backend/cmd/sqlhandler"
	"github.com/winchien/moneyKeeper/backend/cmd/sqlhandler/account"
	"github.com/winchien/moneyKeeper/backend/cmd/sqlhandler/accountType"
	"github.com/winchien/moneyKeeper/backend/cmd/sqlhandler/category"
	"github.com/winchien/moneyKeeper/backend/cmd/sqlhandler/transaction"
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

	switch r.Method {
	case "GET":
		w.Write([]byte("Welcome to MoneyKeeper"))
		return
	default:
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		accounts, err := account.GetAllAccounts(DBConnection)
		if err != nil {
			log.WithError(err).Error("An error has happened during GetAllAccounts DB query")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		respondWithJSON(w, http.StatusOK, accounts)
		return

	case "POST":
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
		return
	default:
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func accountIDHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	log.WithField("id", vars["id"]).Info("accountID to process")

	switch r.Method {
	case "GET":
		// Get account from DB by the account ID
		getAccount, err := account.GetAccountByID(DBConnection, vars["id"])
		if err == account.ErrNoItemResponse {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		respondWithJSON(w, http.StatusOK, getAccount)
		return

	case "PUT":
		d := json.NewDecoder(r.Body)
		accountUpd := &datamodel.Account{}
		err := d.Decode(accountUpd)
		if err != nil {
			log.WithError(err).Error("Could not decode incoming json body to the Account type")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}

		accountUpd.AccountID = vars["id"]

		err = account.UpdateAccountByID(DBConnection, accountUpd)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		respondWithJSON(w, http.StatusOK, accountUpd)
		return

	case "DELETE":
		err := account.DeleteAccountByID(DBConnection, vars["id"])
		if err != nil {
			log.WithError(err).WithField("accountID", vars["id"]).Error("Could not delete the Account")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		respondWithJSON(w, http.StatusNoContent, nil)
		return

	default:
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func accountTypeHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		// Get accountType from DB
		accountTypes, err := accountType.GetAllAccountTypes(DBConnection)
		if err != nil {
			log.WithError(err).Error("An error has happened during GetAllAccountTypes DB query")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusOK, accountTypes)
		return

	case "POST":
		// Decode the JSON in the body and overwrite 'tom' with it
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
		return

	default:
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func accountTypeIDHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	log.WithField("accountTypeID", vars["id"]).Info("Processing accountTypeID")

	switch r.Method {
	case "GET":
		getAccountType, err := accountType.GetAccountTypeByID(DBConnection, vars["id"])
		if err == accountType.ErrNoItemResponse {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusOK, getAccountType)
		return

	case "PUT":
		d := json.NewDecoder(r.Body)
		accountTypeUpd := &datamodel.AccountType{}
		err := d.Decode(accountTypeUpd)
		if err != nil {
			log.WithError(err).Error("Could not decode incoming json body to the AccountType type")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		accountTypeUpd.TypeID = vars["id"]

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
		return

	case "DELETE":
		err := accountType.DeleteAccountTypeByID(DBConnection, vars["id"])

		if err != nil {
			log.WithError(err).WithField("accountTypeID", vars["id"]).Error("Could not delete the AccountType")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusNoContent, nil)
		return

	default:
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func categoryHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		// Get accountType from DB
		categories, err := category.GetAllCategories(DBConnection)
		if err != nil {
			log.WithError(err).Error("Could not do GetAllCategories")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusOK, categories)
		return

	case "POST":
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
		return

	default:
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func categoryIDHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	log.Printf("categoryID: %v", vars["id"])

	switch r.Method {
	case "GET":
		// Get category from DB
		getCategory, err := category.GetCategoryByID(DBConnection, vars["id"])

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
		return

	case "PUT":
		d := json.NewDecoder(r.Body)
		categoryUpd := &datamodel.Category{}
		err := d.Decode(categoryUpd)
		if err != nil {
			log.WithError(err).WithField("IncominBody", r.Body).Error("Could not do Decode the response body")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		*categoryUpd.CategoryID = vars["id"]
		if err != nil {
			log.WithError(err).WithField("categoryUpd", categoryUpd).Error("Could not parse the Category ID")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		err = category.UpdateCategory(DBConnection, categoryUpd)

		if err != nil {
			log.WithError(err).WithField("categoryUpd", categoryUpd).Error("Could not update the Category")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusOK, &categoryUpd)
		return

	case "DELETE":
		err := category.DeleteCategoryByID(DBConnection, vars["id"])
		if err != nil {
			log.WithError(err).WithField("categoryID", vars["id"]).Error("Could the category")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusNoContent, nil)
		return

	default:
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func transactionHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

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

	case "POST":
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
		return

	default:
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func transactionQueryHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		accountFrom := r.FormValue("accountFrom")

		log.WithField("accountFrom", accountFrom).Info("accountFrom")

		// // Get transaction from DB
		// transactions, err := transaction.GetAllTransactions(DBConnection)
		// if err != nil {
		// 	log.WithError(err).Error("Could not do GetAllTransactions")
		// 	respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		// 	return
		// }

		respondWithJSON(w, http.StatusOK, nil)
		return
	default:
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}
func transactionIDHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	log.Printf("transactionID: %v", vars["id"])

	switch r.Method {
	case "GET":
		// Get transaction from DB
		getTransaction, err := transaction.GetTransactionByID(DBConnection, vars["id"])

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
		return

	case "PUT":
		d := json.NewDecoder(r.Body)
		transactionUpd := &datamodel.Transaction{}
		err := d.Decode(transactionUpd)
		if err != nil {
			log.WithError(err).WithField("IncominBody", r.Body).Error("Could not do Decode the response body")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		transactionUpd.TransactionID = vars["id"]

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
		return

	case "DELETE":
		err := transaction.DeleteTransactionByID(DBConnection, vars["id"])
		if err != nil {
			log.WithError(err).WithField("transactionID", vars["id"]).Error("Could not delete the transaction")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		respondWithJSON(w, http.StatusNoContent, nil)
		return

	default:
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func main() {
	DBConnection = sqlhandler.DBConnect()

	r := mux.NewRouter()

	r.HandleFunc("/api/", apiHandler)                                                                                                        // GET saying Hello
	r.HandleFunc("/api/account", accountHandler)                                                                                             // Handle Accout requests: GET ALL accounts, POST new account
	r.HandleFunc("/api/account/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", accountIDHandler)          // Handle Account requests: GET single specified account type, PUT to update account type, DELETE specified account type
	r.HandleFunc("/api/account_type", accountTypeHandler)                                                                                    // Handle AccountType requests: GET all account type, POST new account type
	r.HandleFunc("/api/account_type/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", accountTypeIDHandler) // Handle AccountType requests: GET single specified account type, PUT to update account type, DELETE specified account type
	r.HandleFunc("/api/category", categoryHandler)                                                                                           // GET ALL categories, POST new category
	r.HandleFunc("/api/category/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", categoryIDHandler)        // GET single category info, PUT update category, DELETE category
	r.HandleFunc("/api/transaction", transactionHandler)                                                                                     // GET ALL transactions, POST new transaction, SEARCH transactions by some criteria
	// r.HandleFunc("/api/transaction", transactionQueryHandler).Queries("accountFrom", "{accountFrom:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}") // GET transactions for the specified accountID
	r.HandleFunc("/api/transaction/{id:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}}", transactionIDHandler) // GET single transaction info, PUT update transaction, DELETE transaction

	log.Println("Starting MoneyKeeper backend!")
	srv := &http.Server{
		Handler: r,
		Addr:    ":8000",
		// Good practice: enforce timeouts for servers you create!
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
