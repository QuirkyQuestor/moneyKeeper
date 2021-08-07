package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/winchien/moneyKeeper/backend/cmd/datamodel"
	"github.com/winchien/moneyKeeper/backend/cmd/sqlhandler"
)

var DBConnection *sql.DB

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

func accountTypeHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		// Fake account
		var accountTypes []datamodel.AccountType
		bdStatement, err := DBConnection.Prepare("SELECT typeId, name, description FROM AccountType")
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
			// ...
			var typeId int
			var name string
			var description sql.NullString
			err := rows.Scan(&typeId, &name, &description)

			if err != nil {
				log.Fatal(err)
			}
			aaa := datamodel.AccountType{
				TypeID:      int64(typeId),
				Name:        name,
				Description: description.String,
			}
			accountTypes = append(accountTypes, aaa)
			// log.Println(id, name)

		}
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}

		j, _ := json.Marshal(accountTypes)
		w.Write(j)
	case "POST":
		// Decode the JSON in the body and overwrite 'tom' with it
		d := json.NewDecoder(r.Body)
		p := &datamodel.AccountType{}
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
	DBConnection = sqlhandler.DBConnect()

	http.HandleFunc("/api/", apiHandler)                     // GET saying Hello
	http.HandleFunc("/api/account", accountHandler)          // GET ALL accounts, GET single specified account, POST new account, PUT updates the account, DELETE specified account
	http.HandleFunc("/api/account_type", accountTypeHandler) // GET ALL account types, GET single specified account type, POST new account type, PUT to update account type, DELETE specified account type
	// http.HandleFunc("/api/category", categoryTypeHandler)       // GET ALL categories, GET single category info, POST new category, PUT update category, DELETE category
	// http.HandleFunc("/api/transaction", transactionTypeHandler) // GET ALL transactions, GET single transaction info, POST new transaction, PUT update transaction, DELETE transaction, SEARCH transactions by some criteria

	log.Println("Go!")
	http.ListenAndServe(":8080", nil)
}
