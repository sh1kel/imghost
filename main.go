package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"github.com/labstack/gommon/log"
	"os/signal"
	"syscall"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./auth.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	initSql := `
	create table if not exists auth(id integer primary key autoincrement, username text, password text)
	`
	_, err = db.Exec(initSql)
	if err != nil {
		log.Info(err)
	}
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for s := range sigChan {
			log.Info("Got stop signal: " + s.String())
			db.Close()
			os.Exit(0)
		}
	}()

	router := mux.NewRouter()
	initDB()
	router.Handle("/", http.FileServer(http.Dir("./html/")))
	router.Handle("/signin", SigninHandler).Methods("POST")
	router.Handle("/signup", SignupHandler).Methods("POST")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8081", handlers.LoggingHandler(os.Stdout, router))

}
