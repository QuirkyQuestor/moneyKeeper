package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/winchien/moneyKeeper/backend/cmd/datamodel"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		w.Write([]byte("Welcome to MoneyKeeper"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func accountHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		// Fake account
		fakeAcc := datamodel.Account{
			AccountID:   1,
			TypeID:      1,
			Name:        "MyFakeAccount",
			Description: "",
			Active:      true,
		}

		j, _ := json.Marshal(fakeAcc)
		w.Write(j)
	case "POST":
		// Decode the JSON in the body and overwrite 'tom' with it
		d := json.NewDecoder(r.Body)
		p := &datamodel.Account{}
		err := d.Decode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// tom := p
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func main() {
	http.HandleFunc("/api/", apiHandler)            // GET saying Hello
	http.HandleFunc("/api/account", accountHandler) // GET ALL accounts, GET single specified account, POST new account, PUT updates the account, DELETE specified account
	// http.HandleFunc("/api/account_type", accountTypeHandler)    // GET ALL account types, GET single specified account type, POST new account type, PUT to update account type, DELETE specified account type
	// http.HandleFunc("/api/category", categoryTypeHandler)       // GET ALL categories, GET single category info, POST new category, PUT update category, DELETE category
	// http.HandleFunc("/api/transaction", transactionTypeHandler) // GET ALL transactions, GET single transaction info, POST new transaction, PUT update transaction, DELETE transaction, SEARCH transactions by some criteria

	log.Println("Go!")
	http.ListenAndServe(":8080", nil)
}
