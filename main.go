package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

type Customers struct {
	ID        uint32 `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Role      string `json:"role,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Contacted bool   `json:"contacted"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "example_admin"
	password = "somepassword"
	dbname   = "crm_customers"
)

var DB *sql.DB

func connectDatabase() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	DB = db
	//open connection to the database as open doesnt do that
	if err := DB.Ping(); err != nil {
		fmt.Println("DB Error")
	}

	fmt.Println("Successfully connected!")
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var customers []Customers
	result, err := DB.Query("SELECT * from customers")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var customer Customers
		err := result.Scan(&customer.ID, &customer.Name, &customer.Role, &customer.Email, &customer.Phone, &customer.Contacted)
		if err != nil {
			panic(err.Error())
		}
		customers = append(customers, customer)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	sqlStatement := `SELECT * FROM customers WHERE id=$1;`
	var customer Customers
	row := DB.QueryRow(sqlStatement, id)
	switch err := row.Scan(&customer.ID, &customer.Name, &customer.Role, &customer.Email, &customer.Phone, &customer.Contacted); err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Customer with required id was not found!")
	case nil:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customer)
	default:
		panic(err)
	}

}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := DB.Prepare("INSERT INTO customers VALUES($1, $2, $3, $4, $5, $6)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	customer := &Customers{}
	buffer := []byte(body)
	if err := json.Unmarshal(buffer, customer); err != nil {
		fmt.Printf("error unmarshaling in add customer JSON: %v\n", err)
	}
	_, err = stmt.Exec(customer.ID, customer.Name, customer.Role, customer.Email, customer.Phone, customer.Contacted)
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "New customer was added")
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := DB.Prepare("UPDATE customers SET name = $1, role = $2, email = $3, phone = $4, contacted = $5 WHERE id = $6")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	customer := &Customers{}
	buffer := []byte(body)
	if err := json.Unmarshal(buffer, customer); err != nil {
		fmt.Printf("error unmarshaling in update customer JSON: %v\n", err)
	}
	_, err = stmt.Exec(customer.Name, customer.Role, customer.Email, customer.Phone, customer.Contacted, params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		panic(err.Error())
	} else {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Customer with ID = %s was updated", params["id"])
	}
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := DB.Prepare("DELETE FROM customers WHERE id = $1")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		panic(err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Customer with ID = %s was deleted", params["id"])
	}
}

func main() {
	connectDatabase()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PATCH")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

	fmt.Println("Server is starting on port 3000...")
	http.ListenAndServe(":3000", router)
}
