package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	blockchain "kilimo-chain/block"
)

// Initialize blockchain instance
var blockchainInstance *blockchain.Blockchain

func init() {
	blockchainInstance = blockchain.NewBlockchain()
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve images
	imageServer := http.FileServer(http.Dir("Images"))
	http.Handle("/Images/", http.StripPrefix("/Images/", imageServer))

	// Register handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/collaborators", collaboratorsHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/addTransaction", AddTransactionHandler)
	http.HandleFunc("/blockchain", BlockchainHandler)
	http.HandleFunc("/signup", SignUphandler)
	http.HandleFunc("/signin", Loginpageload)
	http.HandleFunc("/signup/farmers", handleSignupFarmers)
	http.HandleFunc("/signup/cooperates", handleSignupCooperates)
	http.HandleFunc("/signin/farmers", handleSigninFarmers)
	http.HandleFunc("/signin/cooperates", handleSigninCooperates)
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "Dashboard.html")
	})

	// Start the server
	fmt.Println("Server started at :9674")
	log.Fatal(http.ListenAndServe(":9674", nil))
}

// Signup for Farmers
func handleSignupFarmers(w http.ResponseWriter, r *http.Request) {
	log.Println("Received signup request")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	firstName := r.FormValue("first_name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if password != confirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	user := User{
		FirstName: firstName,
		Email:     email,
		Password:  password,
	}

	if err := user.Save(); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

// Signup for Cooperates
func handleSignupCooperates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	company := r.FormValue("company")
	email := r.FormValue("email")
	password := r.FormValue("password")
	//phone := r.FormValue("phone")

	user := User{
		Company:  company,
		Email:    email,
		Password: password,
	}

	if err := user.Save(); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

// Signin for Farmers
func handleSigninFarmers(w http.ResponseWriter, r *http.Request) {

	log.Println("Received signin request")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	users, err := LoadUsers()
	if err != nil {
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		return
	}

	for _, user := range users {
		if user.Email == email && user.Password == password {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

// Signin for Cooperates
func handleSigninCooperates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	users, err := LoadUsers()
	if err != nil {
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		return
	}

	for _, user := range users {
		if user.Email == email && user.Password == password {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

// // Placeholder for dashboard handler
// func handleDashboard(w http.ResponseWriter, r *http.Request) {
// 	// Implement dashboard rendering or redirection logic here
// }

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "about.html")
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "products.html")
}

func collaboratorsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "collaborators.html")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = u.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func AddTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var transaction blockchain.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction.Timestamp = time.Now().Unix()
	blockData := blockchain.BlockData{
		Transactions: []blockchain.Transaction{transaction},
	}
	blockchainInstance.AddBlock(blockData)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

func BlockchainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewEncoder(w).Encode(blockchainInstance)
	if err != nil {
		http.Error(w, "Failed to encode blockchain", http.StatusInternalServerError)
		return
	}
}

// Handlers for signup and signin pages
func SignUphandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("Signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func Loginpageload(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("Signin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
