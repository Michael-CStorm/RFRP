package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"log"
	"net/http"
	"encoding/json"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type Subdomain struct {
	Subdomain string `json:"subdomain"`
	Username string `json:"username"`
}

func main() {
	connectDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/user", addUserEntryPoint)
	mux.HandleFunc("/subdomain", addSubdomainEntryPoint)

	err := http.ListenAndServe(":5000", mux)
	log.Fatal(err)
}

func addUser( user *User ) error {
	_, err := db.Exec("INSERT INTO Users (username, password) VALUES (?, ?)",
							user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("addUser: %v", err)
	}
	return nil
}

func addUserEntryPoint(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = addUser( &user )
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "User: %+v", user)
}

func addSubdomain( subdomain *Subdomain ) error {
	_, err := db.Exec("INSERT INTO Subdomains (subdomain, username) VALUES (?, ?)",
							subdomain.Subdomain, subdomain.Username)
	if err != nil {
		return fmt.Errorf("addSubdomain: %v", err)
	}
	return nil
}

func addSubdomainEntryPoint(w http.ResponseWriter, r *http.Request) {
	var subdomain Subdomain
	err := json.NewDecoder(r.Body).Decode(&subdomain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = addSubdomain( &subdomain )
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Subdomain: %+v", subdomain)
}

func connectDB() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:	"tcp",
		Addr:   os.Getenv("DB_HOST") + ":3306",
		DBName: "user_api",
	}
	// Get a database handle.
	var err error
	for {
		db, err = sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			fmt.Printf("connectDB : %s\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		err := db.Ping()
		if err != nil {
			fmt.Printf("connectDB : %s\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	fmt.Println("connectDB : Connected!")
}

