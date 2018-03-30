package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/michaeljs1990/sqlitestore"

	"database/sql"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
	"os/signal"
	"syscall"
)

var db *sql.DB
var sessionStore *sqlitestore.SqliteStore
const (
	baseUrl = "https://pics.sh1kel.com/files/"
	uploadDir = "./upload/"
)

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
	sessionStore, err = sqlitestore.NewSqliteStore("./auth.sqlite", "sessions", "/", 3600, []byte("lohM2oofaef7eyoophahcohngaihe4ah"))
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
	router.Handle("/auth", AuthHandler).Methods("GET")
	router.Handle("/signin", SignInHandler).Methods("POST")
	router.Handle("/signup", SignUpHandler).Methods("POST")
	router.Handle("/logout", LogoutHandler).Methods("GET")
	router.Handle("/files", FilesHandler).Methods("GET")
	router.Handle("/upload", UploadHandler).Methods("POST")
	router.Handle("/user", RootHandler).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	router.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir("./upload/"))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))





	http.ListenAndServe("127.0.0.1:8081", handlers.LoggingHandler(os.Stdout, router))

}
