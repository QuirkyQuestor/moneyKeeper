package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/winchien/moneyKeeper/backend/cmd/datamodel"
	"github.com/winchien/moneyKeeper/backend/cmd/sqlhandler"
	"github.com/winchien/moneyKeeper/backend/cmd/sqlhandler/account"
	"github.com/winchien/moneyKeeper/backend/cmd/sqlhandler/accountType"
	"github.com/winchien/moneyKeeper/backend/cmd/sqlhandler/category"
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
		if len(accounts) == 0 {
			log.Info("No rows were returned!")
			w.Write([]byte("[]"))
			return
		}
		accountsStr, _ := json.Marshal(accounts)
		log.WithField("accounts", accounts).Info("accounts...")
		w.Write(accountsStr)

	case "POST":
		d := json.NewDecoder(r.Body)
		p := datamodel.Account{}
		err := d.Decode(&p)
		if err != nil {
			log.WithError(err).Error("Could not Decode incoming request body")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		p, err = account.AddAccount(DBConnection, p)
		if err != nil {
			log.WithError(err).Error("AddAccount DB operation returned an error")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		resultDb, err := json.Marshal(p)
		if err != nil {
			log.WithError(err).Error("Could not marshal DB data into the Account type")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Write(resultDb)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func accountIDHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	log.WithField("id", vars["id"]).Info("accountID to get")

	switch r.Method {
	case "GET":
		// Get account from DB by the account ID
		_, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.WithError(err).WithField("IDparamener", vars["id"]).Error("Could not parse the AccountID")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		getAccount, err := account.GetAccountByID(DBConnection, vars["id"])
		if err == account.ErrNoItemResponse {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		j, err := json.Marshal(getAccount)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.Write(j)

	case "PUT":
		d := json.NewDecoder(r.Body)
		accountUpd := &datamodel.Account{}
		err := d.Decode(accountUpd)
		if err != nil {
			log.WithError(err).Error("Could not decode incoming json body to the Account type")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}

		accountUpd.AccountID, err = strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.WithError(err).WithField("IDparamener", vars["id"]).Error("Could not parse the AccountID")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		err = account.UpdateAccountByID(DBConnection, accountUpd)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		// Returning updated Account info
		j, _ := json.Marshal(accountUpd)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.Write(j)

	case "DELETE":
		iDparamener, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.WithError(err).WithField("IDparamener", iDparamener).Error("Could not parse the AccountID")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		err = account.DeleteAccountByID(DBConnection, vars["id"])
		if err != nil {
			log.WithError(err).WithField("accountID", vars["id"]).Error("Could not delete the Account")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
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

		accountTypesStr, _ := json.Marshal(accountTypes)
		log.WithField("accountTypes", accountTypes).Info("accountTypes...")
		w.Write(accountTypesStr)

	case "POST":
		// Decode the JSON in the body and overwrite 'tom' with it
		d := json.NewDecoder(r.Body)
		p := datamodel.AccountType{}
		err := d.Decode(&p)
		if err != nil {
			log.WithError(err).Error("Could not Decode incoming request body")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		accountType, err := accountType.AddAccountType(DBConnection, p)
		if err != nil {
			log.WithError(err).Error("AddAccountType DB operation returned an error")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		accountTypeStr, err := json.Marshal(accountType)
		if err != nil {
			log.WithError(err).Error("Could not marshal DB data into the AccountType type")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		w.Write(accountTypeStr)

	default:
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func accountTypeIDHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	log.Printf("accountTypeID: %v", vars["id"])

	switch r.Method {
	case "GET":
		// Get accountType from DB
		_, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.WithError(err).WithField("IDparamener", vars["id"]).Error("Could not parse the AccountID")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		getAccountType, err := accountType.GetAccountTypeByID(DBConnection, vars["id"])
		if err == accountType.ErrNoItemResponse {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		getAccountTypeStr, _ := json.Marshal(getAccountType)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.Write(getAccountTypeStr)

	case "PUT":
		d := json.NewDecoder(r.Body)
		accountTypeUpd := &datamodel.AccountType{}
		err := d.Decode(accountTypeUpd)
		if err != nil {
			log.WithError(err).Error("Could not decode incoming json body to the AccountType type")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		accountTypeUpd.TypeID, err = strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.WithError(err).WithField("IDparamener", vars["id"]).Error("Could not parse the AccountID")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		err = accountType.UpdateAccountTypeByID(DBConnection, accountTypeUpd)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		// Returning updated AccountType info
		j, _ := json.Marshal(accountTypeUpd)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.Write(j)

	case "DELETE":
		_, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.WithError(err).WithField("IDparamener", vars["id"]).Error("Could not parse the AccountID")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		err = accountType.DeleteAccountTypeByID(DBConnection, vars["id"])

		if err != nil {
			log.WithError(err).WithField("accountTypeID", vars["id"]).Error("Could not delete the AccountType")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
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
		if len(categories) == 0 {
			log.Info("No rows were returned!")
			w.Write([]byte("[]"))
			return
		}

		categoriesStr, _ := json.Marshal(categories)
		w.Write(categoriesStr)
		return

	case "POST":
		d := json.NewDecoder(r.Body)
		p := datamodel.Category{}
		err := d.Decode(&p)
		if err != nil {
			log.WithError(err).WithField("IncominBody", r.Body).Error("Could not do Decode the response body")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		p, err = category.AddCategory(DBConnection, p)
		if err != nil {
			log.WithError(err).WithField("IncominBody", r.Body).Error("An error has happened during the AddCategory DB operation")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		resultDb, err := json.Marshal(p)
		if err != nil {
			log.WithError(err).Error("Could not marshal AddCategory outcoming object for response")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.Write([]byte(resultDb))
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

		jsonStr, _ := json.Marshal(getCategory)
		w.Write(jsonStr)
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
		categoryUpd.CategoryID, err = strconv.ParseInt(vars["id"], 10, 64)
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

		response, err := json.Marshal(&categoryUpd)
		if err != nil {
			log.WithError(err).Error("Error encoding response object")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return

	case "DELETE":
		categoryID, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.WithError(err).WithField("categoryID", categoryID).Error("Could not parse the Category ID")
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		err = category.DeleteCategoryByID(DBConnection, categoryID)
		if err != nil {
			log.WithError(err).WithField("categoryID", categoryID).Error("Could the category")
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return

	default:
		respondWithError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func main() {
	DBConnection = sqlhandler.DBConnect()

	r := mux.NewRouter()

	r.HandleFunc("/api/", apiHandler)                          // GET saying Hello
	r.HandleFunc("/api/account", accountHandler)               // Handle Accout requests: GET ALL accounts, GET single specified account, POST new account, PUT updates the account, DELETE specified account
	r.HandleFunc("/api/account/{id:[0-9]+}", accountIDHandler) // Handle account type requests: GET single specified account type, PUT to update account type, DELETE specified account type

	r.HandleFunc("/api/account_type", accountTypeHandler)               // Handle account type requests: GET all account type, POST new account type
	r.HandleFunc("/api/account_type/{id:[0-9]+}", accountTypeIDHandler) // Handle account type requests: GET single specified account type, PUT to update account type, DELETE specified account type

	http.HandleFunc("/api/category", categoryHandler)               // GET ALL categories, GET single category info, POST new category, PUT update category, DELETE category
	http.HandleFunc("/api/category/{id:[0-9]+}", categoryIDHandler) // GET ALL categories, GET single category info, POST new category, PUT update category, DELETE category

	// http.HandleFunc("/api/transaction", transactionTypeHandler) // GET ALL transactions, GET single transaction info, POST new transaction, PUT update transaction, DELETE transaction, SEARCH transactions by some criteria

	log.Println("Go!")
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
	result, _ := json.Marshal(data)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	response.Write(result)
}

func respondWithError(response http.ResponseWriter, statusCode int, msg string) {
	respondWithJSON(response, statusCode, map[string]string{"error": msg})
}
