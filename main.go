package main

import (
	"database/sql"
	"flag"
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

func init() {
	log.SetFlags(log.Lshortfile)
	var err error
	funcmap := template.FuncMap{
		"formatDuration": formatDuration,
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
	}
	tpl, err = template.New("MyTemplate").Funcs(funcmap).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalln(err)
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connection := fmt.Sprintf("%s:%s@unix(/var/run/mysqld/mysqld.sock)/%s?parseTime=true", user, password, dbName)
	db, err = sql.Open("mysql", connection)
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	filename := flag.String("init-users", "", "The json file of user's data to init database.")
	isDeviced := flag.Bool("init-deviced", false, "Whether should program init the devices data")
	flag.Parse()
	if *filename != "" {
		initUsers(*filename)
	}
	if *isDeviced {
		initDevices()
	}
	initCheck()
	log.Println("Server runs on http://localhost:8080")
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.Handle("/favicon.ico", http.NotFoundHandler())
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/bookings/new", bookingForm)
	r.HandleFunc("/bookings", bookingList).Methods("GET")
	r.HandleFunc("/bookings", newBooking).Methods("POST")
	r.HandleFunc("/bookings/{id:[0-9]+}", bookingPage)
	r.HandleFunc("/bookings/{id:[0-9]+}/lend", lendForm)
	r.HandleFunc("/bookings/{id:[0-9]+}/records", recordList)
	r.HandleFunc("/records", newRecord).Methods("POST")
	r.HandleFunc("/records", handleReturnDevice).Methods("DELETE")
	r.HandleFunc("/return", handleReturnList)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/check", checkPage)
	checkErr(http.ListenAndServe(":8080", r), "Start server fatal: ")
}
