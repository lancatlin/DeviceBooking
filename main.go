package main

import (
	"database/sql"
	"errors"
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

var dbName string
var dbUser string
var dbPassword string
var dbHost string
var port string

func init() {
	log.SetFlags(log.Lshortfile)
	dbName = os.Getenv("DB_NAME")
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbHost = os.Getenv("DB_HOST")
	port = os.Getenv("PORT")
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

func loadDB() (err error) {
	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbName)
	i := 0
	for {
		db, err = sql.Open("mysql", connection)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		if i == 10 {
			log.Fatalln("Cannot connect to database")
		}
		log.Println(err)
		time.Sleep(10 * time.Second)
		i++
	}
	row := db.QueryRow(`
	SELECT ID FROM Users WHERE Name = 'Admin';
	`)
	var id int64
	if err = row.Scan(&id); err == sql.ErrNoRows {
		return errors.New("Admin hasn't set up yet")
	} else if err != nil {
		log.Fatalln("Database didn't set up.")
	}
	return nil
}

func main() {
	var err error
	if err = loadDB(); err != nil {
		log.Print("Init mode: ")
		initMode = true
	}
	defer db.Close()
	initCheck()
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
	r.HandleFunc("/doc/{filename}", docs)
	r.HandleFunc("/init", initPage).Methods("GET")
	r.HandleFunc("/init", initAdmin).Methods("POST")
	log.Println("Server runs on http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalln(err)
	}
}
