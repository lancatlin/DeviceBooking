package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

var db *sql.DB
var tpl *template.Template
var initMode = false

func init() {
	log.SetFlags(log.Lshortfile)
	var err error
	funcmap := template.FuncMap{
		"formatDuration": formatDuration,
		"formatDate": func(t time.Time) string {
			return t.Local().Format("2006-01-02")
		},
		"formatTime": func(t time.Time) string {
			return t.Local().Format("15:04")
		},
	}
	tpl, err = template.New("MyTemplate").Funcs(funcmap).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalln(err)
	}
}

func loadDB() (db *sql.DB, err error) {
	file, err := os.Open("env.json")
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(file)
	vars := new(struct {
		DBName   string
		User     string
		Password string
	})
	if err := dec.Decode(vars); err != nil {
		return nil, err
	}
	connection := fmt.Sprintf("%s:%s@unix(/var/run/mysqld/mysqld.sock)/%s?parseTime=true", vars.User, vars.Password, vars.DBName)
	db, err = sql.Open("mysql", connection)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func handleInitMode() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/init", 303)
	})
	r.HandleFunc("/init", initDBPage).Methods("GET")
	r.HandleFunc("/init", handleInitDB).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var err error
	if db, err = loadDB(); err != nil {
		log.Println(err)
		log.Println("Init mode: server runs on http://localhost:8000/init")
		handleInitMode()
	}
	initCheck()
	log.Println("Server runs on http://localhost:8000")
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.Handle("/favicon.ico", http.NotFoundHandler())
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/bookings/new", bookingForm)
	r.HandleFunc("/bookings/lending", handleLendingList)
	r.HandleFunc("/bookings/overdue", overdue)
	r.HandleFunc("/bookings", bookingList).Methods("GET")
	r.HandleFunc("/bookings", newBooking).Methods("POST")
	r.HandleFunc("/bookings/{id:[0-9]+}", bookingPage)
	r.HandleFunc("/bookings/{id:[0-9]+}/lend", lendForm)
	r.HandleFunc("/bookings/{id:[0-9]+}/records", recordList)
	r.HandleFunc("/records", newRecord).Methods("POST")
	r.HandleFunc("/records", handleReturnDevice).Methods("DELETE")
	r.HandleFunc("/devices", devices)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/check", checkPage)
	r.HandleFunc("/users", users).Methods("GET")
	r.HandleFunc("/users", signUp).Methods("POST")
	r.HandleFunc("/users/upload", importUsers).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", setPermission).Methods("PUT").Queries("permission", "")
	r.HandleFunc("/users/{id:[0-9]+}/set-password", resetPassword).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}/set-password", resetPasswordPage).Methods("GET")
	checkErr(http.ListenAndServe(":8000", r), "Start server fatal: ")
}
