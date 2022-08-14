package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-sql-driver/mysql"
)

var (
	db     *sql.DB
	client *redis.Client
)

type DB_User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DB_Subdomain struct {
	Subdomain string `json:"subdomain"`
	Username  string `json:"username"`
}

type RD_Subdomain struct {
	Subdomain string `json:"subdomain"`
	Operation string `json:"operation"`
}

func main() {
	connectDB()
	connectRedis()

	mux := http.NewServeMux()
	mux.HandleFunc("/user", userHandler)
	mux.HandleFunc("/subdomain", subdomainHandler)
	mux.HandleFunc("/router", routerHandler)

	err := http.ListenAndServe(":"+os.Getenv("API_PORT"), mux)
	log.Fatal(err)
}

func enableSubdomain(subdomain *RD_Subdomain, enabled bool) error {
	enabledStr := "F"
	if enabled {
		enabledStr = "T"
	}
	err := client.HSet("subdomain:"+subdomain.Subdomain, "enabled", enabledStr).Err()

	if err != nil {
		return fmt.Errorf("addSubdomain: %v", err)
	}
	return nil
}

func routerOperation(subdomain *RD_Subdomain) error {
	switch {
	case subdomain.Operation == "enable":
		return enableSubdomain(subdomain, true)
	case subdomain.Operation == "disable":
		return enableSubdomain(subdomain, false)
	}
	return fmt.Errorf("Router operation %s not supported.", subdomain.Operation)
}

func routerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("routerHandler called")
	var subdomain RD_Subdomain
	err := json.NewDecoder(r.Body).Decode(&subdomain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = routerOperation(&subdomain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Router: %+v", subdomain)
}

func addUser(user *DB_User) error {
	_, err := db.Exec("INSERT INTO Users (username, password) VALUES (?, ?)",
		user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("addUser: %v", err)
	}
	return nil
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	var user DB_User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = addUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "User: %+v", user)
}

func addSubdomain(subdomain *DB_Subdomain) error {
	_, err := db.Exec("INSERT INTO Subdomains (subdomain, username) VALUES (?, ?)",
		subdomain.Subdomain, subdomain.Username)
	if err != nil {
		return fmt.Errorf("addSubdomain: %v", err)
	}
	return nil
}

func subdomainHandler(w http.ResponseWriter, r *http.Request) {
	var subdomain DB_Subdomain
	err := json.NewDecoder(r.Body).Decode(&subdomain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = addSubdomain(&subdomain)
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
		Net:    "tcp",
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

func connectRedis() {
	// Capture connection properties
	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RD_HOST") + ":6379",
		Password: "",
		DB:       0,
	})
	var err error
	for {
		_, err = client.Ping().Result()
		if err != nil {
			fmt.Printf("connectRedis : %s\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	fmt.Println("connectRedis : Connected!")
}
